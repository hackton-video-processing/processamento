FROM golang:1.23.4 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

WORKDIR /app/cmd/api
RUN CGO_ENABLED=1 GOOS=linux go build -o /main .

FROM debian:bookworm-slim
WORKDIR /root/
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
COPY --from=builder /main .
EXPOSE 8080
ENTRYPOINT ["./main"]
