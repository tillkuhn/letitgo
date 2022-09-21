.ONESHELL:
.PHONY: help usage clean build image test format fmt lint outdated run

help: ## this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
	awk 'BEGIN {FS = ":.*?## ";{print "[%header,cols=\"1,2\"]\n|===\n| TARGET | DESCRIPTION"}}; {printf "| \033[36m%-30s\033[0m | %s\n", $$1, $$2}; END {print "|==="}'

usage: ## calls app with -h to show envconfig args
	go run main.go -h

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
# if you don't want to rely on external gotest tool with colors etc., simple replace gotest with this line:
# go test -v -coverpkg=./... -coverprofile=coverage.out ./...
test-all: ## run all tests (including long ones) with coverage report
	go install github.com/rakyll/gotest@v0.0.6
	gotest -v -coverpkg=./... -coverprofile=coverage.out ./...
	@go tool cover -func coverage.out | grep "total:"
	go tool cover -html=coverage.out -o coverage.html

test: ## run short tests
	go test -short ./...

imports: ## goimports -w -l .
	goimports -w -l .

lint: imports ## golangci-lint run (./... is implicit, also implies target imports)
	golangci-lint run

# https://github.com/psampaz/go-mod-outdated
outdated: ## show outdated direct dependencies
	go install github.com/psampaz/go-mod-outdated
	go list -u -m -json all | go-mod-outdated -direct

############ Aliases for Cobra Commands ############
.PHONY: run
run: ## runs app w/o args (shows help)
	@go run main.go

.PHONY: charts
charts: ## runs app with charts command
	@echo Running http://localhost:8081
	go run main.go charts

.PHONY: filesystem
filesystem: ## runs app with filesystem command
	go run main.go filesystem

.PHONY: oidc-client
oidc-client: ## runs app with oidcclient command
	go run main.go oidcclient --debug

.PHONY: oidc-server
oidc-server: ## run github.com/zitadel/oidc/example/server, serve http://localhost:9998/.well-known/openid-configuration
	go run github.com/zitadel/oidc/example/server

.PHONY: prometheus
prometheus: ## run app with prometheus command
	go run main.go prometheus

.PHONY: serve
serve: ## runs app with serve command
	go run main.go serve

.PHONY: signal
signal: ## runs app with signal cmd (graceful http shutdown on sigterm)
	go run main.go signal

.PHONY: sqlite
sqlite: ## runs app with sqlite command
	go run main.go sqlite

.PHONY: stack
stack: ## run app with stack command (generics support)
	go run main.go stack

.PHONY: ticker
ticker: ## runs app with ticker command
	go run main.go ticker

.PHONY: worker
worker: ## run app with worker (job queue) command
	go run main.go worker

