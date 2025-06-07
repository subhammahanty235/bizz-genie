package main

import (
	"fmt"
	"log"
	"net"

	pb "github.com/subhammahanty235/bizzgenie/proto"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedHandShakeServer
	pb.UnimplementedHeartBeatServer
}

// type BizzGenie struct {
// 	// internalRedis *redisclient.NewInternalRedisClient()
// 	externalRedis *redisclient.NewExternalRedisClient()
// }

func main() {

	fmt.Println("Bizz Genie is Running and Ready to process")
	//tcp connection
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen to tcp port: %v", err)
	}

	grpcserver := grpc.NewServer()
	pb.RegisterHandShakeServer(grpcserver, &Server{})
	pb.RegisterHeartBeatServer(grpcserver, &Server{})

	//handshake mechanism --> times based

	//health checkup --> timer based on regular intervals

	//redis instance health update --> while start

	//start the tcp connection
	if err := grpcserver.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
