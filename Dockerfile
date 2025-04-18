FROM golang:1.24 AS development
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go install github.com/cespare/reflex@latest
EXPOSE 3000
CMD reflex -g '*.go' go run cmd/server/main.go --start-service

FROM golang:1.24 AS builder
ENV GOOS linux
ENV CGO_ENABLED 0
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o app cmd/server/main.go

FROM alpine:latest AS production
RUN apk add --no-cache ca-certificates
RUN mkdir -p /app/pkg/config
COPY --from=builder /app/app /app
WORKDIR /app
EXPOSE 3000
CMD ./app