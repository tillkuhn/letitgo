# syntax=docker/dockerfile:1
# https://docs.docker.com/language/golang/build-images/
# Build the application from source
FROM golang:1.20 AS build-stage

WORKDIR /app

COPY . .
#COPY go.mod go.sum ./
RUN go mod download

#COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /ltg

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
