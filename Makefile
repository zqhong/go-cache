.PHONY : fmt
fmt :
	@echo "Formatting your Go programs with gofmt..."
	@gofmt -l -w ./

.PHONY : lint
lint :
	@echo "Using golangci-lint to detect your code quality ..."
	@golangci-lint run

.PHONY : test
test :
	@echo "Testing your code, please wait ..."
	@go test -race -cover -coverprofile=coverage.txt -covermode=atomic ./...

.PHONY : release
release: fmt lint test
	@echo "Prepare to release the code, are performing code checking, please wait ..."