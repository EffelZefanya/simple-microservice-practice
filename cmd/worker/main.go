package main

import (
	"encoding/json"
	"fmt"
	"gopher-express/pkg/events"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, _ := amqp.Dial("amqp://guest:guest@localhost:5672/")
	ch, _ := conn.Channel()
	msgs, _ := ch.Consume("orders_queue", "", false, false, false, false, nil)

	fmt.Println("📧 Notification Worker waiting for orders...")

	for d := range msgs {
		var event events.OrderCreatedEvent
		err := json.Unmarshal(d.Body, &event)
    
    if err != nil {
        log.Printf("Error decoding: %v", err)
        d.Nack(false, false) 
        continue
    }

		fmt.Println("-------------------------------------------")
		fmt.Printf("🔔 NOTIFICATION: New Order Received!\n")
		fmt.Printf("📦 Order ID: %s\n", event.OrderID)
		fmt.Printf("👤 Customer: %s\n", event.CustomerID)
		fmt.Printf("💻 Item: %s x %d\n", event.ProductID, event.Quantity)
		fmt.Println("-------------------------------------------")

		d.Ack(false)
	}
}