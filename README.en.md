- [English](README.en.md)
- [Русский](README.ru.md)


# Vault Copy



Vault Copy is a secure secret management system written in Go. It provides a robust API for managing applications, users, and secrets with a focus on security and audit logging.

## Features

- User authentication and registration
- Application management
- Secret storage and retrieval
- Policy-based access control
- Audit logging
- Cryptographic operations

## Project Structure

The project is organized into several packages:

- `db_operations`: Database connection and models
- `services`: Core services including user, app, and logging
- `internal`: Internal API for policy and application management
- `cryptoOperation`: Cryptographic functions

## Getting Started

### Prerequisites

- Go 1.x
- PostgreSQL database

### Installation

1. Clone the repository
2. Set up the PostgreSQL database and update the connection details in `db_operations/db.go`
3. Run `go mod tidy` to install dependencies
4. Build the project with `go build`

### Running the Server

Execute the compiled binary or run:

The server will start on port 8080 by default.

## API Endpoints

- `/`: Hello World endpoint (for testing)

More endpoints to be documented as they are implemented.

## Security Features

- Password hashing with salts
- SHA256 for various hashing operations
- Policy-based access control for applications

## Logging

The project includes both server logging and audit logging for tracking operations and maintaining security records.

## Contributing

Contributions are welcome. Please ensure you follow the existing code structure and add appropriate error handling and logging.

## TODO

- Complete security checks in various functions
- Implement more API endpoints
- Enhance error handling
- Add more comprehensive documentation for each package

## License

[Add your chosen license here]

## Disclaimer

This project is a work in progress and may not be suitable for production use without further security audits and enhancements.