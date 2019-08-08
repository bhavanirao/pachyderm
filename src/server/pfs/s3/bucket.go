package s3

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gogo/protobuf/types"
	"github.com/pachyderm/ohmyglob"
	"github.com/pachyderm/pachyderm/src/client"
	pfsClient "github.com/pachyderm/pachyderm/src/client/pfs"
	pfsServer "github.com/pachyderm/pachyderm/src/server/pfs"
	"github.com/pachyderm/pachyderm/src/server/pkg/errutil"
	"github.com/pachyderm/s2"
	"github.com/sirupsen/logrus"
)

func newContents(fileInfo *pfsClient.FileInfo) (s2.Contents, error) {
	t, err := types.TimestampFromProto(fileInfo.Committed)
	if err != nil {
		return s2.Contents{}, err
	}

	return s2.Contents{
		Key:          fileInfo.File.Path,
		LastModified: t,
		ETag:         fmt.Sprintf("%x", fileInfo.Hash),
		Size:         fileInfo.SizeBytes,
		StorageClass: globalStorageClass,
		Owner:        defaultUser,
	}, nil
}

func newCommonPrefixes(dir string) s2.CommonPrefixes {
	return s2.CommonPrefixes{
		Prefix: fmt.Sprintf("%s/", dir),
		Owner:  defaultUser,
	}
}

type bucketController struct {
	pc     *client.APIClient
	logger *logrus.Entry
}

func newBucketController(pc *client.APIClient, logger *logrus.Entry) *bucketController {
	c := bucketController{
		pc:     pc,
		logger: logger,
	}

	return &c
}

func (c *bucketController) GetLocation(r *http.Request, bucket string) (location string, err error) {
	repo, branch, err := bucketArgs(r, bucket)
	if err != nil {
		return
	}

	_, err = c.pc.InspectBranch(repo, branch)
	if err != nil {
		err = maybeNotFoundError(r, err)
		return
	}

	location = globalLocation
	return
}

func (c *bucketController) ListObjects(r *http.Request, bucket, prefix, marker, delimiter string, maxKeys int) (contents []s2.Contents, commonPrefixes []s2.CommonPrefixes, isTruncated bool, err error) {
	repo, branch, err := bucketArgs(r, bucket)
	if err != nil {
		return
	}

	if delimiter != "" && delimiter != "/" {
		err = invalidDelimiterError(r)
		return
	}

	// ensure the branch exists and has a head
	branchInfo, err := c.pc.InspectBranch(repo, branch)
	if err != nil {
		err = maybeNotFoundError(r, err)
		return
	}
	if branchInfo.Head == nil {
		// if there's no head commit, just print an empty list of files
		return
	}

	recursive := delimiter == ""
	var pattern string
	if recursive {
		pattern = fmt.Sprintf("%s**", glob.QuoteMeta(prefix))
	} else {
		pattern = fmt.Sprintf("%s*", glob.QuoteMeta(prefix))
	}

	err = c.pc.GlobFileF(repo, branch, pattern, func(fileInfo *pfsClient.FileInfo) error {
		if fileInfo.FileType == pfsClient.FileType_DIR {
			if fileInfo.File.Path == "/" {
				// skip the root directory
				return nil
			}
			if recursive {
				// skip directories if recursing
				return nil
			}
		} else if fileInfo.FileType != pfsClient.FileType_FILE {
			// skip anything that isn't a file or dir
			return nil
		}

		fileInfo.File.Path = fileInfo.File.Path[1:] // strip leading slash

		if !strings.HasPrefix(fileInfo.File.Path, prefix) {
			return nil
		}
		if fileInfo.File.Path <= marker {
			return nil
		}

		if len(contents)+len(commonPrefixes) >= maxKeys {
			if maxKeys > 0 {
				isTruncated = true
			}
			return errutil.ErrBreak
		}
		if fileInfo.FileType == pfsClient.FileType_FILE {
			c, err := newContents(fileInfo)
			if err != nil {
				return err
			}

			contents = append(contents, c)
		} else {
			commonPrefixes = append(commonPrefixes, newCommonPrefixes(fileInfo.File.Path))
		}

		return nil
	})

	return
}

func (c *bucketController) CreateBucket(r *http.Request, bucket string) error {
	repo, branch, err := bucketArgs(r, bucket)
	if err != nil {
		return err
	}

	err = c.pc.CreateRepo(repo)
	if err != nil {
		if errutil.IsAlreadyExistError(err) {
			// Bucket already exists - this is not an error so long as the
			// branch being created is new. Verify if that is the case now,
			// since PFS' `CreateBranch` won't error out.
			_, err := c.pc.InspectBranch(repo, branch)
			if err != nil {
				if !pfsServer.IsBranchNotFoundErr(err) {
					return s2.InternalError(r, err)
				}
			} else {
				return s2.BucketAlreadyOwnedByYouError(r)
			}
		} else if errutil.IsInvalidNameError(err) {
			return s2.InvalidBucketNameError(r)
		} else {
			return s2.InternalError(r, err)
		}
	}

	err = c.pc.CreateBranch(repo, branch, "", nil)
	if err != nil {
		if errutil.IsInvalidNameError(err) {
			return s2.InvalidBucketNameError(r)
		}
		return s2.InternalError(r, err)
	}

	return nil
}

func (c *bucketController) DeleteBucket(r *http.Request, bucket string) error {
	repo, branch, err := bucketArgs(r, bucket)
	if err != nil {
		return err
	}

	// `DeleteBranch` does not return an error if a non-existing branch is
	// deleting. So first, we verify that the branch exists so we can
	// otherwise return a 404.
	branchInfo, err := c.pc.InspectBranch(repo, branch)
	if err != nil {
		return maybeNotFoundError(r, err)
	}

	if branchInfo.Head != nil {
		hasFiles := false
		err = c.pc.Walk(branchInfo.Branch.Repo.Name, branchInfo.Head.ID, "", func(fileInfo *pfsClient.FileInfo) error {
			if fileInfo.FileType == pfsClient.FileType_FILE {
				hasFiles = true
				return errutil.ErrBreak
			}
			return nil
		})
		if err != nil {
			return s2.InternalError(r, err)
		}

		if hasFiles {
			return s2.BucketNotEmptyError(r)
		}
	}

	err = c.pc.DeleteBranch(repo, branch, false)
	if err != nil {
		return s2.InternalError(r, err)
	}

	repoInfo, err := c.pc.InspectRepo(repo)
	if err != nil {
		return s2.InternalError(r, err)
	}

	// delete the repo if this was the last branch
	if len(repoInfo.Branches) == 0 {
		err = c.pc.DeleteRepo(repo, false)
		if err != nil {
			return s2.InternalError(r, err)
		}
	}

	return nil
}

func (c *bucketController) ListObjectVersions(r *http.Request, repo, prefix, keyMarker, versionIDMarker string, delimiter string, maxKeys int) (versions []s2.Version, deleteMarkers []s2.DeleteMarker, isTruncated bool, err error) {
	err = s2.NotImplementedError(r)
	return
}

func (c *bucketController) GetBucketVersioning(r *http.Request, repo string) (status string, err error) {
	return s2.VersioningEnabled, nil
}

func (c *bucketController) SetBucketVersioning(r *http.Request, repo, status string) error {
	return s2.NotImplementedError(r)
}
