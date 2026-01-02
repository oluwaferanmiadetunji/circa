DB_URL=postgresql://postgres:password@localhost:5432/circa?sslmode=disable

new_migration:
	migrate create -ext sql -dir internal/db/sqlc/migrations $(name)

sqlc:
	sqlc generate

# Install dependencies: go mod tidy & go mod vendor
install:
	go mod tidy && go mod vendor

# Run the server
run:
	air

# Start services concurrently using Air
start:
	docker compose up --build -d
	

stop:
	docker compose down --remove-orphans -v

lint:
	golangci-lint fmt

# OpenAPI code generation
openapi-generate:
	@echo "Generating API code from OpenAPI spec..."
	@go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest \
		-config openapi/oapi-codegen.yaml \
		openapi/openapi.yaml

# Validate OpenAPI spec
openapi-validate:
	@echo "Validating OpenAPI spec..."
	@go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest \
		-generate types \
		-package api \
		openapi/openapi.yaml > /dev/null && \
		echo "✓ OpenAPI spec is valid" || \
		(echo "✗ OpenAPI spec has errors" && exit 1)

# Generate and validate OpenAPI
openapi: openapi-validate openapi-generate

.PHONY: install start new_migration stop sqlc run lint openapi openapi-generate openapi-validate