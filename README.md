# go-clean-starter

This is a clean architecture starter template for Golang applications, designed with modularity and scalability in mind.

## Features
- Implements **Modular Clean Architecture** to ensure separation of concerns.
- Uses **dig (Dependency Injection)** for managing dependencies efficiently.
- Employs **mux** as the HTTP router.
- Integrates **PostgreSQL** as the database.
- Uses **logrus** for structured logging.

## Project Structure
```
./app        # Application entry point
./core       # Core components shared across modules
./internal   # Internal modules for the application
```
Each module follows a structured approach with the following subdirectories:
```
- domain         # Entity definitions and business rules
- usecase        # Business logic and application use cases
- infrastructure # Database, external APIs, etc.
- delivery       # HTTP handlers and middleware
```

## Available Modules
- **Auth Module**: Handles authentication and token verification.
- **User Module**: Manages user creation and deletion.

## Running the Application

### Start the application:
```sh
make up
```

### Makefile Commands:

#### Development Environment
```sh
make up             # Start the application (Docker + Air for live reload)
make down           # Stop Docker
make destroy        # Remove Docker containers, volumes, and clean temp files
```

#### Docker Management
```sh
make dev-env       # Start Docker environment with PostgreSQL
make dev-env-test  # Build and run the application inside Docker
make docker-stop   # Stop running Docker containers
make docker-teardown # Remove Docker containers and volumes
```

This setup allows for a clean, structured, and scalable Go application. 🚀

