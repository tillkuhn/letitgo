.ONESHELL:
.PHONY: help usage clean build image test format fmt lint outdated run

help: ## this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
	awk 'BEGIN {FS = ":.*?## ";{print "[%header,cols=\"1,2\"]\n|===\n| TARGET | DESCRIPTION"}}; {printf "| \033[36m%-30s\033[0m | %s\n", $$1, $$2}; END {print "|==="}'

usage: ## calls app with -h to show envconfig args
	go run *.go -h

clean: ## cleanup dist/ folder
	rm -rf dist/

build: ## build in dist/app
	mkdir -p dist
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-extldflags '-static'" -o ./dist/app

image: build ## local docker build
	docker build -t $(APPLICATION_NAME) .
	docker images|grep $(APPLICATION_NAME)

# read https://blog.seriesci.com/how-to-measure-code-coverage-in-go/ for total coverage
# and https://stackoverflow.com/questions/33444968/how-to-get-all-packages-code-coverage-together-in-go
test: ## go test everything
	go test -v -coverpkg=./... -coverprofile=coverage.out ./...
	@go tool cover -func coverage.out | grep "total:"
	go tool cover -html=coverage.out -o coverage.html

format: ## goimports -w -l .
	goimports -w -l .

fmt: format  ## alias for format

lint: format ## golangci-lint run (./... is implicit)
	golangci-lint run

# https://github.com/psampaz/go-mod-outdated
outdated: ## show outdated direct dependencies
	go install github.com/psampaz/go-mod-outdated
	go list -u -m -json all | go-mod-outdated -direct

job-queue: ## go run cmd/run/main.go
	go run cmd/run/main.go job-queue


prometheus: ## go run cmd/run/main.go
	go run cmd/run/main.go prometheus

