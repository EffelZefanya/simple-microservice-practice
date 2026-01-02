package order

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CustomerID string             `bson:"customer_id" json:"customer_id"`
	ProductID  string             `bson:"product_id" json:"product_id"`
	Quantity   int                `bson:"quantity" json:"quantity"`
	Status     string             `bson:"status" json:"status"` // e.g., "PENDING", "COMPLETED"
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
}