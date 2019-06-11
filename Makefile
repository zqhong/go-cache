.PHONY : fmt
fmt :
	@echo "Formatting your Go programs with gofmt..."
	@gofmt -l -w ./

.PHONY : test
test :
	@echo "Testing your code, please wait ..."
	@go test -race -cover -coverprofile=coverage.txt -covermode=atomic ./...