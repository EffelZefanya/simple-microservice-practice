package main

import (
	"log"
	"net"
	"google.golang.org/grpc"
	"gopher-express/api/proto/inventory"
	inv "gopher-express/internal/inventory"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(
        grpc.UnaryInterceptor(inv.AuthInterceptor),
    )
	inventory.RegisterInventoryServiceServer(s, &inv.Server{})

	log.Println("Inventory gRPC Server running on :50051...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}