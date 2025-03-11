## DEMO OF CLEAN ARCHITECTURE USING grpc
![CleanGrpc](https://github.com/user-attachments/assets/c4a5ea5d-0f76-4ace-ba53-48d0ac1ca359)
A clean architecture implementation of a gRPC service in Go, demonstrating separation of concerns and testability.

## Project Overview

CleanGrpc is a user management service built using Go and gRPC, following the clean architecture principles. The project demonstrates how to structure a Go application with clear separation of concerns, making it maintainable, testable, and scalable.

### Architecture

The project follows a clean architecture pattern with three distinct layers:

1. **Data Layer (Repository)** - Handles data persistence and retrieval
   - Interacts directly with the database (GORM)
   - Implements the repository interface defined in the domain layer
   - Located in `pkg/v1/Repository`

2. **Domain Layer (Use Case)** - Contains business logic
   - Defines the core business rules and logic
   - Depends on abstractions (interfaces) rather than concrete implementations
   - Located in `pkg/v1/UseCase`

3. **Handler Layer** - Handles external communication
   - Implements the gRPC service interface
   - Transforms data between the domain model and the gRPC protocol buffers
   - Located in `pkg/v1/handler/grpc`

### Key Features

- Complete separation of concerns following clean architecture principles
- gRPC API for efficient client-server communication
- Comprehensive test suite for all layers
- SQLite database integration using GORM
- Command-line client for interacting with the service

## Getting Started

### Prerequisites

- Go 1.24
- Git

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yishak-cs/CleanGrpc.git
   cd CleanGrpc
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Build the project:
   ```bash
   go build -o cleangrpc ./cmd/server
   go build -o client ./cmd/client
   ```

## Running the Application

### Start the Server

```bash
./cleangrpc
# Or directly with Go
go run cmd/server/main.go
```

The server will start on port 50000 by default.

### Using the Client

The client provides a command-line interface to interact with the service:

```bash
# Create a new user
./client create "John Doe" "john@example.com"

# Get a user by ID
./client get 1

# List all users
./client list

# Update a user
./client update 1 "John Updated" "john.updated@example.com"

# Delete a user
./client delete 1
```

## Testing

The project includes comprehensive tests for all layers of the architecture. The tests for the handler and use case layers were developed with assistance from Claude AI.

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests for a specific package
go test ./pkg/v1/Repository/test
go test ./pkg/v1/UseCase/test
go test ./pkg/v1/handler/grpc/test
```

### Test Structure 

- **Repository Tests**: Unit tests that verify the repository layer's interaction with the database using an in-memory SQLite database.
- **Use Case Tests**: Integration tests that verify the business logic using mocked repositories.
- **Handler Tests**: End-to-end tests that verify the gRPC handler using mocked use cases.

## Project Structure

```
CleanGrpc/
├── cmd/
│   ├── client/         # gRPC client implementation
│   └── server/         # Main application entry point
├── Internal/
│   └── model/          # Domain models
├── pkg/
│   └── v1/
│       ├── handler/    # gRPC handlers
│       ├── Repository/ # Data access layer
│       └── UseCase/    # Business logic layer
├── proto/              # Protocol buffer definitions
├── go.mod              # Go module definition
└── README.md           # Project documentation
```


