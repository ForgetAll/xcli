lint:
	@echo Lint start
	@golangci-lint run -v ./...
	@echo Lint finish

test:
	@echo Test start
	@go test ./...
	@echo Test finish

clean:
	@rm -f *.coverprofile
	@rm -f coverage.*
	@echo Clean Finish

all: lint test
.PHONY: lint test cover clean 
