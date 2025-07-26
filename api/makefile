# ai-manager-server/Makefile

# Variables
MIGRATE := migrate

# Load DB variables from .env file if it exists, otherwise use defaults
ifneq (,$(wildcard .env))
	include .env
	# Construct DATABASE_URL from individual variables if they exist
	DB_URL ?= postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)
else
	# Default connection string if no .env file exists
	DB_URL ?= postgres://postgres:secret@localhost:5432/wildbarrens?sslmode=disable
endif

MIGRATIONS_DIR := migrations

# Default target
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make migrate-up        - Apply all up migrations"
	@echo "  make migrate-down      - Roll back the last migration"
	@echo "  make migrate-down-all  - Roll back all migrations"
	@echo "  make migrate-create    - Create a new migration file (usage: make migrate-create NAME=your_migration_name)"
	@echo "  make migrate-force     - Force a specific version (usage: make migrate-force VERSION=123)"
	@echo "  make migrate-version   - Show the current migration version"
	@echo "  make populate-routes   - Populate database with GPX routes from docs/"
	@echo "  make stripe-listen     - Start Stripe webhook listener forwarding to localhost:8000/v1/webhooks/stripe"
	@echo "  make lmia-update       - Fetch and process LMIA data from Open Canada API"
	@echo "  make run               - Start the server with all environment variables loaded"
	@echo "  make docker-run        - Run the server in Docker with all environment variables loaded"
	@echo "  make deploy-helm       - Deploy the server using Helm"

# Apply all up migrations
.PHONY: migrate-up
migrate-up:
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up

# Roll back the last migration
.PHONY: migrate-down
migrate-down:
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down 1

# Roll back all migrations
.PHONY: migrate-down-all
migrate-down-all:
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down -all

# Create a new migration file
.PHONY: migrate-create
migrate-create:
ifndef NAME
	@echo "Error: Please provide a migration name with NAME=your_migration_name"
	@exit 1
endif
	$(MIGRATE) create -ext sql -dir $(MIGRATIONS_DIR) -seq $(NAME)

# Force a specific migration version
.PHONY: migrate-force
migrate-force:
ifndef VERSION
	@echo "Error: Please provide a version with VERSION=123"
	@exit 1
endif
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DB_URL)" force $(VERSION)

# Show the current migration version
.PHONY: migrate-version
migrate-version:
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DB_URL)" version

# Start Stripe webhook listener that forwards to localhost:8000/v1/webhooks/stripe
.PHONY: stripe-listen
stripe-listen:
	tailscale funnel localhost:8000

# Start the server with all environment variables loaded from .env
.PHONY: run
run:
	go run apps/server/main.go

# Run the server in Docker with all environment variables loaded from .env
.PHONY: docker-run
docker-run:
	docker build -t wildbarrens-server .
	docker run --env-file .env --network app-network -p 8000:8000 wildbarrens-server

# Populate database with GPX routes
.PHONY: populate-routes
populate-routes:
	cd scripts && DB_HOST=localhost DB_PORT=5432 DB_USER=postgres DB_PASSWORD=secret DB_NAME=wildbarrens go run populate_routes.go

# Fetch and process LMIA data from Open Canada API
.PHONY: lmia-update
lmia-update:
	@echo "Starting LMIA data update..."
	go run cmd/lmia_update.go

# Deploy helm image
.PHONY: deploy-helm
deploy-helm:
	helm upgrade --install wildbarrens ../helm/wildbarrens --values ../values-production.yaml
