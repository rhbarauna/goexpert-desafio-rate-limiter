.PHONY: start run run-tests

start:
	@echo "Starting Docker containers..."
	docker-compose up -d
	@echo "Docker containers started successfully."
	@echo "Starting the application..."
	go run cmd/main.go cmd/wire_gen.go

run: start

run-tests:
	go test ./... -v
