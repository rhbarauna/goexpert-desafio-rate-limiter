.PHONY: start run

start:
	@echo "Starting Docker containers..."
	docker-compose up -d
	@echo "Docker containers started successfully."
	@echo "Starting the application..."
	go run cmd/main.go

run: start