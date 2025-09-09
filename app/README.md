# 🌍 IP Country Locator Service

This Go application provides a lightweight and efficient service for retrieving the country of origin for a given IP address. It is designed to be a performant and observable microservice, leveraging a caching layer to minimize external API calls.

---

## ✨ Features

* **RESTful API**: Exposes a simple API endpoint to query for a country by IP address.
* **Database Caching**: Utilizes a PostgreSQL database to cache IP-to-country mappings, which significantly reduces latency and reliance on the external API for repeated lookups.
* **External API Integration**: Integrates with the `ipinfo.io` API to fetch country data for new IP addresses.
* **Health Checks**: Includes dedicated endpoints for liveness (`/live`) and readiness (`/ready`) probes, essential for containerized environments.
* **Observability**: Exposes Prometheus metrics on a separate port (`:5050`) to monitor key performance indicators such as request latency, total requests, and in-flight requests.

---

## 🛠️ Technologies Used

* **Language**: Go (Golang).
* **Database**: PostgreSQL with `pgx` and `pgxpool` for connection pooling.
* **External API**: `ipinfo/go/v2/ipinfo` for IP lookup.
* **Metrics**: `prometheus/client_golang` for application monitoring.
* **Containerization**: Docker and Docker Compose.

---

## 📦 Getting Started

These instructions will get a copy of the project up and running on your local machine for development and testing.

### Prerequisites

* Go (version 1.24 or newer).
* Docker and Docker Compose.
* An `IPINFO_TOKEN` from ipinfo.io.

### Installation

1.  Clone the repository:
    ```bash
    git clone [https://github.com/your-username/your-repo.git](https://github.com/your-username/your-repo.git)
    cd your-repo/app
    ```
2.  Set up environment variables. Create a `.env` file from the following template and replace `YOUR_IPINFO_TOKEN` with your actual token:
    ```bash
    # Database Configuration
    DB_HOST=db
    DB_PORT=5432
    DB_USER=test
    DB_PASSWORD=test
    DB_NAME=test

    # IPInfo Token
    IPINFO_TOKEN=YOUR_IPINFO_TOKEN
    ```
3.  Build and run the application using Docker Compose. This command will set up the application and the PostgreSQL database for you.

    ```bash
    docker-compose up --build
    ```

---

## 🐳 Dockerization

The application is containerized for easy deployment and management.

### `Dockerfile`

This file uses a multi-stage build to create a lightweight, production-ready image.

```dockerfile
# Stage 1: Build the Go application
FROM golang:1.24 AS builder

WORKDIR /app

# Copy go.mod and go.sum to cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /go-app ./main.go

# Stage 2: Create the final production image
FROM alpine:3.19

# Set the working directory
WORKDIR /root/

# Copy the compiled binary from the builder stage
COPY --from=builder /go-app .

# Expose the port the app will listen on
EXPOSE 8080 5050

# Run the application
CMD ["./go-app"]
```

### `docker-compose.yaml`

This file orchestrates the application and its database dependencies, ensuring they are started and linked correctly.

```yaml
services:
  app:
    build: .
    ports:
      - "8080:8080"
      - "5050:5050"
    environment:
      - DB_HOST=db
      - DB_PORT=5432
      - DB_USER=test
      - DB_PASSWORD=test
      - DB_NAME=test
      - IPINFO_TOKEN=${IPINFO_TOKEN}
    depends_on:
      db:
        condition: service_healthy
  
  db:
    image: postgres:16-alpine
    restart: always
    environment:
      - POSTGRES_USER=test
      - POSTGRES_PASSWORD=test
      - POSTGRES_DB=test
    volumes:
      - postgres_data:/var/lib/postgresql/data/
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U $$POSTGRES_USER -d $$POSTGRES_DB"]
      interval: 5s
      timeout: 5s
      retries: 5
      start_period: 30s

volumes:
  postgres_data:
```

## 🧪 API Endpoints

The API server exposes the following endpoints:

* `GET /get_country?ip=IP_ADDRESS`

  * Description: Retrieves the country for a given IP address. It first checks the cache (PostgreSQL database) and, if not found, queries the ipinfo.io external API.

  * Example: http://localhost:8080/get_country?ip=8.8.8.8

  * Response: {"country": "US"}

* `GET /live`

  * Description: Liveness probe. Returns 200 OK with "Live!" if the application is running and the database connection is healthy.

* `GET /ready`

  * Description: Readiness probe. Returns 200 OK with "Ready!" if the application is ready to accept traffic and the database connection is healthy.

The metrics server exposes the following endpoint for Prometheus scraping:

* `GET /metrics`

  * Description: Exposes application metrics in a format Prometheus can consume.
