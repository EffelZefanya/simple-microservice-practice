# GopherExpress: Distributed Order Orchestrator

A high-performance microservices-style backend built in Go. This project demonstrates event-driven architecture, inter-service communication via gRPC, and multi-layer data persistence.

## 🏗️ Architecture Overview
- **API Gateway (REST):** Built with Gin, handles customer requests and coordinates services.
- **Inventory Service (gRPC):** A protected internal service for real-time stock validation.
- **Message Broker (RabbitMQ):** Decouples the order logic from the notification system.
- **Caching (Redis):** Implements the Cache-Aside pattern to reduce DB load.
- **Database (MongoDB):** Stores orders as flexible JSON-like documents.



## 🛠️ Tech Stack
- **Language:** Go (Golang)
- **API:** Gin Gonic
- **RPC:** gRPC with Protobuf & Interceptor Auth
- **Messaging:** RabbitMQ (AMQP 0.0.1)
- **Primary DB:** MongoDB
- **Caching:** Redis

## 🚀 Key Features
- **gRPC Authentication:** Internal calls are secured via Metadata Interceptors.
- **Async Notifications:** Orders trigger background events to RabbitMQ; processed by a separate Worker.
- **Cache Eviction:** Deleting an order automatically invalidates the Redis entry to prevent stale data.
- **Observability:** Centralized `/health` endpoint monitoring all infrastructure connections.

## 🚦 Getting Started
1. **Clone & Install:** `go mod download`
2. **Start Infrastructure:** `docker-compose up -d`
3. **Set Environment:** Create a `.env` file (see `.env.example`)
4. **Run Inventory:** `go run cmd/inventory/main.go`
5. **Run Worker:** `go run cmd/worker/main.go`
6. **Run API:** `go run cmd/api/main.go`

## 📈 Performance
- **Database Query:** ~8ms
- **Redis Cache Hit:** ~2ms (~75% faster)