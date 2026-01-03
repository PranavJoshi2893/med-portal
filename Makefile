BINARY_NAME=med-portal
PATH=./cmd/api/main.go

.PHONY: build run clean

build:
		/usr/bin/go build -o $(BINARY_NAME) $(PATH)

run: build
		./$(BINARY_NAME)

clean:
		/usr/bin/go clean
		/usr/bin/rm -f ./$(BINARY_NAME)