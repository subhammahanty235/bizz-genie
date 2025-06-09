package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/subhammahanty235/bizzgenie/internal/redisclient"
	pb "github.com/subhammahanty235/bizzgenie/proto"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedHandShakeServer
	pb.UnimplementedHeartBeatServer
	Redis *RedisClients
}

type RedisClients struct {
	internalInstance *redisclient.InternalRedisClient
	externalInstance *redisclient.ExternalRedisClient
}

func (r *RedisClients) AddHeartBeatSignal(ctx context.Context, req *pb.RequestHeartBeat) {
	key := fmt.Sprintf("heartbeat:%s", req.InstanceId)
	internalRedisClient := r.internalInstance.GetInternalRedisClient()
	currentTime := time.Now().Format(time.RFC3339)
	err := internalRedisClient.Set(ctx, key, currentTime, 0)
	if err != nil {
		log.Println("error occured while stroing hertbeat data")
	}
}

func (s *Server) HeartBeatExchange(ctx context.Context, req *pb.RequestHeartBeat) (*pb.ReplyHeartBeat, error) {
	//heartbeat store mechanism
	s.Redis.AddHeartBeatSignal(ctx, req)

	//Whenever a new heartbeat signal will arrive it will override the time
	log.Printf("Recieved Heartbeat signal from %v with signal id as %v", req.InstanceId, req.SignalId)
	return &pb.ReplyHeartBeat{InstanceId: req.InstanceId, SignalId: req.SignalId}, nil
}

func main() {

	fmt.Println("Bizz Genie is Running and Ready to process")

	//connect redis instances
	instance, err := redisclient.NewInternalRedisClient()
	if err != nil {
		log.Fatalf("failed to start Internal redis client: %v", err)
	}

	exinstance, err1 := redisclient.NewExternalRedisClient()
	if err1 != nil {
		log.Fatalf("failed to start external redis client: %v", err)
	}

	//tcp connection
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen to tcp port: %v", err)
	}

	redisHandler := &RedisClients{
		internalInstance: instance,
		externalInstance: exinstance,
	}

	server := &Server{
		Redis: redisHandler,
	}
	grpcserver := grpc.NewServer()
	pb.RegisterHandShakeServer(grpcserver, server)
	pb.RegisterHeartBeatServer(grpcserver, server)

	//handshake mechanism --> times based

	//health checkup --> timer based on regular intervals

	//redis instance health update --> while start

	//start the tcp connection
	if err := grpcserver.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	fmt.Println("TCP Connection is Running and Ready to process")

}
