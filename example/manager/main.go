package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/tuanuet/bwork"
	pb "github.com/tuanuet/bwork/proto"
	"google.golang.org/grpc"
)

func main() {
	sm := bwork.NewShardManager()
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterShardServiceServer(s, sm)
	log.Printf("ShardManager gRPC server listening at %v", lis.Addr())

	// Simulate some activity
	go func() {
		time.Sleep(5 * time.Second)
		for i := 1; i <= 10; i++ {
			recordID := fmt.Sprintf("record%d", i)
			sm.AddRecord(recordID)
		}
	}()

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
