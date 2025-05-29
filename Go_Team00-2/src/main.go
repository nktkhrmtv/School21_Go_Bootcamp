package main

import (
	"log"
	"math/rand"
	"net"
	"time"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"team00/message" 
)

type server struct {
	message.UnimplementedMessageServiceServer
}

func (s *server) StreamMessages(req *message.StreamMessagesRequest, stream message.MessageService_StreamMessagesServer) error {
	mean := rand.Float64()*20 - 10     
	stdDev := rand.Float64()*1.2 + 0.3 
	sessionID := uuid.New().String()

	log.Printf("New session: %s, Mean: %.2f, StdDev: %.2f", sessionID, mean, stdDev)

	for {
		frequency := rand.NormFloat64()*stdDev + mean
		timestamp := time.Now().UTC().Unix()
		if err := stream.Send(&message.Message{
			SessionId: sessionID,
			Frequency: frequency,
			Timestamp: timestamp,
		}); err != nil {
			return err
		}

		time.Sleep(time.Millisecond * 100)
	}
}

func main() {
	lis, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	message.RegisterMessageServiceServer(s, &server{})
	reflection.Register(s)
	log.Printf("Server is running on port 8888")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
