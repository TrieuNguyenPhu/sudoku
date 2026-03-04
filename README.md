# Sudoku Game Website

A web-based Sudoku game built with Go and the Gin framework, featuring puzzle generation at multiple difficulty levels, an automatic solver, and an interactive web interface.

## Features

- **Puzzle Generation**: Generate Sudoku puzzles at three difficulty levels (Easy, Medium, Hard)
- **Automatic Solver**: Solve any valid Sudoku puzzle
- **Interactive Web Interface**: Play Sudoku directly in your browser with server-side rendering
- **Clean Architecture**: Well-structured codebase following Go best practices
- **Graceful Shutdown**: Proper handling of termination signals

## Project Structure

```
sudoku-game-website/
├── cmd/
│   └── server/
│       └── main.go                 # Application entry point
├── internal/
│   ├── handlers/                   # HTTP request handlers
│   ├── services/                   # Business logic services
│   ├── models/                     # Data models and structures
│   ├── middleware/                 # HTTP middleware components
│   └── config/                     # Application configuration
├── web/
│   ├── templates/                  # HTML templates
│   └── static/                     # CSS and JavaScript assets
├── go.mod                          # Go module definition
├── go.sum                          # Go module checksums
└── README.md                       # This file
```

## Technology Stack

- **Backend**: Go 1.21+ with Gin web framework
- **Templates**: Go's built-in html/template package
- **Frontend**: HTML, CSS, and vanilla JavaScript
- **Architecture**: Clean architecture with separation of concerns

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Git

### Installation

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd sudoku-game-website
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Run the application:
   ```bash
   go run cmd/server/main.go
   ```

4. Open your browser and navigate to `http://localhost:8080`

### Configuration

The application can be configured using environment variables:

- `PORT`: Server port (default: 8080)
- `GIN_MODE`: Gin mode - debug, release, or test (default: debug)

Example:
```bash
PORT=3000 GIN_MODE=release go run cmd/server/main.go
```

## API Endpoints

### Puzzle Generation
- `GET /api/easy` - Generate an easy Sudoku puzzle (40-45 pre-filled cells)
- `GET /api/medium` - Generate a medium Sudoku puzzle (30-35 pre-filled cells)
- `GET /api/hard` - Generate a hard Sudoku puzzle (20-25 pre-filled cells)

### Solver
- `POST /api/sudoku-solver` - Solve a Sudoku puzzle

### Web Interface
- `GET /` - Main game interface
- `GET /static/*` - Static assets (CSS, JavaScript)

## Development Status

This project is currently under development. The basic project structure and entry point have been set up with graceful shutdown handling. The following components are planned for implementation:

1. ✅ Project structure and core dependencies
2. ⏳ Core data models and validation
3. ⏳ Sudoku solver service
4. ⏳ Puzzle generation service
5. ⏳ HTTP handlers and API endpoints
6. ⏳ Web interface templates and assets
7. ⏳ Comprehensive testing suite

## Requirements

The application fulfills the following key requirements:

- **Graceful Shutdown**: Implements proper handling of termination signals (SIGINT, SIGTERM)
- **Performance**: Designed to handle concurrent requests and meet response time requirements
- **Error Handling**: Comprehensive error handling with appropriate HTTP status codes
- **Clean Architecture**: Modular design with clear separation of concerns

## Contributing

This project follows Go best practices and clean architecture principles. When contributing:

1. Follow the existing project structure
2. Write tests for new functionality
3. Ensure proper error handling
4. Document public APIs
5. Use meaningful commit messages

## License

[License information to be added]