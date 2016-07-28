// Code generated by protoc-gen-go.
// source: server/pfs/db/persist/persist.proto
// DO NOT EDIT!

/*
Package persist is a generated protocol buffer package.

It is generated from these files:
	server/pfs/db/persist/persist.proto

It has these top-level messages:
	Clock
	ClockID
	BranchClock
	Repo
	Branch
	BlockRef
	Diff
	Commit
*/
package persist

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "go.pedge.io/pb/go/google/protobuf"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Clock struct {
	// a document either has these two fields
	Branch string `protobuf:"bytes,1,opt,name=branch" json:"branch,omitempty"`
	Clock  uint64 `protobuf:"varint,2,opt,name=clock" json:"clock,omitempty"`
}

func (m *Clock) Reset()                    { *m = Clock{} }
func (m *Clock) String() string            { return proto.CompactTextString(m) }
func (*Clock) ProtoMessage()               {}
func (*Clock) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type ClockID struct {
	ID     string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Repo   string `protobuf:"bytes,2,opt,name=repo" json:"repo,omitempty"`
	Branch string `protobuf:"bytes,3,opt,name=branch" json:"branch,omitempty"`
	Clock  uint64 `protobuf:"varint,4,opt,name=clock" json:"clock,omitempty"`
}

func (m *ClockID) Reset()                    { *m = ClockID{} }
func (m *ClockID) String() string            { return proto.CompactTextString(m) }
func (*ClockID) ProtoMessage()               {}
func (*ClockID) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type BranchClock struct {
	Clocks []*Clock `protobuf:"bytes,1,rep,name=clocks" json:"clocks,omitempty"`
}

func (m *BranchClock) Reset()                    { *m = BranchClock{} }
func (m *BranchClock) String() string            { return proto.CompactTextString(m) }
func (*BranchClock) ProtoMessage()               {}
func (*BranchClock) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *BranchClock) GetClocks() []*Clock {
	if m != nil {
		return m.Clocks
	}
	return nil
}

type Repo struct {
	Name    string                     `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Created *google_protobuf.Timestamp `protobuf:"bytes,2,opt,name=created" json:"created,omitempty"`
}

func (m *Repo) Reset()                    { *m = Repo{} }
func (m *Repo) String() string            { return proto.CompactTextString(m) }
func (*Repo) ProtoMessage()               {}
func (*Repo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *Repo) GetCreated() *google_protobuf.Timestamp {
	if m != nil {
		return m.Created
	}
	return nil
}

type Branch struct {
	ID   string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Repo string `protobuf:"bytes,2,opt,name=repo" json:"repo,omitempty"`
	Name string `protobuf:"bytes,3,opt,name=name" json:"name,omitempty"`
}

func (m *Branch) Reset()                    { *m = Branch{} }
func (m *Branch) String() string            { return proto.CompactTextString(m) }
func (*Branch) ProtoMessage()               {}
func (*Branch) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

type BlockRef struct {
	Hash  string `protobuf:"bytes,1,opt,name=hash" json:"hash,omitempty"`
	Lower uint64 `protobuf:"varint,2,opt,name=lower" json:"lower,omitempty"`
	Upper uint64 `protobuf:"varint,3,opt,name=upper" json:"upper,omitempty"`
}

func (m *BlockRef) Reset()                    { *m = BlockRef{} }
func (m *BlockRef) String() string            { return proto.CompactTextString(m) }
func (*BlockRef) ProtoMessage()               {}
func (*BlockRef) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

type Diff struct {
	ID       string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Repo     string `protobuf:"bytes,2,opt,name=repo" json:"repo,omitempty"`
	CommitID string `protobuf:"bytes,3,opt,name=commit_id,json=commitId" json:"commit_id,omitempty"`
	Path     string `protobuf:"bytes,4,opt,name=path" json:"path,omitempty"`
	// block_refs and delete cannot both be set
	BlockRefs    []*BlockRef    `protobuf:"bytes,5,rep,name=block_refs,json=blockRefs" json:"block_refs,omitempty"`
	Delete       bool           `protobuf:"varint,6,opt,name=delete" json:"delete,omitempty"`
	Size         uint64         `protobuf:"varint,7,opt,name=size" json:"size,omitempty"`
	BranchClocks []*BranchClock `protobuf:"bytes,8,rep,name=branch_clocks,json=branchClocks" json:"branch_clocks,omitempty"`
}

func (m *Diff) Reset()                    { *m = Diff{} }
func (m *Diff) String() string            { return proto.CompactTextString(m) }
func (*Diff) ProtoMessage()               {}
func (*Diff) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *Diff) GetBlockRefs() []*BlockRef {
	if m != nil {
		return m.BlockRefs
	}
	return nil
}

func (m *Diff) GetBranchClocks() []*BranchClock {
	if m != nil {
		return m.BranchClocks
	}
	return nil
}

type Commit struct {
	ID           string                     `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Repo         string                     `protobuf:"bytes,2,opt,name=repo" json:"repo,omitempty"`
	BranchClocks []*BranchClock             `protobuf:"bytes,3,rep,name=branch_clocks,json=branchClocks" json:"branch_clocks,omitempty"`
	Started      *google_protobuf.Timestamp `protobuf:"bytes,4,opt,name=started" json:"started,omitempty"`
	Finished     *google_protobuf.Timestamp `protobuf:"bytes,5,opt,name=finished" json:"finished,omitempty"`
	Cancelled    bool                       `protobuf:"varint,6,opt,name=cancelled" json:"cancelled,omitempty"`
	Provenance   []string                   `protobuf:"bytes,7,rep,name=provenance" json:"provenance,omitempty"`
}

func (m *Commit) Reset()                    { *m = Commit{} }
func (m *Commit) String() string            { return proto.CompactTextString(m) }
func (*Commit) ProtoMessage()               {}
func (*Commit) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *Commit) GetBranchClocks() []*BranchClock {
	if m != nil {
		return m.BranchClocks
	}
	return nil
}

func (m *Commit) GetStarted() *google_protobuf.Timestamp {
	if m != nil {
		return m.Started
	}
	return nil
}

func (m *Commit) GetFinished() *google_protobuf.Timestamp {
	if m != nil {
		return m.Finished
	}
	return nil
}

func init() {
	proto.RegisterType((*Clock)(nil), "Clock")
	proto.RegisterType((*ClockID)(nil), "ClockID")
	proto.RegisterType((*BranchClock)(nil), "BranchClock")
	proto.RegisterType((*Repo)(nil), "Repo")
	proto.RegisterType((*Branch)(nil), "Branch")
	proto.RegisterType((*BlockRef)(nil), "BlockRef")
	proto.RegisterType((*Diff)(nil), "Diff")
	proto.RegisterType((*Commit)(nil), "Commit")
}

func init() { proto.RegisterFile("server/pfs/db/persist/persist.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 460 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x8c, 0x52, 0x4d, 0x8f, 0xd3, 0x40,
	0x0c, 0x55, 0x9b, 0x34, 0x4d, 0xdc, 0x85, 0xc3, 0x08, 0xa1, 0x68, 0x41, 0x80, 0xc2, 0xa5, 0x17,
	0x12, 0xb1, 0x7c, 0x9c, 0xd1, 0xee, 0x5e, 0x96, 0x13, 0x1a, 0x71, 0xe3, 0x50, 0x25, 0xa9, 0xb3,
	0x8d, 0x48, 0x9a, 0x68, 0x26, 0x5b, 0x24, 0xfe, 0x02, 0x7f, 0x93, 0x1f, 0x82, 0xc7, 0x93, 0xa1,
	0x91, 0x40, 0x6a, 0x4f, 0xb1, 0x5f, 0xfc, 0xec, 0x37, 0xcf, 0x86, 0xd7, 0x1a, 0xd5, 0x01, 0x55,
	0xd6, 0x57, 0x3a, 0xdb, 0x16, 0x59, 0x8f, 0x4a, 0xd7, 0x7a, 0x70, 0xdf, 0xb4, 0x57, 0xdd, 0xd0,
	0x5d, 0xbe, 0xbc, 0xef, 0xba, 0xfb, 0x06, 0x33, 0xce, 0x8a, 0x87, 0x2a, 0x1b, 0xea, 0x16, 0xf5,
	0x90, 0xb7, 0xbd, 0x2d, 0x48, 0x3e, 0xc0, 0xe2, 0xa6, 0xe9, 0xca, 0xef, 0xe2, 0x29, 0x04, 0x85,
	0xca, 0xf7, 0xe5, 0x2e, 0x9e, 0xbd, 0x9a, 0xad, 0x23, 0x39, 0x66, 0xe2, 0x09, 0x2c, 0x4a, 0x53,
	0x10, 0xcf, 0x09, 0xf6, 0xa5, 0x4d, 0x92, 0x6f, 0xb0, 0x64, 0xda, 0xdd, 0xad, 0x78, 0x0c, 0xf3,
	0x7a, 0x3b, 0x92, 0x28, 0x12, 0x02, 0x7c, 0x85, 0x7d, 0xc7, 0xf5, 0x91, 0xe4, 0x78, 0xd2, 0xdc,
	0xfb, 0x7f, 0x73, 0x7f, 0xda, 0xfc, 0x0d, 0xac, 0xae, 0xf9, 0xbf, 0x55, 0xf6, 0x02, 0x02, 0xc6,
	0x35, 0x0d, 0xf1, 0xd6, 0xab, 0xab, 0x20, 0x65, 0x5c, 0x8e, 0x68, 0xf2, 0x05, 0x7c, 0x69, 0x86,
	0xd0, 0xe0, 0x7d, 0xde, 0xe2, 0x28, 0x85, 0x63, 0xf1, 0x1e, 0x96, 0xa5, 0xc2, 0x7c, 0xc0, 0x2d,
	0xeb, 0x59, 0x5d, 0x5d, 0xa6, 0xd6, 0x91, 0xd4, 0x39, 0x92, 0x7e, 0x75, 0x8e, 0x48, 0x57, 0x9a,
	0x7c, 0x82, 0xc0, 0x0a, 0x38, 0xeb, 0x71, 0x6e, 0xae, 0x77, 0x9c, 0x9b, 0x7c, 0x86, 0xf0, 0x9a,
	0x45, 0x62, 0x65, 0xfe, 0xef, 0x72, 0xed, 0x7c, 0xe5, 0xd8, 0x3c, 0xbc, 0xe9, 0x7e, 0xa0, 0x72,
	0xae, 0x72, 0x62, 0xd0, 0x87, 0x9e, 0x16, 0xc8, 0xad, 0x08, 0xe5, 0x24, 0xf9, 0x3d, 0x03, 0xff,
	0xb6, 0xae, 0xaa, 0xb3, 0xc4, 0x3c, 0x83, 0xa8, 0xec, 0xda, 0xb6, 0x1e, 0x36, 0x54, 0x6a, 0x15,
	0x85, 0x16, 0xb8, 0x63, 0x42, 0x9f, 0x0f, 0x3b, 0x76, 0x9b, 0x08, 0x26, 0x16, 0x6b, 0x80, 0xc2,
	0x28, 0xdd, 0x28, 0xac, 0x74, 0xbc, 0x60, 0x87, 0xa3, 0xd4, 0x89, 0x97, 0x51, 0x31, 0x46, 0xda,
	0x2c, 0x71, 0x8b, 0x0d, 0x0e, 0x18, 0x07, 0xc4, 0x0f, 0xe5, 0x98, 0x99, 0xae, 0xba, 0xfe, 0x89,
	0xf1, 0x92, 0x45, 0x73, 0x2c, 0xde, 0xc2, 0x23, 0xbb, 0xe2, 0xcd, 0xb8, 0xba, 0x90, 0x1b, 0x5f,
	0xa4, 0x93, 0xc5, 0xca, 0x8b, 0xe2, 0x98, 0xe8, 0xe4, 0xd7, 0x1c, 0x82, 0x1b, 0x56, 0x7a, 0xd6,
	0x43, 0xff, 0x99, 0xe0, 0x9d, 0x9a, 0x60, 0x8e, 0x81, 0x16, 0xad, 0xcc, 0x31, 0xf8, 0xa7, 0x8f,
	0x61, 0x2c, 0x15, 0x1f, 0x21, 0xac, 0xea, 0x7d, 0xad, 0x77, 0x44, 0x5b, 0x9c, 0xa4, 0xfd, 0xad,
	0x15, 0xcf, 0x69, 0x13, 0x34, 0x1c, 0x9b, 0x86, 0x88, 0xd6, 0xb1, 0x23, 0x40, 0x47, 0x0d, 0xc4,
	0x3e, 0xe0, 0xde, 0x20, 0x64, 0x9d, 0x47, 0x0f, 0x9b, 0x20, 0x45, 0xc0, 0xbd, 0xdf, 0xfd, 0x09,
	0x00, 0x00, 0xff, 0xff, 0x31, 0xad, 0xe0, 0x48, 0xe6, 0x03, 0x00, 0x00,
}
