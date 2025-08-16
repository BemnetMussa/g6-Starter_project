FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum from root (project root)
COPY go.mod go.sum ./
RUN go mod download

# Copy everything (including Delivery folder)
COPY . .

# Build your main.go inside Delivery/
RUN go build -o g6blog ./Delivery/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/g6blog .

EXPOSE 8080

CMD ["./g6blog"]
