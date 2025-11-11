# API GO HEALTH CHECKER

## Prerequisites

Before you begin, ensure you have met the following requirements:
* You have installed the latest version 1.23.9 of [Go](https://go.dev/dl/)
* You have installed [Docker](https://docs.docker.com/desktop/)

## Installation Dependencies
```bash
  make modd
```

### Generate Swag (Swagger):
Swag converts Go annotations to Swagger Documentation 2.0. We've created a variety of plugins for popular Go web frameworks. This allows you to quickly integrate with an existing Go project (using Swagger UI).

**Link:** https://github.com/swaggo/swag

#### Install Swagger:
```bash
  make install-swag
```

##### command generate Swagger doc
```bash
  make swag
```

### Linter (golangci-lint):
Golangci-lint is a fast linters runner for Go.
It runs linters in parallel, uses caching, supports YAML configuration, integrates with all major IDEs, and includes over a hundred linters.

**Link:** https://github.com/golangci/golangci-lint

#### Install by Binary:
```bash
  make install-lint-binaries
```
##### Windows

###### Install by Chocolatey:
```bash
  make install-lint-windows-chocolatey
```
###### Install by Scoop:
```bash
  make install-lint-windows-scoop
```

##### MAC

###### Install by Homebrew:
```bash
  make install-lint-mac-silicon
```

#### Install Dependencies for Linter:
```bash
  make install-lint-dependencies
```

### Generate Mock Interface
This is an automatic mock generator using mockery, the first thing we must do is go to the path of the file that we want to autogenerate:

**Link:** https://github.com/vektra/mockery

Install the library
```bash
  make mockery-install
```
#### Command:
```bash
  mockery --config .mockery.yml
```

Generate all mocks with expected feat include, see official documentation. This line would go to the beginning of the test file

## Execute Test
```bash
  make go-test
```

#### Execute Report Coverage:
```bash
  make go-test-report
```

## Run application

prior to execution, you must include the following environment variables
```
DB_USER=postgres
DB_PASSWORD=admin
DB_NAME=beer-api
DB_HOST=localhost
DB_PORT=5432
DB_MAX_OPEN=50
DB_MAX_IDLE=25
DB_CONN_MAX_LIFETIME='15m'
DB_CONN_MAX_IDLE_TIME='5m'
DB_SSL_MODE=disable

LOGGER_DEBUG=false
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
SERVER_POSTFIX=dev

RABBITMQ_HOST=localhost
RABBITMQ_PORT=5672
RABBITMQ_USERNAME=guest
RABBITMQ_PASSWORD=guest

HAZEL_SERVER=localhost:5701
```

### Execute go build
```bash
  make build
```

##### build apple-silicon
```bash
  make build-apple-silicon
```