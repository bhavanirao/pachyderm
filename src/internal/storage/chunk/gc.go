package chunk

import (
	"context"
	"time"

	"github.com/pachyderm/pachyderm/v2/src/internal/errors"
	"github.com/sirupsen/logrus"
)

// GarbageCollector removes unused chunks from object storage
type GarbageCollector struct {
	s      *Storage
	log    *logrus.Logger
	period time.Duration
}

// NewGC returns a new garbage collector operating on s
func NewGC(s *Storage, d time.Duration, log *logrus.Logger) *GarbageCollector {
	return &GarbageCollector{s: s, log: log, period: d}
}

// RunForever calls RunOnce until the context is cancelled, logging any errors.
func (gc *GarbageCollector) RunForever(ctx context.Context) error {
	ticker := time.NewTicker(gc.period)
	defer ticker.Stop()
	for {
		if err := gc.RunOnce(ctx); err != nil {
			select {
			case <-ctx.Done():
				return err
			default:
			}
			gc.log.Errorf("during chunk GC: %v", err)
		}
		select {
		case <-ctx.Done():
			return errors.EnsureStack(ctx.Err())
		case <-ticker.C:
		}
	}
}

// RunOnce runs 1 cycle of garbage collection.
func (gc *GarbageCollector) RunOnce(ctx context.Context) (retErr error) {
	rows, err := gc.s.db.QueryxContext(ctx, `
	SELECT chunk_id, gen, uploaded FROM storage.chunk_objects
	WHERE tombstone = true
	`)
	if err != nil {
		return errors.EnsureStack(err)
	}
	defer func() {
		if err := rows.Close(); retErr == nil {
			retErr = err
		}
	}()
	for rows.Next() {
		var ent Entry
		if err := rows.StructScan(&ent); err != nil {
			return errors.EnsureStack(err)
		}
		if !ent.Uploaded {
			gc.log.Warnf("possibility for untracked chunk %s", chunkPath(ent.ChunkID, ent.Gen))
		}
		if err := gc.deleteOne(ctx, ent); err != nil {
			return err
		}
		gc.log.WithFields(logrus.Fields{
			"chunk_id": ent.ChunkID,
			"gen":      ent.Gen,
		}).Infof("deleting object for chunk entry")
	}
	return errors.EnsureStack(rows.Err())
}

func (gc *GarbageCollector) deleteOne(ctx context.Context, ent Entry) error {
	if err := gc.deleteObject(ctx, ent.ChunkID, ent.Gen); err != nil {
		return err
	}
	return gc.deleteEntry(ctx, ent.ChunkID, ent.Gen)
}

func (gc *GarbageCollector) deleteObject(ctx context.Context, chunkID ID, gen uint64) error {
	return errors.EnsureStack(gc.s.store.Delete(ctx, chunkKey(chunkID, gen)))
}

func (gc *GarbageCollector) deleteEntry(ctx context.Context, chunkID ID, gen uint64) error {
	_, err := gc.s.db.ExecContext(ctx, `
	DELETE FROM storage.chunk_objects
	WHERE chunk_id = $1 AND gen = $2 AND tombstone = TRUE
	`, chunkID, gen)
	return errors.EnsureStack(err)
}
