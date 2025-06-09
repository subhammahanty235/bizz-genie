package main

import (
	"context"
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

func (s *Server) HeartBeatExchange(ctx context.Context, req *pb.RequestHeartBeat) (*pb.ReplyHeartBeat, error) {
	log.Printf("Recieved Heartbeat signal from %v with signal id as ", req.InstanceId, req.SignalId)
	return &pb.ReplyHeartBeat{InstanceId: req.InstanceId, SignalId: req.SignalId}, nil
}

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
	fmt.Println("TCP Connection is Running and Ready to process")

}
