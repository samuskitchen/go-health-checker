# This starts a new build stage from the official Go image,
# which is based on the lightweight Alpine Linux
FROM golang:1.23.9-alpine3.22 AS deps
# This sets an environment variable to ensure Go
# uses Go Modules for package management.
ENV GO111MODULE=on
# This sets the working directory inside the container
# to /app. All subsequent commands will be run from this path.
WORKDIR /app
# This copies the two files that manage your project's dependencies
# (go.mod and go.sum) into the /app directory.
COPY go.mod go.sum ./
# This command downloads all the dependencies defined in your
# go.mod file and stores them in the image layer.
RUN go mod download

# This starts another stage from the same Go base image and names it build.
FROM golang:1.23.9-alpine3.21 AS build
# This sets the working directory inside the container
# to /app. All subsequent commands will be run from this path.
WORKDIR /app
# Copies all your project's source code (from your local machine)
# into the /app directory in the container.
COPY . .
# This is the build command.
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go build -ldflags="-s -w" -o go-health-checker cmd/main.go

# This starts the final image from a clean, minimal alpine:latest base image
FROM alpine:latest
# This copies only one file—the compiled go-health-checker binary—from the build stage
# into the root directory (.) of this final image.
COPY --from=build /app/go-health-checker .
# This informs Docker that the container listens on port 8080 at runtime.
# It's primarily documentation and helps with networking configuration.
EXPOSE 8080
# This sets the default command to run when a container starts from this image.
# It simply executes your compiled Go application.
CMD ["./go-health-checker"]