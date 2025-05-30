up:
	docker-compose up -d

down:
	docker-compose down

lint:
	golangci-lint run ./...