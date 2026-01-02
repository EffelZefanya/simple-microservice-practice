package events

type OrderCreatedEvent struct {
    OrderID    string  `json:"order_id"`
    CustomerID string  `json:"customer_id"`
    ProductID  string  `json:"product_id"`
    Quantity   int     `json:"quantity"`
}