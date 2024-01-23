.PHONY: start run

start:
	docker-compose up -d
	go run cmd/main.go

run: start