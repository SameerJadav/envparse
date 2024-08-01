# Change these variables as necessary.
PACKAGE_PATH := ./...

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## test: run all tests with coverage
.PHONY: test
test:
	@echo "Running tests..."
	@go test -v -race -cover -buildvcs ${PACKAGE_PATH}

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## tidy: format code and tidy modfile
.PHONY: tidy
tidy:
	@echo "Formatting code..."
	@gofumpt -l -w .
	@echo "Tidying Go mod..."
	@go mod tidy -v
