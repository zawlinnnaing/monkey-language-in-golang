# Monkey Language in Golang

Implement your own interpretor in Golang.

## Installation & Setup

### Prerequisites
- Go 1.19 or higher
- Git

### Setup
1. Install dependencies:
   ```bash
   go mod tidy
   ```

2. Run tests to verify setup:
   ```bash
   go test ./...
   ```

### Running the REPL
To start the Monkey language REPL (Read-Eval-Print Loop):
```bash
go run main.go
```

### Language Specification
The Monkey language specification and examples can be found in the test files throughout the project. These tests serve as both documentation and validation of the language features.


## TODOs
- [ ] Add support for `<=` and `>=` infix operators
- [ ] Add stacktrace on errors
- [ ] Add support for optional function parameters