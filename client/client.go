package client

import (
	"fmt"
	"log"

	pb "github.com/alejbv/SistemaDeFicherosDistribuido/tagFileSystempb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewTagFileSystemServiceClient(conn)
	fmt.Printf("The client is: %v", client)

}
