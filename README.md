# Rate Limiter

## Overview

A tool designed to control the rate of incoming requests to a web server.
It works by analyzing tokens and/or IP addresses and applying rate-limiting policies based on the configuration provided.

The limiter performs analysis based on the provided token in the request header with the key API_KEY.
If a configuration is found for that token, the specified limits and blocking time will be applied.

If no token is passed or the token does not have a configuration, default values will be applied based on the IP address of the request.

Feel free to customize the configuration to suit your application's needs.

## Technologies Used

The rate limiter is built using the following technologies:

- [viper](https://github.com/spf13/viper): A configuration management library for Go.
- [redis](https://redis.io/): A powerful, open-source, in-memory data structure store.
- [go-redis](https://github.com/redis/go-redis): A Go client library for the Redis database.
- [go-chi](https://github.com/go-chi/chi): A lightweight, idiomatic, and composable router for building Go HTTP services.
- [testify](https://github.com/stretchr/testify): A testing toolkit for Go.
- [wire](https://github.com/google/wire): Automated Initialization in Go.

## Customizable Storage

The rate limiter uses Redis as the default storage backend.
However, you have the flexibility to substitute Redis with a different storage
solution, as long as the new storage implements the `Storage` interface found in the `storage` package:

```go
type Storage interface {
    Increment(ctx context.Context, key string, ttl int) (int, error)
    Get(ctx context.Context, key string) (interface{}, error)
    Set(ctx context.Context, key string, ttl int) error
    Exists(ctx context.Context, key string) (bool, error)
    IsBlocked(ctx context.Context, key string) (bool, error)
}
```

## Getting Started

### Prerequisites

Before you begin, make sure you have Docker and Go installed on your machine.

### Configuration

The application can be configured via the `.env` file located in the `/cmd` directory. Here is an example of the `.env` file:

```env
MAX_REQUESTS=10
TOKENS=[{"name": "tkn_123", "max_requests": 20, "cooldown_seconds": 3}, {"name": "tkn_456", "max_requests": 30, "cooldown_seconds": 4}]
TTL_SECONDS=1
COOLDOWN_SECONDS=5
WEB_SERVER_PORT=:8080
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
```

**Adjust the values in the `.env` file according to your desired configuration.**

### Explanation of `.env` Items

- **MAX_REQUESTS**: Specifies the maximum number of requests allowed within the time window defined by `TTL_SECONDS`.
- **TTL_SECONDS**: Sets the time window (in seconds) during which the rate limiter analyzes and enforces the maximum number of requests specified by `MAX_REQUESTS`.
- **COOLDOWN_SECONDS**: Sets the default blocking time (in seconds) applied when the request limit is exceeded.
- **TOKENS**: Configures a list of tokens in the following format:

  ```json
  [
    { "name": "token_1", "max_requests": 20, "cooldown_seconds": 3 },
    { "name": "token_2", "max_requests": 30, "cooldown_seconds": 4 }
  ]
  ```

- **name**: The token identifier.
- **max_requests**: The maximum number of requests allowed within the time window (`TTL_SECONDS`) for this specific token.
- **cooldown_seconds**: The blocking time applied when the request limit is exceeded for this specific token.

### Usage

To start the application, run the following command:

```bash
make run # OR make start
```

This command will start the application using Docker Compose and then run the main.go file.

To execute all tests, run the following command:

```bash
make run-tests
```
