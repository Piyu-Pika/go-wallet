# Wallet Management REST API

A simple RESTful backend service built in Go to manage user wallets and transactions. This project allows users to create accounts, manage wallets, check balances, transfer funds, and view transaction history.

## Table of Contents

- [Features](#features)
- [Technologies](#technologies)
- [Project Structure](#project-structure)
- [Setup Instructions](#setup-instructions)
  - [Prerequisites](#prerequisites)
  - [Local Setup](#local-setup)
  - [Docker Setup](#docker-setup)
- [API Documentation](#api-documentation)
  - [Endpoints](#endpoints)
- [Running Tests](#running-tests)
- [Contributing](#contributing)
- [License](#license)

## Features

- Create a user with an associated wallet
- Create additional wallets for existing users
- Check wallet balances
- Transfer funds between wallets
- List all transactions for a wallet
- Input validation and error handling
- Database transactions for data consistency
- Unit tests with coverage reporting
- Dockerized application for easy deployment

## Technologies

- **Language**: Go (v1.20 or later)
- **Web Framework**: Gin (v1.10.0)
- **Database**: CockroachDB (PostgreSQL-compatible) for production, SQLite for testing
- **ORM**: GORM (v1.25.12)
- **Dependency Management**: Go Modules
- **Testing**: Go's testing package, testify for assertions
- **Containerization**: Docker
- **Environment Variables**: Managed via godotenv

## Project Structure

```
go-wallet/
├── cmd/
│   └── go-wallet/
│       └── main.go              # Entry point for the application
├── internal/
│   ├── database/
│   │   └── database.go          # Database initialization and operations
│   ├── handler/
│   │   └── handler.go           # HTTP handlers for API endpoints
│   ├── models/
│   │   └── models.go            # Data models (User, Wallet, Transaction)
│   └── routers/
│       └── router.go            # API route definitions
├── Dockerfile                   # Docker configuration
├── go.mod                       # Go module dependencies
├── go.sum                       # Dependency checksums
├── run_tests.sh                 # Bash script to run tests
└── README.md                    # Project documentation
```

## Setup Instructions

### Prerequisites

- **Go**: Version 1.20 or later ([Download](https://golang.org/dl/))
- **CockroachDB**: A running instance for production ([Setup Guide](https://www.cockroachlabs.com/docs/stable/install-cockroachdb.html))
- **Docker**: Optional, for containerized setup ([Install Docker](https://docs.docker.com/get-docker/))
- **Git**: For cloning the repository ([Install Git](https://git-scm.com/downloads))

### Local Setup

1. **Clone the Repository**:
   ```bash
   git clone https://github.com/Piyu-Pika/go-wallet.git
   cd go-wallet
   ```

2. **Install Dependencies**:
   ```bash
   go mod tidy
   ```

3. **Set Up Environment Variables**:
   - Create a `.env` file in the project root:
   ```bash
   echo "COCKROACH_DSN=postgresql://user:password@localhost:26257/wallet_db?sslmode=disable" > .env
   ```
   - Replace `user`, `password`, `localhost:26257`, and `wallet_db` with your CockroachDB connection details.

4. **Run the Application**:
   ```bash
   go run cmd/go-wallet/main.go
   ```
   - The API will be available at `http://localhost:8080`.

### Docker Setup

1. **Build the Docker Image**:
   ```bash
   docker build -t go-wallet .
   ```

2. **Run the Container**:
   ```bash
   docker run --rm -p 8080:8080 -e COCKROACH_DSN="postgresql://user:password@host:26257/wallet_db?sslmode=disable" go-wallet
   ```
   - Replace the `COCKROACH_DSN` value with your CockroachDB connection string.
   - Ensure your CockroachDB instance is accessible from the Docker container.

3. **Access the API**:
   - The API will be available at `http://localhost:8080`.
  
## Try at
  -  Try the deployed API at  `https://go-wallet.leapcell.app`

## API Documentation

### Endpoints

All endpoints are prefixed with `/api/`. The API expects JSON payloads for POST requests and returns JSON responses.

#### 1. Create User

- **Method**: `POST`
- **Path**: `/users`
- **Request Body**:
  ```json
  {
    "name": "John Doe",
    "email": "john@example.com"
  }
  ```
- **Response (201 Created)**:
  ```json
  {
    "ID": 1,
    "Name": "John Doe",
    "Email": "john@example.com",
    "CreatedAt": "2025-04-30T12:00:00Z",
    "UpdatedAt": "2025-04-30T12:00:00Z",
    "Wallets": [
      {
        "ID": 1,
        "UserID": 1,
        "Balance": 0,
        "Currency": "USD",
        "CreatedAt": "2025-04-30T12:00:00Z",
        "UpdatedAt": "2025-04-30T12:00:00Z"
      }
    ]
  }
  ```
- **Errors**:
  - `400 Bad Request`: Invalid input
  - `409 Conflict`: Email already exists
  - `500 Internal Server Error`: Database failure

#### 2. Get User Info

- **Method**: `GET`
- **Path**: `/users/:id`
- **Response (200 OK)**: Same as Create User response
- **Errors**:
  - `400 Bad Request`: Invalid user ID
  - `404 Not Found`: User not found

#### 3. Create Wallet

- **Method**: `POST`
- **Path**: `/wallets`
- **Request Body**:
  ```json
  {
    "user_id": 1
  }
  ```
- **Response (201 Created)**:
  ```json
  {
    "ID": 2,
    "UserID": 1,
    "Balance": 0,
    "Currency": "USD",
    "CreatedAt": "2025-04-30T12:01:00Z",
    "UpdatedAt": "2025-04-30T12:01:00Z"
  }
  ```
- **Errors**:
  - `400 Bad Request`: Invalid input
  - `404 Not Found`: User not found
  - `500 Internal Server Error`: Database failure

#### 4. Get Wallet Info

- **Method**: `GET`
- **Path**: `/wallets/:id`
- **Response (200 OK)**: Same as Create Wallet response
- **Errors**:
  - `400 Bad Request`: Invalid wallet ID
  - `404 Not Found`: Wallet not found

#### 5. Deposit 
- **Method**: `POST`
- **Path**: `/deposit`
- **Request Body**:
  ```json
  {
    "wallet_id": 1,
    "amount": 100.0
  }
  ```
- **Response (201 Created)**:
  ```json
  {
    "message": "Deposit processed"
  }
  ```
- **Errors**:
  - `400 Bad Request`: Invalid input
  - `404 Not Found`: Wallet not found
  - `500 Internal Server Error`: Database failure


#### 6. Get Wallet Balance

- **Method**: `GET`
- **Path**: `/wallets/:id/balance`
- **Response (200 OK)**:
  ```json
  {
    "balance": 100.0
  }
  ```
- **Errors**:
  - `400 Bad Request`: Invalid wallet ID
  - `404 Not Found`: Wallet not found

#### 7. Create Transaction

- **Method**: `POST`
- **Path**: `/transactions`
- **Request Body**:
  ```json
  {
    "source_wallet_id": 1,
    "dest_wallet_id": 2,
    "amount": 50.0
  }
  ```
- **Response (201 Created)**:
  ```json
  {
    "message": "Transaction processed"
  }
  ```
- **Errors**:
  - `400 Bad Request`: Invalid input
  - `404 Not Found`: Source or destination wallet not found
  - `500 Internal Server Error`: Insufficient funds or database failure

#### 8. List Transactions

- **Method**: `GET`
- **Path**: `/transactions/wallet/:wallet_id`
- **Response (200 OK)**:
  ```json
  [
    {
      "ID": 1,
      "SourceWalletID": 1,
      "DestWalletID": 2,
      "Amount": 50.0,
      "Status": "COMPLETED",
      "CreatedAt": "2025-04-30T12:02:00Z",
      "UpdatedAt": "2025-04-30T12:02:00Z"
    }
  ]
  ```
- **Errors**:
  - `400 Bad Request`: Invalid wallet ID
  - `500 Internal Server Error`: Database failure




## Running Tests

The project includes unit tests for all handlers in `internal/handler/handlers_test.go`.

1. **Install Test Dependencies**:
   ```bash
   go get gorm.io/driver/sqlite
   go mod tidy
   ```

2. **Run Tests Locally**:
   ```bash
   cd internal/handler
   go test -v
   ```

3. **Run Tests with Coverage**:
   ```bash
   go test -v -coverprofile=coverage.out
   go tool cover -html=coverage.out -o coverage.html
   ```
   - Open `coverage.html` in a browser to view coverage details.



## License

This project is licensed under the MIT License. See the LICENSE file for details.
