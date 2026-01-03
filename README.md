# Running the App

## Prerequisites

* Docker and Docker Compose installed
* Go installed
* No other PostgreSQL instance running on port 5432

## Setup

### 1. Create environment file

Create a `.env` file from `example.env` in the project root.

This file is required for both local development and production.

### 2. Start PostgreSQL

Start the PostgreSQL container before running the application:

```bash
docker-compose up -d
```

To stop and remove the container along with its volumes:

```bash
docker-compose down -v
```

### 3. Build the application

Build the Go binary:

```bash
make build
```

If Go is not located at `/usr/bin/go`, update the path in the `Makefile`
(for example, `/bin/go`).

If you want a different binary name, update it in the `Makefile`
to match your application’s desired name.

### 4. Clean build artifacts

Remove the generated binary:

```bash
make clean
```

### 5. Run the application

Build and run the binary:

```bash
make run
```
