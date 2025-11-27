.PHONY: build
build:
	go build -v ./cmd/apiserver

.PHONY: test
test:
	go test -race -timeout 30s ./ ...

.PHONY: migrate-install
migrate-install:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

.PHONY: migrate-create
migrate-create:
	migrate create -seq -ext sql -dir ./migrations $(name)

.PHONY: migrate-up
migrate-up:
	migrate -path ./migrations -database "postgres://postgres:password@localhost:5432/ft?sslmode=disable" up

.PHONY: migrate-down
migrate-down:
	migrate -path ./migrations -database "postgres://postgres:password@localhost:5432/ft?sslmode=disable" down

.PHONY: migrate-status
migrate-status:
	migrate -path ./migrations -database "postgres://postgres:password@localhost:5432/ft?sslmode=disable" version

.PHONY: migrate-force
migrate-force:
	migrate -path ./migrations -database "postgres://postgres:password@localhost:5432/ft?sslmode=disable" force $(version)

.DEFAULT_GOAL := build