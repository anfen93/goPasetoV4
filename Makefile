test:
	go test -v -cover ./...
checkTestCoverage:
	@echo "Running tests with coverage..."
	@go test -coverprofile=coverage.out ./...
	@echo "Checking for missing test coverage..."
	@go tool cover -func=coverage.out | grep -v "100.0%"
	@echo "Cleaning up..."
ifeq ($(OS),Windows_NT)
	@del coverage.out
else
	@rm -f coverage.out
endif


.PHONY: test checkTestCoverage