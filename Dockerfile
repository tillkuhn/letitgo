# syntax=docker/dockerfile:1
# Build the application from source
# https://docs.docker.com/language/golang/build-images/
# Document setting ldflags to embed version into docker image / go binary
# https://github.com/ko-build/ko/issues/167

FROM golang:1.20 AS build-stage

# Build arguments for this image (to be used in ldflags)
ARG commit_hash=""
ARG commit_date=""
ARG commit_tag=""
ARG build_date=""

WORKDIR /app

# use .dockerignore so .idea,.git etc. won't be pushed to the build context
COPY . .
RUN go mod download

# go build -ldflags="-help"
# -X definition    add string value definition of the form importpath.name=value
# https://programmingpercy.tech/blog/modify-variables-during-build/
# CAUTION: make sure the import-path matches the module name in go.mod for -X
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath \
    -ldflags="-w -s \
    -X 'github.com/tillkuhn/letitgo/cmd.CommitHash=${commit_hash}' \
    -X 'github.com/tillkuhn/letitgo/cmd.CommitDate=${commit_date}' \
    -X 'github.com/tillkuhn/letitgo/cmd.CommitTag=${commit_tag}' \
    -X 'github.com/tillkuhn/letitgo/cmd.BuildDate=${build_date}' \
    -extldflags '-static'" \
    -a -o /ltg

# Run the tests in the container
FROM build-stage AS run-test-stage
RUN go test -v ./...

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /ltg /ltg

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/ltg"]
