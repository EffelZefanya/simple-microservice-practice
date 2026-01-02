package main

import (
	"context"
	"fmt"
	"gopher-express/internal/order"
	"gopher-express/internal/platform"
	"gopher-express/pkg/events"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type OrderRequest struct {
	ProductID  string `json:"product_id" binding:"required"`
	CustomerID string `json:"customer_id" binding:"required"`
	Quantity   int    `json:"quantity" binding:"required"`
}

func main() {
	repo, err := order.NewRepository("mongodb://localhost:27017")
	if err != nil {
		log.Fatal(err)
	}

	rabbit, _ := platform.NewRabbitMQ("amqp://guest:guest@localhost:5672/")

	cache := order.NewCache("localhost:6379")

	r := gin.Default()

	r.POST("/orders", func(c *gin.Context) {
		var req OrderRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 1. Call gRPC Inventory Service
		available, err := order.GetInventoryStatus(req.ProductID)
		if err != nil || !available {
			c.JSON(http.StatusConflict, gin.H{"message": "Item out of stock or service down"})
			return
		}

		// 2. Logic for MongoDB and RabbitMQ will go here next...
		
		newOrder := order.Order{
			CustomerID: req.CustomerID,
			ProductID:  req.ProductID,
			Quantity:   req.Quantity,
			Status:     "PLACED",
			CreatedAt:  time.Now(),
		}

		id, err := repo.CreateOrder(context.Background(), newOrder)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save order"})
			return
		}

		event := events.OrderCreatedEvent{
    OrderID:    id,
    CustomerID: req.CustomerID,
    ProductID:  req.ProductID,
    Quantity:   req.Quantity,
}
_ = rabbit.PublishOrder(context.Background(), event)

		c.JSON(http.StatusCreated, gin.H{
			"message":  "Order placed successfully!",
			"order_id": id,
		})
	})

	r.GET("/orders/:id", func(c *gin.Context) {
    id := c.Param("id")
    ctx := context.Background()

    // 1. Try to get from Redis
    cachedOrder, err := cache.GetOrder(ctx, id)
    if err == nil {
        fmt.Println("🚀 Cache Hit!")
        c.JSON(http.StatusOK, cachedOrder)
        return
    }

    // 2. Cache Miss - Get from MongoDB
    fmt.Println("🐌 Cache Miss! Fetching from DB...")
    dbOrder, err := repo.FindByID(ctx, id) // You'll need to add FindByID to your repo
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
        return
    }

    // 3. Store in Redis for next time
    _ = cache.SetOrder(ctx, *dbOrder)

    c.JSON(http.StatusOK, dbOrder)
})

r.DELETE("/orders/:id", func(c *gin.Context) {
    id := c.Param("id")
    ctx := context.Background()

    // 1. Delete from MongoDB
    err := repo.DeleteOrder(ctx, id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete from DB"})
        return
    }

    // 2. IMPORTANT: Evict from Redis Cache
    err = cache.DeleteOrder(ctx, id)
    if err != nil {
        // We log this but don't necessarily fail the request, 
        // though in high-stakes apps, this is a critical sync point.
        fmt.Printf("⚠️ Warning: Failed to evict cache for %s\n", id)
    }

    c.JSON(http.StatusOK, gin.H{"message": "Order deleted and cache cleared"})
})

	r.Run(":8080") // Listen on port 8080
}