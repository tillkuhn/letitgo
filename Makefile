.ONESHELL:
.PHONY: help usage clean build image test format fmt lint outdated run
APPLICATION_NAME=ltg
# unified sed command, use like $(SED) -i -e 's/Id/ID/' some-file.txt
SED?=$(shell command -v gsed || command -v sed)

help: ## this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
	awk 'BEGIN {FS = ":.*?## ";{print "[%header,cols=\"1,2\"]\n|===\n| TARGET | DESCRIPTION"}}; {printf "| \033[36m%-30s\033[0m | %s\n", $$1, $$2}; END {print "|==="}'

usage: ## calls app with -h to show envconfig args
	go run main.go -h

clean: ## cleanup dist/ folder
	rm -rf dist/

# GOOS=linux
# CAUTION: make sure the import-path matches the module name in go.mod for -X
build: ## build in dist/app
	mkdir -p dist
	GOARCH=amd64 CGO_ENABLED=0 go build -ldflags \
	"-X 'github.com/tillkuhn/letitgo/cmd.CommitTag=$(shell git describe --tags --abbrev=0)' -X 'github.com/tillkuhn/letitgo/cmd.CommitHash=$(shell git rev-parse --short HEAD)' -X 'github.com/tillkuhn/letitgo/cmd.BuildDate=$(shell date +'%Y-%m-%dT%H:%M:%S')' -extldflags '-static'" \
	 -a -o ./dist/app

docker: ## local build with docker multistage
	docker build \
	--build-arg commit_hash=$(shell git rev-parse --short HEAD) \
	--build-arg commit_date="$(shell git show -s --format=%ci)" \
	--build-arg commit_tag="$(shell git describe --tags --abbrev=0)" \
	--build-arg build_date="$(shell date +'%Y-%m-%dT%H:%M:%S')" \
	-t $(APPLICATION_NAME) .
	docker images | grep $(APPLICATION_NAME)

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
	golangci-lint run --fix

.PHONY: update
update: ## update go dependencies
	go get -u

# https://github.com/psampaz/go-mod-outdated
outdated: ## show outdated direct dependencies
	go install github.com/psampaz/go-mod-outdated
	go list -u -m -json all | go-mod-outdated -direct

############ Aliases for Cobra Commands ############
.PHONY: run
run: ## runs app w/o args (shows help)
	@go run main.go

.PHONY: cache
cache: ## runs app with cache command
	go run main.go cache

.PHONY: charts
db: ## runs app with db command
	PGHOST=$(shell cat ~/.secret/cockroachdb/horsthost) PGDATABASE="horstdb" PGUSER="horst" \
	PGPASSWORD=$(shell cat ~/.secret/cockroachdb/horst) PGPORT="26257" \
	go run main.go db

.PHONY: db
charts: ## runs app with charts command
	@echo Running http://localhost:8081
	go run main.go db

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

.PHONY: rpc
rpc: ## run app with rpc (server) command
	go run main.go rpc

.PHONY: rpc-client
rpc-client: ## run app with rpc (client) command
	go run main.go rpc --client

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

.PHONY: version
version: ## run app with version command
	go run main.go version

.PHONY: worker
worker: ## run app with worker (job queue) command
	go run main.go worker

