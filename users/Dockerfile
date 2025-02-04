# Build Stage
FROM golang:1.22.4-alpine3.19 AS builder

WORKDIR /app
COPY . .
RUN go build -o bin/go-parking-lot-users-service cmd/http/main.go

# Run Stage
FROM alpine:3.19

WORKDIR /app

COPY ./db ./db
COPY --from=builder /app/bin/go-parking-lot-users-service .

EXPOSE 8080

ENTRYPOINT [ "/app/go-parking-lot-users-service" ]