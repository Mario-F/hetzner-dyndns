# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:1.20-buster AS build

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./

ARG BUILD_VERSION=development
RUN go build -ldflags="-X 'github.com/Mario-F/hetzner-dyndns/cmd.Version=${BUILD_VERSION}'" -o /hetzner-dyndns

##
## Deploy
##
FROM gcr.io/distroless/base-debian11

WORKDIR /

COPY --from=build /hetzner-dyndns /hetzner-dyndns

USER nonroot:nonroot

ENTRYPOINT ["/hetzner-dyndns"]
