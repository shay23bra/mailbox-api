.PHONY: build clean test run docker-up docker-down

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=mailbox-api

all: test build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

test:
	$(GOTEST) -v ./...

test-coverage:
	$(GOTEST) -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

run:
	$(GOBUILD) -o $(BINARY_NAME) -v
	./$(BINARY_NAME)

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-build:
	docker-compose build

# Database-related commands
db-init:
	cd scripts && ./init_db.sh

db-seed:
	cd scripts && ./seed_db.sh

# For development
dev: docker-up db-init db-seed run

# For CI/CD
ci: test build