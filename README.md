# JokeCLI

JokeCLI is a command-line interface tool written in Go that fetches jokes and funny facts from various sources across the internet. It's built using the Cobra CLI framework and provides a simple, intuitive interface for accessing different types of jokes.

## Features

Currently supported joke sources:
- Chuck Norris facts from [api.chucknorris.io](https://api.chucknorris.io/)
- Dad jokes from [icanhazdadjoke.com](https://icanhazdadjoke.com/)
- Interactive mode for selecting joke types

## Installation

### Prerequisites
- Go 1.21 or higher

### Building from source

1. Clone the repository:
```bash
git clone https://github.com/brianhays/jokecli.git
cd jokecli
```

2. Install dependencies:
```bash
go mod tidy
```

3. Build the binary:
```bash
go build
```

## Usage

### Basic Commands

- Interactive mode (select joke type):
```bash
./jokecli
```

- Get a random Chuck Norris fact:
```bash
./jokecli chuck
```

- Get a random dad joke:
```bash
./jokecli dad
```

- Display help information:
```bash
./jokecli --help
```

- Display help for a specific command:
```bash
./jokecli chuck --help
```

## Testing

The project includes both unit tests and integration tests.

### Unit Tests
Run the unit test suite:
```bash
make test-unit
```

Unit tests cover:
#### Joke Commands (Chuck Norris & Dad Jokes)
- Success cases with valid responses
- Error handling for invalid JSON responses
- Server error handling
- Missing field validation
- HTTP header validation
- Command execution validation

#### Root Command
- Help text validation
- Command registration
- Basic command structure

#### HTTP Client Mocking
- Custom test client for reliable testing
- Response mocking
- Header validation
- Status code handling

### Integration Tests
Run the integration test suite (requires internet connection):
```bash
make test-integration
```

Integration tests verify:
- End-to-end functionality with real APIs
- Rate limiting compliance
- Command-line interface behavior
- Help text and command registration

### Running All Tests
```bash
# Run all tests (unit + integration)
make test-all

# Run tests with verbose output
go test -v ./...

# Run tests and generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run tests for specific commands
go test -v ./cmd -run TestGetChuckNorrisJoke
go test -v ./cmd -run TestGetDadJoke
```

## Project Structure

```
jokecli/
├── cmd/                    # Command-line interface definitions
│   ├── root.go            # Root command and interactive mode
│   ├── chuck.go           # Chuck Norris jokes command
│   ├── dad.go             # Dad jokes command
│   ├── root_test.go       # Root command tests
│   ├── chuck_test.go      # Chuck Norris command tests
│   └── dad_test.go        # Dad jokes command tests
├── internal/              # Internal packages
│   ├── jokesapi/         # Core joke fetching logic
│   │   ├── client.go     # Shared HTTP client interface
│   │   ├── chuck.go      # Chuck Norris API interaction
│   │   └── dad.go        # Dad jokes API interaction
│   └── testutils/        # Testing utilities and HTTP mocks
├── tests/
│   └── integration/      # Integration tests
├── main.go               # Application entry point
├── Makefile             # Build and test automation
├── go.mod               # Go module file
└── README.md            # Project documentation
```

The project follows standard Go project layout:
- `cmd/`: Contains the CLI commands and their tests. Each command file focuses on command definition and user interaction.
- `internal/`: Houses packages that are internal to the project:
  - `jokesapi/`: Core logic for fetching jokes from various APIs, separated by source
  - `testutils/`: Shared testing utilities used across the project
- `tests/`: Contains integration tests that verify end-to-end functionality
- Root level files handle project configuration and documentation

## Contributing

Contributions are welcome! Feel free to submit issues and pull requests.

### Development Guidelines
1. Write tests for new features
   - Include unit tests for components
   - Add integration tests for API interactions
   - Mock external dependencies
   - Validate HTTP interactions
   - Test command execution
2. Ensure all tests pass before submitting PRs
3. Update documentation as needed
4. Follow Go best practices and coding standards
5. Add appropriate error handling
6. Include header validation for API requests

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Upcoming Features

- Support for joke categories
- Additional joke sources (programming jokes, etc.)
- Joke formatting options
- Offline mode with cached jokes
- Save favorite jokes
- Share jokes via clipboard 