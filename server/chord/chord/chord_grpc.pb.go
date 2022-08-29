// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.15.8
// source: server/proto/chord.proto

package chord

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// ChordClient is the client API for Chord service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ChordClient interface {
	// GetPredecessor devuelve el nodo que se cree que es el predecesor.
	GetPredecessor(ctx context.Context, in *GetPredecessorRequest, opts ...grpc.CallOption) (*GetPredecessorResponse, error)
	// GetSuccessor devuelve el nodo que se cree que es el sucesor.
	GetSuccessor(ctx context.Context, in *GetSuccessorRequest, opts ...grpc.CallOption) (*GetSuccessorResponse, error)
	// SetPredecessor encuentra el predecesor para el nodo actual.
	SetPredecessor(ctx context.Context, in *SetPredecessorRequest, opts ...grpc.CallOption) (*SetPredecessorResponse, error)
	// SetSuccessor  encuentra el sucesor para el nodo actual.
	SetSuccessor(ctx context.Context, in *SetSuccessorRequest, opts ...grpc.CallOption) (*SetSuccessorResponse, error)
	// FindSuccessor encuentra el nodo que sucede a esta ID.
	FindSuccessor(ctx context.Context, in *FindSuccesorRequest, opts ...grpc.CallOption) (*FindSuccesorResponse, error)
	// Notify notifica al nodo que puede tener un nuevo sucesor.
	Notify(ctx context.Context, in *NotifyRequest, opts ...grpc.CallOption) (*NotifyResponse, error)
	// Check comprueba si el nodo esta vivo.
	Check(ctx context.Context, in *CheckRequest, opts ...grpc.CallOption) (*CheckResponse, error)
	//
	//Manda un fichero hacia el sistema y estos son guardados con
	//las etiquetas contenidas en tag-list
	AddFile(ctx context.Context, in *AddFileRequest, opts ...grpc.CallOption) (*AddFileResponse, error)
	//
	//Manda a eliminar todos los archivos que cumplen con una determinada query
	DeleteFileByQuery(ctx context.Context, in *DeleteFileByQueryRequest, opts ...grpc.CallOption) (*DeleteFileByQueryResponse, error)
	//
	//Manda a eliminar un archivo
	DeleteFile(ctx context.Context, in *DeleteFileRequest, opts ...grpc.CallOption) (*DeleteFileResponse, error)
	//
	//Manda a guardar una de las etiquetas de un archivo y la informacion relevante de dicho archivo
	AddTag(ctx context.Context, in *AddTagRequest, opts ...grpc.CallOption) (*AddTagResponse, error)
	GetTag(ctx context.Context, in *GetTagRequest, opts ...grpc.CallOption) (Chord_GetTagClient, error)
	DeleteTag(ctx context.Context, in *DeleteTagRequest, opts ...grpc.CallOption) (*DeleteTagResponse, error)
	// Get recupera el valor asociado a la clave.
	Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error)
	// Set almacena un par <clave, valor>.
	Set(ctx context.Context, in *SetRequest, opts ...grpc.CallOption) (*SetResponse, error)
	// Delete elimina el par <clave, valor> especifico.
	Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteResponse, error)
	// Partition devuelve todos los pares <clave, valor> en el almacenamiento dado  un internal.
	Partition(ctx context.Context, in *PartitionRequest, opts ...grpc.CallOption) (*PartitionResponse, error)
	// Extend establece una lista de pares <clave, valor> en el diccionario de almacenamiento
	Extend(ctx context.Context, in *ExtendRequest, opts ...grpc.CallOption) (*ExtendResponse, error)
	// Discard elimina todos los pares <clave, valor> en el almacenamiento.
	Discard(ctx context.Context, in *DiscardRequest, opts ...grpc.CallOption) (*DiscardResponse, error)
}

type chordClient struct {
	cc grpc.ClientConnInterface
}

func NewChordClient(cc grpc.ClientConnInterface) ChordClient {
	return &chordClient{cc}
}

func (c *chordClient) GetPredecessor(ctx context.Context, in *GetPredecessorRequest, opts ...grpc.CallOption) (*GetPredecessorResponse, error) {
	out := new(GetPredecessorResponse)
	err := c.cc.Invoke(ctx, "/chord.Chord/GetPredecessor", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chordClient) GetSuccessor(ctx context.Context, in *GetSuccessorRequest, opts ...grpc.CallOption) (*GetSuccessorResponse, error) {
	out := new(GetSuccessorResponse)
	err := c.cc.Invoke(ctx, "/chord.Chord/GetSuccessor", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chordClient) SetPredecessor(ctx context.Context, in *SetPredecessorRequest, opts ...grpc.CallOption) (*SetPredecessorResponse, error) {
	out := new(SetPredecessorResponse)
	err := c.cc.Invoke(ctx, "/chord.Chord/SetPredecessor", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chordClient) SetSuccessor(ctx context.Context, in *SetSuccessorRequest, opts ...grpc.CallOption) (*SetSuccessorResponse, error) {
	out := new(SetSuccessorResponse)
	err := c.cc.Invoke(ctx, "/chord.Chord/SetSuccessor", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chordClient) FindSuccessor(ctx context.Context, in *FindSuccesorRequest, opts ...grpc.CallOption) (*FindSuccesorResponse, error) {
	out := new(FindSuccesorResponse)
	err := c.cc.Invoke(ctx, "/chord.Chord/FindSuccessor", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chordClient) Notify(ctx context.Context, in *NotifyRequest, opts ...grpc.CallOption) (*NotifyResponse, error) {
	out := new(NotifyResponse)
	err := c.cc.Invoke(ctx, "/chord.Chord/Notify", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chordClient) Check(ctx context.Context, in *CheckRequest, opts ...grpc.CallOption) (*CheckResponse, error) {
	out := new(CheckResponse)
	err := c.cc.Invoke(ctx, "/chord.Chord/Check", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chordClient) AddFile(ctx context.Context, in *AddFileRequest, opts ...grpc.CallOption) (*AddFileResponse, error) {
	out := new(AddFileResponse)
	err := c.cc.Invoke(ctx, "/chord.Chord/AddFile", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chordClient) DeleteFileByQuery(ctx context.Context, in *DeleteFileByQueryRequest, opts ...grpc.CallOption) (*DeleteFileByQueryResponse, error) {
	out := new(DeleteFileByQueryResponse)
	err := c.cc.Invoke(ctx, "/chord.Chord/DeleteFileByQuery", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chordClient) DeleteFile(ctx context.Context, in *DeleteFileRequest, opts ...grpc.CallOption) (*DeleteFileResponse, error) {
	out := new(DeleteFileResponse)
	err := c.cc.Invoke(ctx, "/chord.Chord/DeleteFile", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chordClient) AddTag(ctx context.Context, in *AddTagRequest, opts ...grpc.CallOption) (*AddTagResponse, error) {
	out := new(AddTagResponse)
	err := c.cc.Invoke(ctx, "/chord.Chord/AddTag", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chordClient) GetTag(ctx context.Context, in *GetTagRequest, opts ...grpc.CallOption) (Chord_GetTagClient, error) {
	stream, err := c.cc.NewStream(ctx, &Chord_ServiceDesc.Streams[0], "/chord.Chord/GetTag", opts...)
	if err != nil {
		return nil, err
	}
	x := &chordGetTagClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Chord_GetTagClient interface {
	Recv() (*GetTagResponse, error)
	grpc.ClientStream
}

type chordGetTagClient struct {
	grpc.ClientStream
}

func (x *chordGetTagClient) Recv() (*GetTagResponse, error) {
	m := new(GetTagResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *chordClient) DeleteTag(ctx context.Context, in *DeleteTagRequest, opts ...grpc.CallOption) (*DeleteTagResponse, error) {
	out := new(DeleteTagResponse)
	err := c.cc.Invoke(ctx, "/chord.Chord/DeleteTag", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chordClient) Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error) {
	out := new(GetResponse)
	err := c.cc.Invoke(ctx, "/chord.Chord/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chordClient) Set(ctx context.Context, in *SetRequest, opts ...grpc.CallOption) (*SetResponse, error) {
	out := new(SetResponse)
	err := c.cc.Invoke(ctx, "/chord.Chord/Set", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chordClient) Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteResponse, error) {
	out := new(DeleteResponse)
	err := c.cc.Invoke(ctx, "/chord.Chord/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chordClient) Partition(ctx context.Context, in *PartitionRequest, opts ...grpc.CallOption) (*PartitionResponse, error) {
	out := new(PartitionResponse)
	err := c.cc.Invoke(ctx, "/chord.Chord/Partition", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chordClient) Extend(ctx context.Context, in *ExtendRequest, opts ...grpc.CallOption) (*ExtendResponse, error) {
	out := new(ExtendResponse)
	err := c.cc.Invoke(ctx, "/chord.Chord/Extend", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chordClient) Discard(ctx context.Context, in *DiscardRequest, opts ...grpc.CallOption) (*DiscardResponse, error) {
	out := new(DiscardResponse)
	err := c.cc.Invoke(ctx, "/chord.Chord/Discard", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ChordServer is the server API for Chord service.
// All implementations must embed UnimplementedChordServer
// for forward compatibility
type ChordServer interface {
	// GetPredecessor devuelve el nodo que se cree que es el predecesor.
	GetPredecessor(context.Context, *GetPredecessorRequest) (*GetPredecessorResponse, error)
	// GetSuccessor devuelve el nodo que se cree que es el sucesor.
	GetSuccessor(context.Context, *GetSuccessorRequest) (*GetSuccessorResponse, error)
	// SetPredecessor encuentra el predecesor para el nodo actual.
	SetPredecessor(context.Context, *SetPredecessorRequest) (*SetPredecessorResponse, error)
	// SetSuccessor  encuentra el sucesor para el nodo actual.
	SetSuccessor(context.Context, *SetSuccessorRequest) (*SetSuccessorResponse, error)
	// FindSuccessor encuentra el nodo que sucede a esta ID.
	FindSuccessor(context.Context, *FindSuccesorRequest) (*FindSuccesorResponse, error)
	// Notify notifica al nodo que puede tener un nuevo sucesor.
	Notify(context.Context, *NotifyRequest) (*NotifyResponse, error)
	// Check comprueba si el nodo esta vivo.
	Check(context.Context, *CheckRequest) (*CheckResponse, error)
	//
	//Manda un fichero hacia el sistema y estos son guardados con
	//las etiquetas contenidas en tag-list
	AddFile(context.Context, *AddFileRequest) (*AddFileResponse, error)
	//
	//Manda a eliminar todos los archivos que cumplen con una determinada query
	DeleteFileByQuery(context.Context, *DeleteFileByQueryRequest) (*DeleteFileByQueryResponse, error)
	//
	//Manda a eliminar un archivo
	DeleteFile(context.Context, *DeleteFileRequest) (*DeleteFileResponse, error)
	//
	//Manda a guardar una de las etiquetas de un archivo y la informacion relevante de dicho archivo
	AddTag(context.Context, *AddTagRequest) (*AddTagResponse, error)
	GetTag(*GetTagRequest, Chord_GetTagServer) error
	DeleteTag(context.Context, *DeleteTagRequest) (*DeleteTagResponse, error)
	// Get recupera el valor asociado a la clave.
	Get(context.Context, *GetRequest) (*GetResponse, error)
	// Set almacena un par <clave, valor>.
	Set(context.Context, *SetRequest) (*SetResponse, error)
	// Delete elimina el par <clave, valor> especifico.
	Delete(context.Context, *DeleteRequest) (*DeleteResponse, error)
	// Partition devuelve todos los pares <clave, valor> en el almacenamiento dado  un internal.
	Partition(context.Context, *PartitionRequest) (*PartitionResponse, error)
	// Extend establece una lista de pares <clave, valor> en el diccionario de almacenamiento
	Extend(context.Context, *ExtendRequest) (*ExtendResponse, error)
	// Discard elimina todos los pares <clave, valor> en el almacenamiento.
	Discard(context.Context, *DiscardRequest) (*DiscardResponse, error)
	mustEmbedUnimplementedChordServer()
}

// UnimplementedChordServer must be embedded to have forward compatible implementations.
type UnimplementedChordServer struct {
}

func (UnimplementedChordServer) GetPredecessor(context.Context, *GetPredecessorRequest) (*GetPredecessorResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPredecessor not implemented")
}
func (UnimplementedChordServer) GetSuccessor(context.Context, *GetSuccessorRequest) (*GetSuccessorResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSuccessor not implemented")
}
func (UnimplementedChordServer) SetPredecessor(context.Context, *SetPredecessorRequest) (*SetPredecessorResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetPredecessor not implemented")
}
func (UnimplementedChordServer) SetSuccessor(context.Context, *SetSuccessorRequest) (*SetSuccessorResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetSuccessor not implemented")
}
func (UnimplementedChordServer) FindSuccessor(context.Context, *FindSuccesorRequest) (*FindSuccesorResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindSuccessor not implemented")
}
func (UnimplementedChordServer) Notify(context.Context, *NotifyRequest) (*NotifyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Notify not implemented")
}
func (UnimplementedChordServer) Check(context.Context, *CheckRequest) (*CheckResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Check not implemented")
}
func (UnimplementedChordServer) AddFile(context.Context, *AddFileRequest) (*AddFileResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddFile not implemented")
}
func (UnimplementedChordServer) DeleteFileByQuery(context.Context, *DeleteFileByQueryRequest) (*DeleteFileByQueryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteFileByQuery not implemented")
}
func (UnimplementedChordServer) DeleteFile(context.Context, *DeleteFileRequest) (*DeleteFileResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteFile not implemented")
}
func (UnimplementedChordServer) AddTag(context.Context, *AddTagRequest) (*AddTagResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddTag not implemented")
}
func (UnimplementedChordServer) GetTag(*GetTagRequest, Chord_GetTagServer) error {
	return status.Errorf(codes.Unimplemented, "method GetTag not implemented")
}
func (UnimplementedChordServer) DeleteTag(context.Context, *DeleteTagRequest) (*DeleteTagResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteTag not implemented")
}
func (UnimplementedChordServer) Get(context.Context, *GetRequest) (*GetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedChordServer) Set(context.Context, *SetRequest) (*SetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Set not implemented")
}
func (UnimplementedChordServer) Delete(context.Context, *DeleteRequest) (*DeleteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedChordServer) Partition(context.Context, *PartitionRequest) (*PartitionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Partition not implemented")
}
func (UnimplementedChordServer) Extend(context.Context, *ExtendRequest) (*ExtendResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Extend not implemented")
}
func (UnimplementedChordServer) Discard(context.Context, *DiscardRequest) (*DiscardResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Discard not implemented")
}
func (UnimplementedChordServer) mustEmbedUnimplementedChordServer() {}

// UnsafeChordServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ChordServer will
// result in compilation errors.
type UnsafeChordServer interface {
	mustEmbedUnimplementedChordServer()
}

func RegisterChordServer(s grpc.ServiceRegistrar, srv ChordServer) {
	s.RegisterService(&Chord_ServiceDesc, srv)
}

func _Chord_GetPredecessor_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPredecessorRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChordServer).GetPredecessor(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chord.Chord/GetPredecessor",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChordServer).GetPredecessor(ctx, req.(*GetPredecessorRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Chord_GetSuccessor_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetSuccessorRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChordServer).GetSuccessor(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chord.Chord/GetSuccessor",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChordServer).GetSuccessor(ctx, req.(*GetSuccessorRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Chord_SetPredecessor_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetPredecessorRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChordServer).SetPredecessor(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chord.Chord/SetPredecessor",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChordServer).SetPredecessor(ctx, req.(*SetPredecessorRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Chord_SetSuccessor_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetSuccessorRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChordServer).SetSuccessor(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chord.Chord/SetSuccessor",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChordServer).SetSuccessor(ctx, req.(*SetSuccessorRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Chord_FindSuccessor_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FindSuccesorRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChordServer).FindSuccessor(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chord.Chord/FindSuccessor",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChordServer).FindSuccessor(ctx, req.(*FindSuccesorRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Chord_Notify_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NotifyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChordServer).Notify(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chord.Chord/Notify",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChordServer).Notify(ctx, req.(*NotifyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Chord_Check_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChordServer).Check(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chord.Chord/Check",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChordServer).Check(ctx, req.(*CheckRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Chord_AddFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddFileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChordServer).AddFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chord.Chord/AddFile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChordServer).AddFile(ctx, req.(*AddFileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Chord_DeleteFileByQuery_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteFileByQueryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChordServer).DeleteFileByQuery(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chord.Chord/DeleteFileByQuery",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChordServer).DeleteFileByQuery(ctx, req.(*DeleteFileByQueryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Chord_DeleteFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteFileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChordServer).DeleteFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chord.Chord/DeleteFile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChordServer).DeleteFile(ctx, req.(*DeleteFileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Chord_AddTag_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddTagRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChordServer).AddTag(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chord.Chord/AddTag",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChordServer).AddTag(ctx, req.(*AddTagRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Chord_GetTag_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(GetTagRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ChordServer).GetTag(m, &chordGetTagServer{stream})
}

type Chord_GetTagServer interface {
	Send(*GetTagResponse) error
	grpc.ServerStream
}

type chordGetTagServer struct {
	grpc.ServerStream
}

func (x *chordGetTagServer) Send(m *GetTagResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _Chord_DeleteTag_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteTagRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChordServer).DeleteTag(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chord.Chord/DeleteTag",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChordServer).DeleteTag(ctx, req.(*DeleteTagRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Chord_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChordServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chord.Chord/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChordServer).Get(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Chord_Set_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChordServer).Set(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chord.Chord/Set",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChordServer).Set(ctx, req.(*SetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Chord_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChordServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chord.Chord/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChordServer).Delete(ctx, req.(*DeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Chord_Partition_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PartitionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChordServer).Partition(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chord.Chord/Partition",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChordServer).Partition(ctx, req.(*PartitionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Chord_Extend_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExtendRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChordServer).Extend(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chord.Chord/Extend",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChordServer).Extend(ctx, req.(*ExtendRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Chord_Discard_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DiscardRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChordServer).Discard(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chord.Chord/Discard",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChordServer).Discard(ctx, req.(*DiscardRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Chord_ServiceDesc is the grpc.ServiceDesc for Chord service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Chord_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "chord.Chord",
	HandlerType: (*ChordServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetPredecessor",
			Handler:    _Chord_GetPredecessor_Handler,
		},
		{
			MethodName: "GetSuccessor",
			Handler:    _Chord_GetSuccessor_Handler,
		},
		{
			MethodName: "SetPredecessor",
			Handler:    _Chord_SetPredecessor_Handler,
		},
		{
			MethodName: "SetSuccessor",
			Handler:    _Chord_SetSuccessor_Handler,
		},
		{
			MethodName: "FindSuccessor",
			Handler:    _Chord_FindSuccessor_Handler,
		},
		{
			MethodName: "Notify",
			Handler:    _Chord_Notify_Handler,
		},
		{
			MethodName: "Check",
			Handler:    _Chord_Check_Handler,
		},
		{
			MethodName: "AddFile",
			Handler:    _Chord_AddFile_Handler,
		},
		{
			MethodName: "DeleteFileByQuery",
			Handler:    _Chord_DeleteFileByQuery_Handler,
		},
		{
			MethodName: "DeleteFile",
			Handler:    _Chord_DeleteFile_Handler,
		},
		{
			MethodName: "AddTag",
			Handler:    _Chord_AddTag_Handler,
		},
		{
			MethodName: "DeleteTag",
			Handler:    _Chord_DeleteTag_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _Chord_Get_Handler,
		},
		{
			MethodName: "Set",
			Handler:    _Chord_Set_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _Chord_Delete_Handler,
		},
		{
			MethodName: "Partition",
			Handler:    _Chord_Partition_Handler,
		},
		{
			MethodName: "Extend",
			Handler:    _Chord_Extend_Handler,
		},
		{
			MethodName: "Discard",
			Handler:    _Chord_Discard_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetTag",
			Handler:       _Chord_GetTag_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "server/proto/chord.proto",
}
