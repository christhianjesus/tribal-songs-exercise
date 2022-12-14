# syntax=docker/dockerfile:1
FROM golang:1.19-bullseye AS build

WORKDIR /usr/src/app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /echo-server cmd/main.go

## Deploy
FROM gcr.io/distroless/base-debian11

WORKDIR /

COPY --from=build /echo-server /echo-server

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/echo-server"]
