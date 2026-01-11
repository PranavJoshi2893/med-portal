include .env
export

BINARY_NAME=med-portal
MAIN_PATH=./cmd/api/main.go
DOCKER := /usr/bin/docker

.PHONY: build run clean

build:
		/usr/bin/go build -o $(BINARY_NAME) $(MAIN_PATH)

run: build
		./$(BINARY_NAME)

clean:
		/usr/bin/go clean
		/usr/bin/rm -f ./$(BINARY_NAME)

migrate-up:
		$(DOCKER) run -v ./migrations:/migrations --network host migrate/migrate \
		-path=/migrations/ \
		-database "postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:5432/$(POSTGRES_DB)?sslmode=disable" \
		up 1

migrate-up-all:
		$(DOCKER) run -v ./migrations:/migrations --network host migrate/migrate \
		-path=/migrations/ \
		-database "postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:5432/$(POSTGRES_DB)?sslmode=disable" \
		up

migrate-down:
		$(DOCKER) run -v ./migrations:/migrations --network host migrate/migrate \
		-path=/migrations/ \
		-database "postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:5432/$(POSTGRES_DB)?sslmode=disable" \
		down 1