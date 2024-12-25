## MusicBrainz Feed Integration with ZITADEL Authentication

A web application that demonstrates ZITADEL SDK capabilities by integrating the MusicBrainz Feed API with authentication.

### Features

- Authentication via ZITADEL
- Role-based access control
- MusicBrainz feed api integration
- Custom feed selection for admin users
- Health check integration into the webapp 

### Project Structure

The project consists of two main components:
- [app](app) - Authorization API service 
- [web](web) - Web application frontend 

### Application Routes

### WebApp

| Route | Description |
|-------|-------------|
| `/` | Home page (redirects to `/feed` if logged in) 
| `/feed` | Feed display page with admin controls for feed selection |
| `/select_feed` | Select another musicbrainz feed |

### API Service

| Route | Description | Authentication |
|-------|-------------|----------------|
| `/api/healthz` | Health check endpoint | None |
| `/api/feed` | Feed data endpoint with health monitoring | Required |
| `/api/select_feed` | Feed selection endpoint | Required + Admin role |


## Setup

### Prerequisites

- Docker
- Go 1.x 
- ZITADEL account and project configuration

### Environment Configuration

Create a `.env` file in the root directory based on the provided [`env.example`](env.example) template. Required variables:

```bash
API_PORT=        # API service port
WEB_PORT=        # Web application port
CLIENT_ID=       # ZITADEL project client ID
WEB_KEY=         # Webapp encryption key
API_HOSTNAME=    # API hostname for webapp communication
REDIRECT_URI=    # OAuth redirect URI
KEY_FILE=        # API private key path
```

## Development/Deployment Options

### Docker

Start both applications using Docker Compose:

```bash
docker compose up -d
```

### Local Development

Run the web application:

```bash
go run cmd/web/main.go \
    -domain ${DOMAIN} \
    --clientID ${CLIENT_ID} \
    -key ${WEB_KEY} \
    -redirectURI ${REDIRECT_URI} \
    -port ${WEB_PORT}
```

Run the API service:

```bash
go run cmd/api/main.go \
    -domain ${DOMAIN} \
    -port ${API_PORT} \
    -key ${KEY_FILE}
```

### Building From Source

Generate binaries in the `bin/` directory:

```bash
# Build both applications
make build

# Build individual components
make build-api
make build-web
```

### Running tests

```bash
make test
```

```bash
# run tests with race detector
make test-race
```

## Documentation

- [Securing APIs with ZITADEL](https://zitadel.com/docs/examples/secure-api/go)
- [ZITADEL Login Integration](https://zitadel.com/docs/examples/login/go)
- [ZITADEL Go SDK](https://github.com/zitadel/zitadel-go)
- [Additional Examples](https://github.com/zitadel/zitadel-go/tree/next/example)

## Next

### Tests

- Add unit tests for both web and API
- Add integration tests for API endpoints
- Integrate those tests into a CI/CD pipelines

### Logging

- Add request/response logging middleware
- Add structured logging for both web and API

### MusicBrainz API

MusicBrainz API is rate-limited. 

- Implement rate limiting
- Cache MusicBrainz responses 

## Notes 

This project was inspired by examples from the [ZITADEL Go SDK repository](https://github.com/zitadel/zitadel-go/tree/next/example).
