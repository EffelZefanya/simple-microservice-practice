package inventory

import (
	"context"
	"gopher-express/api/proto/inventory" // Update this path to match your module name
)

type Server struct {
	inventory.UnimplementedInventoryServiceServer
}

// CheckStock implements the logic defined in your .proto file
func (s *Server) CheckStock(ctx context.Context, req *inventory.StockRequest) (*inventory.StockResponse, error) {
	// For now, let's hardcode a logic: 
	// If ProductID is "laptop", we have 10. Otherwise, 0.
	// Later, you will replace this with a Redis lookup.
	
	if req.ProductId == "laptop" {
		return &inventory.StockResponse{
			Quantity:    10,
			IsAvailable: true,
		}, nil
	}

	return &inventory.StockResponse{
		Quantity:    0,
		IsAvailable: false,
	}, nil
}