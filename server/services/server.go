package services

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/alejbv/SistemaDeFicherosDistribuido/server/chord"
	pb "github.com/alejbv/SistemaDeFicherosDistribuido/tagFileSystempb"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedTagFileSystemServiceServer
}

/*
Falta implemetar todos los servicos RPC
*/

// Server Stream Service

func (*Server) ListFiles(req *pb.ListFilesRequest, stream pb.TagFileSystemService_ListFilesServer) error {

	return nil
}

// Client Stream Service

func (*Server) AddFiles(stream pb.TagFileSystemService_AddFilesServer) error {
	// Implementar- Falta el chord
	return nil
}

func (*Server) AddTags(stream pb.TagFileSystemService_AddTagsServer) error {
	// Implementar- Falta el chord
	return nil
}

// Unary Service
func (*Server) DeleteFiles(ctx context.Context, req *pb.DeleteFilesRequest) (*pb.DeleteFilesResponse, error) {

	// Implementar- Falta el chord
	return &pb.DeleteFilesResponse{
		Response: "Not implemented",
	}, nil
}

func (*Server) DeleteTags(ctx context.Context, req *pb.DeleteTagsRequest) (*pb.DeleteTagsResponse, error) {

	// Implementar- Falta el chord
	return &pb.DeleteTagsResponse{
		Response: "Not implemented",
	}, nil
}

//var collection *mongo.Collection

func NewServerListening(addr string) {
	lis, err := net.Listen("tcp", addr)

	if err != nil {
		log.Fatalf("Failed to listen %v", err)
	}
	// Closing the  listener at  the end
	defer func() {
		fmt.Println("Closing the listener")
		lis.Close()
	}()

	opts := []grpc.ServerOption{}
	grpcserver := grpc.NewServer(opts...)
	pb.RegisterTagFileSystemServiceServer(grpcserver, &Server{})

	fmt.Println("Starting Server...")
	if err := grpcserver.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

	/*
		go func() {
			fmt.Println("Starting Server...")
			if err := grpcserver.Serve(lis); err != nil {
				log.Fatalf("Failed to serve: %v", err)
			}
		}()

		// Wait for Control C to exit
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt)

		 Block until a signal is received
		<-ch
		fmt.Println("Stopping the Server...")
		grpcserver.Stop()
		fmt.Println("Closing the listener")
		lis.Close()
	*/
}

var (
	node       *chord.Node
	rsaPrivate string
	rsaPublic  string
)

func StartServer(network string, rsaPrivateKeyPath string, rsaPublicteKeyPath string) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	fmt.Println("The Tag File System Service has started")

	var err error
	rsaPrivate = rsaPrivateKeyPath
	rsaPublic = rsaPublicteKeyPath

	//NewServerListening("0.0.0.0:50051")
	//fmt.Println("Ending the service")
	// Completar esto
	node, err = chord.DefaultNode("50050")
	if err != nil {
		log.Fatalf("No se pudo crear el nodo")
	}

	err = node.Start()

	if err != nil {
		log.Fatalf("No se pudo iniciar el nodo")
	}

}
