# Build Stage
FROM golang:1.24.1-alpine3.21 AS builder

WORKDIR /app
COPY . .
RUN go build -o bin/go-parking-lot-users-service cmd/http/main.go

# Run Stage
FROM alpine:3.21

WORKDIR /app

COPY ./db ./db
COPY --from=builder /app/bin/go-parking-lot-users-service .

EXPOSE 8080

ENTRYPOINT [ "/app/go-parking-lot-users-service" ]
