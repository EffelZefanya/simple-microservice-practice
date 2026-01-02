package order

import (
	"context"
	"encoding/json"
	"time"
	"github.com/go-redis/redis/v8"
)

type Cache struct {
	rdb *redis.Client
}

func NewCache(addr string) *Cache {
	return &Cache{
		rdb: redis.NewClient(&redis.Options{Addr: addr}),
	}
}

func (c *Cache) SetOrder(ctx context.Context, order Order) error {
	data, _ := json.Marshal(order)
	return c.rdb.Set(ctx, "order:"+order.ID.Hex(), data, 10*time.Minute).Err()
}

func (c *Cache) GetOrder(ctx context.Context, id string) (*Order, error) {
	val, err := c.rdb.Get(ctx, "order:"+id).Result()
	if err != nil {
		return nil, err
	}

	var order Order
	err = json.Unmarshal([]byte(val), &order)
	return &order, err
}

func (c *Cache) DeleteOrder(ctx context.Context, id string) error {
	return c.rdb.Del(ctx, "order:"+id).Err()
}