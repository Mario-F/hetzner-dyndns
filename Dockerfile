# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:1.15-buster AS build

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./

RUN go build -o /hetzner-dyndns

##
## Deploy
##
FROM gcr.io/distroless/base-debian11

WORKDIR /

COPY --from=build /hetzner-dyndns /hetzner-dyndns

USER nonroot:nonroot

ENTRYPOINT ["/hetzner-dyndns"]
