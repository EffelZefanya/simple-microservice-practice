package main

import (
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
	// 1. Initialize Repository
	repo, err := order.NewRepository("mongodb://localhost:27017")
	if err != nil {
		log.Fatal(err)
	}

	// 2. Initialize RabbitMQ (Fixed redeclaration)
	rabbit, err := platform.NewRabbitMQ("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal(err)
	}

	cache := order.NewCache("localhost:6379")

	r := gin.Default()

	// --- CREATE ORDER ---
	r.POST("/orders", func(c *gin.Context) {
		var req OrderRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		available, err := order.GetInventoryStatus(req.ProductID)
		if err != nil || !available {
			c.JSON(http.StatusConflict, gin.H{"message": "Item out of stock or service down"})
			return
		}

		newOrder := order.Order{
			CustomerID: req.CustomerID,
			ProductID:  req.ProductID,
			Quantity:   req.Quantity,
			Status:     "PLACED",
			CreatedAt:  time.Now(),
		}

		id, err := repo.CreateOrder(c.Request.Context(), newOrder)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save order"})
			return
		}

		// 3. Publish Event (Cleaned up duplicates)
		event := events.OrderCreatedEvent{
			OrderID:    id,
			CustomerID: req.CustomerID,
			ProductID:  req.ProductID,
			Quantity:   req.Quantity,
		}

		if err := rabbit.PublishOrder(c.Request.Context(), event); err != nil {
			log.Printf("Failed to publish event for order %s: %v", id, err)
		}

		c.JSON(http.StatusCreated, gin.H{
			"message":  "Order placed successfully!",
			"order_id": id,
		})
	})

	// --- GET ORDER ---
	r.GET("/orders/:id", func(c *gin.Context) {
		id := c.Param("id")
		ctx := c.Request.Context()

		cachedOrder, err := cache.GetOrder(ctx, id)
		if err == nil {
			fmt.Println("🚀 Cache Hit!")
			c.JSON(http.StatusOK, cachedOrder)
			return
		}

		log.Println("🐌 Cache Miss! Fetching from DB...")
		dbOrder, err := repo.FindByID(ctx, id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}

		_ = cache.SetOrder(ctx, *dbOrder)
		c.JSON(http.StatusOK, dbOrder)
	})

	// --- DELETE ORDER ---
	r.DELETE("/orders/:id", func(c *gin.Context) {
		id := c.Param("id")
		ctx := c.Request.Context()

		err := repo.DeleteOrder(ctx, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete from DB"})
			return
		}

		err = cache.DeleteOrder(ctx, id)
		if err != nil {
			log.Printf("⚠️ Warning: Failed to evict cache for %s\n", id)
		}

		c.JSON(http.StatusOK, gin.H{"message": "Order deleted and cache cleared"})
	})

	r.Run(":8080")
}