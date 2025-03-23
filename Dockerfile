FROM golang:1.23.5-alpine AS builder
RUN apk add --no-cache git
WORKDIR /app
COPY . .

RUN go install github.com/swaggo/swag/cmd/swag@latest && swag init -g cmd/onlineLibrary/main.go
RUN go mod download
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o libmusic ./cmd/onlineLibrary/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/libmusic .

EXPOSE 8080

CMD ["./libmusic"]