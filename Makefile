MAIN_PACKAGE_PATH := ./cmd/main.go
BINARY_NAME := mvc

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

.PHONY: no-dirty
no-dirty:
	git diff --exit-code

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## tidy: format code and tidy modfile
.PHONY: tidy
tidy:
	go fmt ./...
	go mod tidy -v

## audit: run quality control checks
.PHONY: audit
audit:
	go mod verify
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...
	go test -race -buildvcs -vet=off ./...

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## test: run all tests
.PHONY: test
test:
	go test -v -race -buildvcs ./...

## test/cover: run all tests and display coverage
.PHONY: test/cover
test/cover:
	go test -v -race -buildvcs -coverprofile=/tmp/coverage.out ./...
	go tool cover -html=/tmp/coverage.out

## build: build the application
.PHONY: build
build:
	# Include additional build steps, like TypeScript, SCSS or Tailwind compilation here...
	go build -o=/tmp/bin/${BINARY_NAME} ${MAIN_PACKAGE_PATH}

## run: run the  application
.PHONY: run
run: build-env migrate-up build check-migrate
	/tmp/bin/${BINARY_NAME}

## run/live: run the application with reloading on file changes
.PHONY: run/live
run/live:
	go run github.com/cosmtrek/air@v1.43.0 \
		--build.cmd "make build" --build.bin "/tmp/bin/${BINARY_NAME}" --build.delay "100" \
		--build.exclude_dir "" \
		--build.include_ext "go, tpl, tmpl, html, css, scss, js, ts, sql, jpeg, jpg, gif, png, bmp, svg, webp, ico" \
		--misc.clean_on_exit "true"

# ==================================================================================== #
# DATABASE
# ==================================================================================== #

MIGRATE ?= $(shell go env GOPATH)/bin/migrate

.PHONY: migrate-up migrate-down migrate-new

migrate-up:
	@echo "Running database migrations up..."
	$(MIGRATE) -database "$(DB)" -path ./migrations up

migrate-down:
	@echo "Running database migrations down..."
	$(MIGRATE) -database "$(DB)" -path ./migrations down

migrate-new:
	@read -p "Enter migration name: " name; \
	$(MIGRATE) create -ext sql -dir ./migrations -seq $$name

# ==================================================================================== #
# OPERATIONS
# ==================================================================================== #

## build-env: prompt for environment variables and construct DB URL
.PHONY: build-env
build-env:
	@read -p "Enter port (default: 3000): " PORT; \
	read -p "Enter DB username (default: root): " DB_USERNAME; \
	read -sp "Enter DB password (default: password): " DB_PASSWORD; echo; \
	read -p "Enter DB host (default: localhost): " DB_HOST; \
	read -p "Enter DB port (default: 3306): " DB_PORT; \
	read -p "Enter DB name (default: bookstore): " DB_NAME; \
	read -p "Enter JWT secret (default: your_jwt_secret): " SECRET; \
	export PORT=$${PORT:-3000}; \
	export DB_USERNAME=$${DB_USERNAME:-root}; \
	export DB_PASSWORD=$${DB_PASSWORD:-password}; \
	export DB_HOST=$${DB_HOST:-localhost}; \
	export DB_PORT=$${DB_PORT:-3306}; \
	export DB_NAME=$${DB_NAME:-bookstore}; \
	export SECRET=$${SECRET:-your_jwt_secret}; \
	export DB="$$DB_USERNAME:$$DB_PASSWORD@tcp($$DB_HOST:$$DB_PORT)/$$DB_NAME?charset=utf8mb4&parseTime=True&loc=Local"; \
	echo "PORT=$$PORT" > .env; \
	echo "DB_USERNAME=$$DB_USERNAME" >> .env; \
	echo "DB_PASSWORD=$$DB_PASSWORD" >> .env; \
	echo "DB_HOST=$$DB_HOST" >> .env; \
	echo "DB_PORT=$$DB_PORT" >> .env; \
	echo "DB_NAME=$$DB_NAME" >> .env; \
	echo "DB=$$DB" >> .env; \
	echo "SECRET=$$SECRET" >> .env; \
	echo "Environment variables set and .env file created."

## push: push changes to the remote Git repository
.PHONY: push
push: tidy audit no-dirty
	git push

## production/deploy: deploy the application to production
.PHONY: production/deploy
production/deploy: confirm tidy audit no-dirty
	GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o=/tmp/bin/linux_amd64/${BINARY_NAME} ${MAIN_PACKAGE_PATH}
	upx -5 /tmp/bin/linux_amd64/${BINARY_NAME}

# ==================================================================================== #
# INSTALLATION
# ==================================================================================== #

## check-migrate: check if golang-migrate is installed, and install if not
.PHONY: check-migrate
check-migrate:
	@if ! command -v migrate &> /dev/null; then \
		echo "golang-migrate is not installed. Installing..."; \
		go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest; \
	fi