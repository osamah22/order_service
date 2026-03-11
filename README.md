# Order Service

A simple Order Management API written in Go.
The service provides endpoints for managing products and orders, following a clean architecture approach with clear separation between handlers, services, and models.

The API also includes Swagger documentation for easy exploration and testing.

## Features

- Clean architecture separation
- Swagger API documentation
- Docker-ready project structure

## Tech Stack

- Go
- Gin (HTTP framework)
- GORM (ORM)
- SQLite (default database)
- Swagger (Swaggo)

## Running the Project (Docker)

1. Clone the repository

```bash
git clone https://github.com/yourusername/order_service.git
cd order_service
```

2. Build the Docker image

```bash
docker build -t order-service .
```

3. Run the container

```bash
docker run -p 7070:8080 order-service
```

4. Open swagger docs on browser
   `http://localhost:7070/swagger/index.html`

---

## API Endpoints

### Products

| Method | Endpoint         | Description                  |
| ------ | ---------------- | ---------------------------- |
| POST   | `/products`      | Create a new product         |
| GET    | `/products`      | List all products            |
| GET    | `/products/{id}` | Get a specific product by ID |
| PUT    | `/products/{id}` | Update an existing product   |
| DELETE | `/products/{id}` | Delete a product by ID       |

### Orders

| Method | Endpoint       | Description                |
| ------ | -------------- | -------------------------- |
| POST   | `/orders`      | Create a new order         |
| GET    | `/orders`      | List all orders            |
| GET    | `/orders/{id}` | Get a specific order by ID |
| DELETE | `/orders/{id}` | Delete an order by ID      |

Relationships:

```
orders
└── line_items
└── product_id
```

## Folder Architecture

```
order-service/
├── cmd/
│ └── api/
│ ├── main.go
│ ├── routes.go
│ ├── products.go
│ └── orders.go
│
├── internal/
│ ├── models/
│ │ ├── product.go
│ │ └── order.go
│ │
│ ├── services/
│ │ ├── product_service.go
│ │ └── order_service.go
│ │
│ └── dtos/
│ ├── product_request.go
│ └── order_request.go
│
├── docs/
│ └── swagger files
│
├── go.mod
├── go.sum
└── README.md
```
