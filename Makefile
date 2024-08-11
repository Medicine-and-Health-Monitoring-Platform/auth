CURRENT_DIR := $(shell pwd)
DATABASE_URL :=  postgres://postgres:123@localhost:5432/medicine?sslmode=disable

run-router:
	@go run cmd/router/main.go

run-service:
	@go run cmd/service/main.go

gen-proto:
	./scripts/gen-proto.sh "$(CURRENT_DIR)"

tidy:
	@go mod tidy
	@go mod vendor

mig-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir migrations -seq "$$name"

mig-up:
	@migrate -database "$(DATABASE_URL)" -path migrations up

mig-down:
	@migrate -database "$(DATABASE_URL)" -path migrations down

mig-force:
	@read -p "Enter migration version: " version; \
	migrate -database "$(DATABASE_URL)" -path migrations force "$$version"

permission:
	@chmod +x scripts/gen-proto.sh

test:
	@go test ./storage/postgres

swag-gen:
	~/go/bin/swag init -g ./api/router.go -o api/docs

run:
	@go run cmd/main.go