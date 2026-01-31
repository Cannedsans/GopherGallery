#!/bin/bash

# Define the container name
CONTAINER_NAME="testeLocalStack"

# Check if the container is running
if docker ps --filter "name=$CONTAINER_NAME" --format '{{.Names}}' | grep -q $CONTAINER_NAME; then
  echo "Container $CONTAINER_NAME is running. Starting the Go application..."
  go run ./cmd/api/main.go
else
  echo "Container $CONTAINER_NAME is not running. Starting it with docker-compose..."
  docker compose up -d
  # Optional: Add a small delay to allow the container to start properly
  sleep 5
  echo "Container $CONTAINER_NAME started. Starting the Go application..."
  go run ./cmd/api/main.go
fi