run:
	docker-compose build && docker-compose up

migrate:
	migrate -path migrations -database 'postgres://postgres:postgres@localhost:5436/postgres?sslmode=disable' up

mod:
	go mod tidy
	go mod vendor

lint:
	golangci-lint run --fix