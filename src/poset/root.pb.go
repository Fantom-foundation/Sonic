// Code generated by protoc-gen-go. DO NOT EDIT.
// source: root.proto

package poset

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type RootEvent struct {
	Hash                 []byte   `protobuf:"bytes,1,opt,name=Hash,proto3" json:"Hash,omitempty"`
	CreatorID            uint64   `protobuf:"varint,2,opt,name=CreatorID,proto3" json:"CreatorID,omitempty"`
	Index                int64    `protobuf:"varint,3,opt,name=Index,proto3" json:"Index,omitempty"`
	LamportTimestamp     int64    `protobuf:"varint,4,opt,name=LamportTimestamp,proto3" json:"LamportTimestamp,omitempty"`
	Round                int64    `protobuf:"varint,5,opt,name=Round,proto3" json:"Round,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RootEvent) Reset()         { *m = RootEvent{} }
func (m *RootEvent) String() string { return proto.CompactTextString(m) }
func (*RootEvent) ProtoMessage()    {}
func (*RootEvent) Descriptor() ([]byte, []int) {
	return fileDescriptor_08a043f6ee9336a8, []int{0}
}

func (m *RootEvent) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RootEvent.Unmarshal(m, b)
}
func (m *RootEvent) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RootEvent.Marshal(b, m, deterministic)
}
func (m *RootEvent) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RootEvent.Merge(m, src)
}
func (m *RootEvent) XXX_Size() int {
	return xxx_messageInfo_RootEvent.Size(m)
}
func (m *RootEvent) XXX_DiscardUnknown() {
	xxx_messageInfo_RootEvent.DiscardUnknown(m)
}

var xxx_messageInfo_RootEvent proto.InternalMessageInfo

func (m *RootEvent) GetHash() []byte {
	if m != nil {
		return m.Hash
	}
	return nil
}

func (m *RootEvent) GetCreatorID() uint64 {
	if m != nil {
		return m.CreatorID
	}
	return 0
}

func (m *RootEvent) GetIndex() int64 {
	if m != nil {
		return m.Index
	}
	return 0
}

func (m *RootEvent) GetLamportTimestamp() int64 {
	if m != nil {
		return m.LamportTimestamp
	}
	return 0
}

func (m *RootEvent) GetRound() int64 {
	if m != nil {
		return m.Round
	}
	return 0
}

type Root struct {
	NextRound            int64                 `protobuf:"varint,1,opt,name=NextRound,proto3" json:"NextRound,omitempty"`
	SelfParent           *RootEvent            `protobuf:"bytes,2,opt,name=SelfParent,proto3" json:"SelfParent,omitempty"`
	Others               map[string]*RootEvent `protobuf:"bytes,3,rep,name=Others,proto3" json:"Others,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *Root) Reset()         { *m = Root{} }
func (m *Root) String() string { return proto.CompactTextString(m) }
func (*Root) ProtoMessage()    {}
func (*Root) Descriptor() ([]byte, []int) {
	return fileDescriptor_08a043f6ee9336a8, []int{1}
}

func (m *Root) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Root.Unmarshal(m, b)
}
func (m *Root) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Root.Marshal(b, m, deterministic)
}
func (m *Root) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Root.Merge(m, src)
}
func (m *Root) XXX_Size() int {
	return xxx_messageInfo_Root.Size(m)
}
func (m *Root) XXX_DiscardUnknown() {
	xxx_messageInfo_Root.DiscardUnknown(m)
}

var xxx_messageInfo_Root proto.InternalMessageInfo

func (m *Root) GetNextRound() int64 {
	if m != nil {
		return m.NextRound
	}
	return 0
}

func (m *Root) GetSelfParent() *RootEvent {
	if m != nil {
		return m.SelfParent
	}
	return nil
}

func (m *Root) GetOthers() map[string]*RootEvent {
	if m != nil {
		return m.Others
	}
	return nil
}

func init() {
	proto.RegisterType((*RootEvent)(nil), "poset.RootEvent")
	proto.RegisterType((*Root)(nil), "poset.Root")
	proto.RegisterMapType((map[string]*RootEvent)(nil), "poset.Root.OthersEntry")
}

func init() { proto.RegisterFile("root.proto", fileDescriptor_08a043f6ee9336a8) }

var fileDescriptor_08a043f6ee9336a8 = []byte{
	// 267 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x90, 0xc1, 0x4a, 0xfb, 0x40,
	0x10, 0xc6, 0xd9, 0x26, 0x29, 0x64, 0xf2, 0x3f, 0x84, 0xe1, 0x0f, 0x2e, 0xe2, 0x21, 0xf4, 0x20,
	0xc1, 0x43, 0x94, 0x7a, 0x11, 0xaf, 0x5a, 0xb0, 0x28, 0x2a, 0xab, 0x2f, 0xb0, 0xd2, 0x91, 0x8a,
	0x4d, 0x26, 0x6c, 0xa6, 0xa5, 0x7d, 0x10, 0xdf, 0xcc, 0x07, 0x92, 0xec, 0x8a, 0x0d, 0x88, 0xb7,
	0x99, 0x6f, 0x7e, 0xcc, 0x7c, 0xdf, 0x00, 0x38, 0x66, 0xa9, 0x5a, 0xc7, 0xc2, 0x98, 0xb4, 0xdc,
	0x91, 0x4c, 0x3e, 0x14, 0xa4, 0x86, 0x59, 0x66, 0x1b, 0x6a, 0x04, 0x11, 0xe2, 0x1b, 0xdb, 0x2d,
	0xb5, 0x2a, 0x54, 0xf9, 0xcf, 0xf8, 0x1a, 0x8f, 0x20, 0xbd, 0x72, 0x64, 0x85, 0xdd, 0xfc, 0x5a,
	0x8f, 0x0a, 0x55, 0xc6, 0x66, 0x2f, 0xe0, 0x7f, 0x48, 0xe6, 0xcd, 0x82, 0xb6, 0x3a, 0x2a, 0x54,
	0x19, 0x99, 0xd0, 0xe0, 0x09, 0xe4, 0x77, 0xb6, 0x6e, 0xd9, 0xc9, 0xf3, 0x5b, 0x4d, 0x9d, 0xd8,
	0xba, 0xd5, 0xb1, 0x07, 0x7e, 0xe9, 0xfd, 0x06, 0xc3, 0xeb, 0x66, 0xa1, 0x93, 0xb0, 0xc1, 0x37,
	0x93, 0x4f, 0x05, 0x71, 0xef, 0xab, 0x3f, 0x7f, 0x4f, 0x5b, 0x09, 0x88, 0xf2, 0xc8, 0x5e, 0xc0,
	0x33, 0x80, 0x27, 0x5a, 0xbd, 0x3e, 0x5a, 0x47, 0x8d, 0x78, 0x77, 0xd9, 0x34, 0xaf, 0x7c, 0xb4,
	0xea, 0x27, 0x96, 0x19, 0x30, 0x78, 0x0a, 0xe3, 0x07, 0x59, 0x92, 0xeb, 0x74, 0x54, 0x44, 0x65,
	0x36, 0x3d, 0x18, 0xd0, 0x55, 0x98, 0xcc, 0x1a, 0x71, 0x3b, 0xf3, 0x8d, 0x1d, 0xde, 0x42, 0x36,
	0x90, 0x31, 0x87, 0xe8, 0x9d, 0x76, 0xde, 0x49, 0x6a, 0xfa, 0x12, 0x8f, 0x21, 0xd9, 0xd8, 0xd5,
	0x9a, 0xfe, 0x3c, 0x1f, 0xc6, 0x97, 0xa3, 0x0b, 0xf5, 0x32, 0xf6, 0xcf, 0x3f, 0xff, 0x0a, 0x00,
	0x00, 0xff, 0xff, 0x6e, 0xb3, 0x69, 0x45, 0x8a, 0x01, 0x00, 0x00,
}