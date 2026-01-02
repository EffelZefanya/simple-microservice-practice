package order

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	db *mongo.Collection
}

func NewRepository(uri string) (*Repository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	col := client.Database("gopher_express").Collection("orders")
	return &Repository{db: col}, nil
}

func (r *Repository) CreateOrder(ctx context.Context, o Order) (string, error) {
	result, err := r.db.InsertOne(ctx, o)
	if err != nil {
		return "", err
	}
	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (r *Repository) FindByID(ctx context.Context, idStr string) (*Order, error) {
	objID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return nil, err
	}

	var result Order
	err = r.db.FindOne(ctx, bson.M{"_id": objID}).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *Repository) DeleteOrder(ctx context.Context, idStr string) error {
	objID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return err
	}

	_, err = r.db.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}