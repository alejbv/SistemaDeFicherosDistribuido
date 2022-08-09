// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.15.8
// source: server/proto/chord.proto

package chord

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Empty request for null parameters.
type EmptyRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *EmptyRequest) Reset() {
	*x = EmptyRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_server_proto_chord_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EmptyRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EmptyRequest) ProtoMessage() {}

func (x *EmptyRequest) ProtoReflect() protoreflect.Message {
	mi := &file_server_proto_chord_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EmptyRequest.ProtoReflect.Descriptor instead.
func (*EmptyRequest) Descriptor() ([]byte, []int) {
	return file_server_proto_chord_proto_rawDescGZIP(), []int{0}
}

// Empty response for null returns.
type EmptyResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *EmptyResponse) Reset() {
	*x = EmptyResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_server_proto_chord_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EmptyResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EmptyResponse) ProtoMessage() {}

func (x *EmptyResponse) ProtoReflect() protoreflect.Message {
	mi := &file_server_proto_chord_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EmptyResponse.ProtoReflect.Descriptor instead.
func (*EmptyResponse) Descriptor() ([]byte, []int) {
	return file_server_proto_chord_proto_rawDescGZIP(), []int{1}
}

// Identifier of a node.
type ID struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID []byte `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
}

func (x *ID) Reset() {
	*x = ID{}
	if protoimpl.UnsafeEnabled {
		mi := &file_server_proto_chord_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ID) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ID) ProtoMessage() {}

func (x *ID) ProtoReflect() protoreflect.Message {
	mi := &file_server_proto_chord_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ID.ProtoReflect.Descriptor instead.
func (*ID) Descriptor() ([]byte, []int) {
	return file_server_proto_chord_proto_rawDescGZIP(), []int{2}
}

func (x *ID) GetID() []byte {
	if x != nil {
		return x.ID
	}
	return nil
}

// Node contains an ID and an address.
type Node struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID   []byte `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	IP   string `protobuf:"bytes,2,opt,name=IP,proto3" json:"IP,omitempty"`
	Port string `protobuf:"bytes,3,opt,name=port,proto3" json:"port,omitempty"`
}

func (x *Node) Reset() {
	*x = Node{}
	if protoimpl.UnsafeEnabled {
		mi := &file_server_proto_chord_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Node) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Node) ProtoMessage() {}

func (x *Node) ProtoReflect() protoreflect.Message {
	mi := &file_server_proto_chord_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Node.ProtoReflect.Descriptor instead.
func (*Node) Descriptor() ([]byte, []int) {
	return file_server_proto_chord_proto_rawDescGZIP(), []int{3}
}

func (x *Node) GetID() []byte {
	if x != nil {
		return x.ID
	}
	return nil
}

func (x *Node) GetIP() string {
	if x != nil {
		return x.IP
	}
	return ""
}

func (x *Node) GetPort() string {
	if x != nil {
		return x.Port
	}
	return ""
}

// GetRequest contains the key of a desired value.
type GetRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key  string `protobuf:"bytes,1,opt,name=Key,proto3" json:"Key,omitempty"`
	Lock bool   `protobuf:"varint,2,opt,name=Lock,proto3" json:"Lock,omitempty"`
	IP   string `protobuf:"bytes,3,opt,name=IP,proto3" json:"IP,omitempty"`
}

func (x *GetRequest) Reset() {
	*x = GetRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_server_proto_chord_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetRequest) ProtoMessage() {}

func (x *GetRequest) ProtoReflect() protoreflect.Message {
	mi := &file_server_proto_chord_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetRequest.ProtoReflect.Descriptor instead.
func (*GetRequest) Descriptor() ([]byte, []int) {
	return file_server_proto_chord_proto_rawDescGZIP(), []int{4}
}

func (x *GetRequest) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *GetRequest) GetLock() bool {
	if x != nil {
		return x.Lock
	}
	return false
}

func (x *GetRequest) GetIP() string {
	if x != nil {
		return x.IP
	}
	return ""
}

// GetResponse contains the value of a requested key.
type GetResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value []byte `protobuf:"bytes,1,opt,name=Value,proto3" json:"Value,omitempty"`
}

func (x *GetResponse) Reset() {
	*x = GetResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_server_proto_chord_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetResponse) ProtoMessage() {}

func (x *GetResponse) ProtoReflect() protoreflect.Message {
	mi := &file_server_proto_chord_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetResponse.ProtoReflect.Descriptor instead.
func (*GetResponse) Descriptor() ([]byte, []int) {
	return file_server_proto_chord_proto_rawDescGZIP(), []int{5}
}

func (x *GetResponse) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

// SetRequest contains the <key, value> pair to set on storage.
type SetRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key     string `protobuf:"bytes,1,opt,name=Key,proto3" json:"Key,omitempty"`
	Value   []byte `protobuf:"bytes,2,opt,name=Value,proto3" json:"Value,omitempty"`
	Replica bool   `protobuf:"varint,3,opt,name=Replica,proto3" json:"Replica,omitempty"`
	Lock    bool   `protobuf:"varint,4,opt,name=Lock,proto3" json:"Lock,omitempty"`
	IP      string `protobuf:"bytes,5,opt,name=IP,proto3" json:"IP,omitempty"`
}

func (x *SetRequest) Reset() {
	*x = SetRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_server_proto_chord_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetRequest) ProtoMessage() {}

func (x *SetRequest) ProtoReflect() protoreflect.Message {
	mi := &file_server_proto_chord_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SetRequest.ProtoReflect.Descriptor instead.
func (*SetRequest) Descriptor() ([]byte, []int) {
	return file_server_proto_chord_proto_rawDescGZIP(), []int{6}
}

func (x *SetRequest) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *SetRequest) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

func (x *SetRequest) GetReplica() bool {
	if x != nil {
		return x.Replica
	}
	return false
}

func (x *SetRequest) GetLock() bool {
	if x != nil {
		return x.Lock
	}
	return false
}

func (x *SetRequest) GetIP() string {
	if x != nil {
		return x.IP
	}
	return ""
}

// DeleteRequest contains the key to eliminate.
type DeleteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key     string `protobuf:"bytes,1,opt,name=Key,proto3" json:"Key,omitempty"`
	Replica bool   `protobuf:"varint,2,opt,name=Replica,proto3" json:"Replica,omitempty"`
	Lock    bool   `protobuf:"varint,3,opt,name=Lock,proto3" json:"Lock,omitempty"`
	IP      string `protobuf:"bytes,4,opt,name=IP,proto3" json:"IP,omitempty"`
}

func (x *DeleteRequest) Reset() {
	*x = DeleteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_server_proto_chord_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteRequest) ProtoMessage() {}

func (x *DeleteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_server_proto_chord_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteRequest.ProtoReflect.Descriptor instead.
func (*DeleteRequest) Descriptor() ([]byte, []int) {
	return file_server_proto_chord_proto_rawDescGZIP(), []int{7}
}

func (x *DeleteRequest) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *DeleteRequest) GetReplica() bool {
	if x != nil {
		return x.Replica
	}
	return false
}

func (x *DeleteRequest) GetLock() bool {
	if x != nil {
		return x.Lock
	}
	return false
}

func (x *DeleteRequest) GetIP() string {
	if x != nil {
		return x.IP
	}
	return ""
}

// PartitionResponse contains the contains the <key, value> pairs to return.
type PartitionResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	In  map[string][]byte `protobuf:"bytes,1,rep,name=in,proto3" json:"in,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Out map[string][]byte `protobuf:"bytes,2,rep,name=out,proto3" json:"out,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *PartitionResponse) Reset() {
	*x = PartitionResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_server_proto_chord_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PartitionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PartitionResponse) ProtoMessage() {}

func (x *PartitionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_server_proto_chord_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PartitionResponse.ProtoReflect.Descriptor instead.
func (*PartitionResponse) Descriptor() ([]byte, []int) {
	return file_server_proto_chord_proto_rawDescGZIP(), []int{8}
}

func (x *PartitionResponse) GetIn() map[string][]byte {
	if x != nil {
		return x.In
	}
	return nil
}

func (x *PartitionResponse) GetOut() map[string][]byte {
	if x != nil {
		return x.Out
	}
	return nil
}

// ExtendRequest contains the <key, value> pairs to set on the storage.
type ExtendRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Dictionary map[string][]byte `protobuf:"bytes,1,rep,name=dictionary,proto3" json:"dictionary,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *ExtendRequest) Reset() {
	*x = ExtendRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_server_proto_chord_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ExtendRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExtendRequest) ProtoMessage() {}

func (x *ExtendRequest) ProtoReflect() protoreflect.Message {
	mi := &file_server_proto_chord_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExtendRequest.ProtoReflect.Descriptor instead.
func (*ExtendRequest) Descriptor() ([]byte, []int) {
	return file_server_proto_chord_proto_rawDescGZIP(), []int{9}
}

func (x *ExtendRequest) GetDictionary() map[string][]byte {
	if x != nil {
		return x.Dictionary
	}
	return nil
}

// DiscardRequest contains the lower and upper bound of the interval storage to delete.
type DiscardRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Keys []string `protobuf:"bytes,1,rep,name=keys,proto3" json:"keys,omitempty"`
}

func (x *DiscardRequest) Reset() {
	*x = DiscardRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_server_proto_chord_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DiscardRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DiscardRequest) ProtoMessage() {}

func (x *DiscardRequest) ProtoReflect() protoreflect.Message {
	mi := &file_server_proto_chord_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DiscardRequest.ProtoReflect.Descriptor instead.
func (*DiscardRequest) Descriptor() ([]byte, []int) {
	return file_server_proto_chord_proto_rawDescGZIP(), []int{10}
}

func (x *DiscardRequest) GetKeys() []string {
	if x != nil {
		return x.Keys
	}
	return nil
}

var File_server_proto_chord_proto protoreflect.FileDescriptor

var file_server_proto_chord_proto_rawDesc = []byte{
	0x0a, 0x18, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63,
	0x68, 0x6f, 0x72, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x63, 0x68, 0x6f, 0x72,
	0x64, 0x22, 0x0e, 0x0a, 0x0c, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x22, 0x0f, 0x0a, 0x0d, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x14, 0x0a, 0x02, 0x49, 0x44, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x02, 0x49, 0x44, 0x22, 0x3a, 0x0a, 0x04, 0x4e, 0x6f, 0x64, 0x65,
	0x12, 0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x02, 0x49, 0x44,
	0x12, 0x0e, 0x0a, 0x02, 0x49, 0x50, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x49, 0x50,
	0x12, 0x12, 0x0a, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x70, 0x6f, 0x72, 0x74, 0x22, 0x42, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x4b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x4b, 0x65, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x4c, 0x6f, 0x63, 0x6b, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x04, 0x4c, 0x6f, 0x63, 0x6b, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x50, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x49, 0x50, 0x22, 0x23, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x56, 0x61, 0x6c, 0x75, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x72, 0x0a,
	0x0a, 0x53, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x4b,
	0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x4b, 0x65, 0x79, 0x12, 0x14, 0x0a,
	0x05, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x56, 0x61,
	0x6c, 0x75, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x12, 0x12, 0x0a,
	0x04, 0x4c, 0x6f, 0x63, 0x6b, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x04, 0x4c, 0x6f, 0x63,
	0x6b, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x50, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x49,
	0x50, 0x22, 0x5f, 0x0a, 0x0d, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x4b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x4b, 0x65, 0x79, 0x12, 0x18, 0x0a, 0x07, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x12, 0x12,
	0x0a, 0x04, 0x4c, 0x6f, 0x63, 0x6b, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x04, 0x4c, 0x6f,
	0x63, 0x6b, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x50, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02,
	0x49, 0x50, 0x22, 0xe9, 0x01, 0x0a, 0x11, 0x50, 0x61, 0x72, 0x74, 0x69, 0x74, 0x69, 0x6f, 0x6e,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x30, 0x0a, 0x02, 0x69, 0x6e, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x63, 0x68, 0x6f, 0x72, 0x64, 0x2e, 0x50, 0x61, 0x72,
	0x74, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x49,
	0x6e, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x02, 0x69, 0x6e, 0x12, 0x33, 0x0a, 0x03, 0x6f, 0x75,
	0x74, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x63, 0x68, 0x6f, 0x72, 0x64, 0x2e,
	0x50, 0x61, 0x72, 0x74, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x2e, 0x4f, 0x75, 0x74, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x03, 0x6f, 0x75, 0x74, 0x1a,
	0x35, 0x0a, 0x07, 0x49, 0x6e, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65,
	0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0x36, 0x0a, 0x08, 0x4f, 0x75, 0x74, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x94,
	0x01, 0x0a, 0x0d, 0x45, 0x78, 0x74, 0x65, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x44, 0x0a, 0x0a, 0x64, 0x69, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x61, 0x72, 0x79, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x63, 0x68, 0x6f, 0x72, 0x64, 0x2e, 0x45, 0x78, 0x74,
	0x65, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x44, 0x69, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x61, 0x72, 0x79, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x0a, 0x64, 0x69, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x61, 0x72, 0x79, 0x1a, 0x3d, 0x0a, 0x0f, 0x44, 0x69, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x61, 0x72, 0x79, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x24, 0x0a, 0x0e, 0x44, 0x69, 0x73, 0x63, 0x61, 0x72, 0x64,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6b, 0x65, 0x79, 0x73, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x6b, 0x65, 0x79, 0x73, 0x32, 0x9d, 0x05, 0x0a, 0x05,
	0x43, 0x68, 0x6f, 0x72, 0x64, 0x12, 0x32, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x50, 0x72, 0x65, 0x64,
	0x65, 0x63, 0x65, 0x73, 0x73, 0x6f, 0x72, 0x12, 0x13, 0x2e, 0x63, 0x68, 0x6f, 0x72, 0x64, 0x2e,
	0x45, 0x6d, 0x70, 0x74, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0b, 0x2e, 0x63,
	0x68, 0x6f, 0x72, 0x64, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x12, 0x30, 0x0a, 0x0c, 0x47, 0x65, 0x74,
	0x53, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x6f, 0x72, 0x12, 0x13, 0x2e, 0x63, 0x68, 0x6f, 0x72,
	0x64, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0b,
	0x2e, 0x63, 0x68, 0x6f, 0x72, 0x64, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x12, 0x33, 0x0a, 0x0e, 0x53,
	0x65, 0x74, 0x50, 0x72, 0x65, 0x64, 0x65, 0x63, 0x65, 0x73, 0x73, 0x6f, 0x72, 0x12, 0x0b, 0x2e,
	0x63, 0x68, 0x6f, 0x72, 0x64, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x1a, 0x14, 0x2e, 0x63, 0x68, 0x6f,
	0x72, 0x64, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x31, 0x0a, 0x0c, 0x53, 0x65, 0x74, 0x53, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x6f, 0x72,
	0x12, 0x0b, 0x2e, 0x63, 0x68, 0x6f, 0x72, 0x64, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x1a, 0x14, 0x2e,
	0x63, 0x68, 0x6f, 0x72, 0x64, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x27, 0x0a, 0x0d, 0x46, 0x69, 0x6e, 0x64, 0x53, 0x75, 0x63, 0x63, 0x65,
	0x73, 0x73, 0x6f, 0x72, 0x12, 0x09, 0x2e, 0x63, 0x68, 0x6f, 0x72, 0x64, 0x2e, 0x49, 0x44, 0x1a,
	0x0b, 0x2e, 0x63, 0x68, 0x6f, 0x72, 0x64, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x12, 0x2b, 0x0a, 0x06,
	0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x12, 0x0b, 0x2e, 0x63, 0x68, 0x6f, 0x72, 0x64, 0x2e, 0x4e,
	0x6f, 0x64, 0x65, 0x1a, 0x14, 0x2e, 0x63, 0x68, 0x6f, 0x72, 0x64, 0x2e, 0x45, 0x6d, 0x70, 0x74,
	0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x32, 0x0a, 0x05, 0x43, 0x68, 0x65,
	0x63, 0x6b, 0x12, 0x13, 0x2e, 0x63, 0x68, 0x6f, 0x72, 0x64, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e, 0x63, 0x68, 0x6f, 0x72, 0x64, 0x2e,
	0x45, 0x6d, 0x70, 0x74, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2c, 0x0a,
	0x03, 0x47, 0x65, 0x74, 0x12, 0x11, 0x2e, 0x63, 0x68, 0x6f, 0x72, 0x64, 0x2e, 0x47, 0x65, 0x74,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x12, 0x2e, 0x63, 0x68, 0x6f, 0x72, 0x64, 0x2e,
	0x47, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2e, 0x0a, 0x03, 0x53,
	0x65, 0x74, 0x12, 0x11, 0x2e, 0x63, 0x68, 0x6f, 0x72, 0x64, 0x2e, 0x53, 0x65, 0x74, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e, 0x63, 0x68, 0x6f, 0x72, 0x64, 0x2e, 0x45, 0x6d,
	0x70, 0x74, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x34, 0x0a, 0x06, 0x44,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x12, 0x14, 0x2e, 0x63, 0x68, 0x6f, 0x72, 0x64, 0x2e, 0x44, 0x65,
	0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e, 0x63, 0x68,
	0x6f, 0x72, 0x64, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x3a, 0x0a, 0x09, 0x50, 0x61, 0x72, 0x74, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x13,
	0x2e, 0x63, 0x68, 0x6f, 0x72, 0x64, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x63, 0x68, 0x6f, 0x72, 0x64, 0x2e, 0x50, 0x61, 0x72, 0x74,
	0x69, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x34, 0x0a,
	0x06, 0x45, 0x78, 0x74, 0x65, 0x6e, 0x64, 0x12, 0x14, 0x2e, 0x63, 0x68, 0x6f, 0x72, 0x64, 0x2e,
	0x45, 0x78, 0x74, 0x65, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e,
	0x63, 0x68, 0x6f, 0x72, 0x64, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x36, 0x0a, 0x07, 0x44, 0x69, 0x73, 0x63, 0x61, 0x72, 0x64, 0x12, 0x15,
	0x2e, 0x63, 0x68, 0x6f, 0x72, 0x64, 0x2e, 0x44, 0x69, 0x73, 0x63, 0x61, 0x72, 0x64, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e, 0x63, 0x68, 0x6f, 0x72, 0x64, 0x2e, 0x45, 0x6d,
	0x70, 0x74, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x16, 0x5a, 0x14, 0x2e,
	0x2f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2f, 0x63, 0x68, 0x6f, 0x72, 0x64, 0x2f, 0x63, 0x68,
	0x6f, 0x72, 0x64, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_server_proto_chord_proto_rawDescOnce sync.Once
	file_server_proto_chord_proto_rawDescData = file_server_proto_chord_proto_rawDesc
)

func file_server_proto_chord_proto_rawDescGZIP() []byte {
	file_server_proto_chord_proto_rawDescOnce.Do(func() {
		file_server_proto_chord_proto_rawDescData = protoimpl.X.CompressGZIP(file_server_proto_chord_proto_rawDescData)
	})
	return file_server_proto_chord_proto_rawDescData
}

var file_server_proto_chord_proto_msgTypes = make([]protoimpl.MessageInfo, 14)
var file_server_proto_chord_proto_goTypes = []interface{}{
	(*EmptyRequest)(nil),      // 0: chord.EmptyRequest
	(*EmptyResponse)(nil),     // 1: chord.EmptyResponse
	(*ID)(nil),                // 2: chord.ID
	(*Node)(nil),              // 3: chord.Node
	(*GetRequest)(nil),        // 4: chord.GetRequest
	(*GetResponse)(nil),       // 5: chord.GetResponse
	(*SetRequest)(nil),        // 6: chord.SetRequest
	(*DeleteRequest)(nil),     // 7: chord.DeleteRequest
	(*PartitionResponse)(nil), // 8: chord.PartitionResponse
	(*ExtendRequest)(nil),     // 9: chord.ExtendRequest
	(*DiscardRequest)(nil),    // 10: chord.DiscardRequest
	nil,                       // 11: chord.PartitionResponse.InEntry
	nil,                       // 12: chord.PartitionResponse.OutEntry
	nil,                       // 13: chord.ExtendRequest.DictionaryEntry
}
var file_server_proto_chord_proto_depIdxs = []int32{
	11, // 0: chord.PartitionResponse.in:type_name -> chord.PartitionResponse.InEntry
	12, // 1: chord.PartitionResponse.out:type_name -> chord.PartitionResponse.OutEntry
	13, // 2: chord.ExtendRequest.dictionary:type_name -> chord.ExtendRequest.DictionaryEntry
	0,  // 3: chord.Chord.GetPredecessor:input_type -> chord.EmptyRequest
	0,  // 4: chord.Chord.GetSuccessor:input_type -> chord.EmptyRequest
	3,  // 5: chord.Chord.SetPredecessor:input_type -> chord.Node
	3,  // 6: chord.Chord.SetSuccessor:input_type -> chord.Node
	2,  // 7: chord.Chord.FindSuccessor:input_type -> chord.ID
	3,  // 8: chord.Chord.Notify:input_type -> chord.Node
	0,  // 9: chord.Chord.Check:input_type -> chord.EmptyRequest
	4,  // 10: chord.Chord.Get:input_type -> chord.GetRequest
	6,  // 11: chord.Chord.Set:input_type -> chord.SetRequest
	7,  // 12: chord.Chord.Delete:input_type -> chord.DeleteRequest
	0,  // 13: chord.Chord.Partition:input_type -> chord.EmptyRequest
	9,  // 14: chord.Chord.Extend:input_type -> chord.ExtendRequest
	10, // 15: chord.Chord.Discard:input_type -> chord.DiscardRequest
	3,  // 16: chord.Chord.GetPredecessor:output_type -> chord.Node
	3,  // 17: chord.Chord.GetSuccessor:output_type -> chord.Node
	1,  // 18: chord.Chord.SetPredecessor:output_type -> chord.EmptyResponse
	1,  // 19: chord.Chord.SetSuccessor:output_type -> chord.EmptyResponse
	3,  // 20: chord.Chord.FindSuccessor:output_type -> chord.Node
	1,  // 21: chord.Chord.Notify:output_type -> chord.EmptyResponse
	1,  // 22: chord.Chord.Check:output_type -> chord.EmptyResponse
	5,  // 23: chord.Chord.Get:output_type -> chord.GetResponse
	1,  // 24: chord.Chord.Set:output_type -> chord.EmptyResponse
	1,  // 25: chord.Chord.Delete:output_type -> chord.EmptyResponse
	8,  // 26: chord.Chord.Partition:output_type -> chord.PartitionResponse
	1,  // 27: chord.Chord.Extend:output_type -> chord.EmptyResponse
	1,  // 28: chord.Chord.Discard:output_type -> chord.EmptyResponse
	16, // [16:29] is the sub-list for method output_type
	3,  // [3:16] is the sub-list for method input_type
	3,  // [3:3] is the sub-list for extension type_name
	3,  // [3:3] is the sub-list for extension extendee
	0,  // [0:3] is the sub-list for field type_name
}

func init() { file_server_proto_chord_proto_init() }
func file_server_proto_chord_proto_init() {
	if File_server_proto_chord_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_server_proto_chord_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EmptyRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_server_proto_chord_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EmptyResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_server_proto_chord_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ID); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_server_proto_chord_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Node); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_server_proto_chord_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_server_proto_chord_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_server_proto_chord_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SetRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_server_proto_chord_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_server_proto_chord_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PartitionResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_server_proto_chord_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ExtendRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_server_proto_chord_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DiscardRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_server_proto_chord_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   14,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_server_proto_chord_proto_goTypes,
		DependencyIndexes: file_server_proto_chord_proto_depIdxs,
		MessageInfos:      file_server_proto_chord_proto_msgTypes,
	}.Build()
	File_server_proto_chord_proto = out.File
	file_server_proto_chord_proto_rawDesc = nil
	file_server_proto_chord_proto_goTypes = nil
	file_server_proto_chord_proto_depIdxs = nil
}
