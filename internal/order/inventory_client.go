package order

import (
	"context"
	"gopher-express/api/proto/inventory"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func GetInventoryStatus(productID string) (bool, error) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		return false, err
	}
	defer conn.Close()

	client := inventory.NewInventoryServiceClient(conn)

	md := metadata.Pairs("authorization", "Bearer my-secret-token")
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	resp, err := client.CheckStock(ctx, &inventory.StockRequest{ProductId: productID})
	if err != nil {
		return false, err
	}

	return resp.IsAvailable, nil
}