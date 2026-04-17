# CinemaAbyss API Tests

This directory contains Postman tests for the CinemaAbyss microservices architecture. The tests are designed to be run using Newman, the command-line collection runner for Postman.

## Structure

- `CinemaAbyss.postman_collection.json` - The main Postman collection containing all API tests
- `local.environment.json` - Environment variables for running tests against locally running services
- `docker.environment.json` - Environment variables for running tests against services in Docker containers
- `run-tests.js` - Node.js script to run the tests using Newman
- `package.json` - Node.js package configuration with dependencies and scripts

## Test Coverage

Supported test scope (implemented services in repository):

1. **Monolith Service**
   - User management (create, get)
   - Movie management (create, get)
   - Payment processing (create, get)
   - Subscription management (create, get)

2. **Movies Microservice**
   - Health check
   - Movie management (create, get)

`Events Microservice` and `Proxy Service` folders are kept in the collection as planned migration stages and are not part of CI execution.

## Prerequisites

- Node.js (v14 or later)
- npm (v6 or later)

## Installation

```bash
# Navigate to the tests directory
cd tests/postman

# Install dependencies
npm install
```

## Running Tests

### Basic Usage

```bash
# Run all tests against the local environment (default)
npm test

# Run all tests against the Docker environment
npm run test:docker
```

### Running Specific Test Folders

```bash
# Run only Monolith Service tests
npm run test:docker -- --folder "Monolith Service"

# Run only Movies Microservice tests
npm run test:docker -- --folder "Movies Microservice"

# Optional planned-stage tests
npm run test:docker -- --folder "Events Microservice"
npm run test:docker -- --folder "Proxy Service"
```

### Advanced Usage

The `run-tests.js` script supports several command-line options:

```bash
node run-tests.js --environment <env> --folder <folder> --reporters <reporters> --bail --timeout <ms>
```

Options:
- `--environment`, `-e`: Environment to run tests against (default: 'local')
- `--collection`, `-c`: Collection to run (default: 'CinemaAbyss')
- `--folder`, `-f`: Specific folder in the collection to run
- `--reporters`, `-r`: Reporters to use, comma-separated (default: 'cli,htmlextra,junit')
- `--bail`, `-b`: Stop on first error (default: false)
- `--timeout`, `-t`: Request timeout in ms (default: 10000)

Example:
```bash
node run-tests.js --environment docker --folder "Movies Microservice" --reporters cli,htmlextra --bail
```

## Test Reports

After running the tests, HTML and JUnit XML reports will be generated in the `reports` directory. These reports can be used for CI/CD integration and documentation.

## CI/CD Integration

The current GitHub Actions workflow runs only the supported folders (`Monolith Service`, `Movies Microservice`) against a dedicated CI compose stack.

Example command sequence:

```yaml
- name: Run API Tests
  run: |
    cd tests/postman
    npm install
    npm run test:docker -- --folder "Monolith Service"
    npm run test:docker -- --folder "Movies Microservice"
```

## Troubleshooting

If you encounter issues running the tests:

1. Ensure all services are running and accessible
2. Check the environment configuration in the environment JSON files
3. Verify that the API endpoints match those in the collection
4. Increase the timeout value if requests are timing out