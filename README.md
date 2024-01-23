# Rate Limiter

## Overview

The Rate Limiter restricts access to the web server based on the defined configuration in the .env file.

It limits the number of requests that exceed the configured values.

The limiter performs analysis based on the provided token in the request header with the key API_KEY.

If a configuration is found for that token, the specified limits and blocking time will be applied.

If no token is passed or the token does not have a configuration, default values will be applied based on the IP address of the request.

The Rate Limiter is a tool designed to control the rate of incoming requests to a web server.
It works by analyzing tokens and/or IP addresses and applying rate-limiting policies based on the configuration provided.

Feel free to customize the configuration to suit your application's needs.

## Getting Started

### Prerequisites

Before you begin, make sure you have Docker and Go installed on your machine.

### Configuration

The application can be configured via the `.env` file located in the `/cmd` directory. Here is an example of the `.env` file:

```env
MAX_REQUESTS=10
TOKENS=[{"name": "tkn_123", "max_requests": 20, "cooldown_seconds": 3}, {"name": "tkn_456", "max_requests": 30, "cooldown_seconds":4}]
TTL_SECONDS=1
COOLDOWN_SECONDS=5
WEB_SERVER_PORT=:8080
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
```

**Adjust the values in the `.env` file according to your desired configuration.**

### Usage

To start the application, run the following command:

```bash
make run # OR make start
```

This command will start the application using Docker Compose and then run the main.go file.
