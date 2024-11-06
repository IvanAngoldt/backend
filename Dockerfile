FROM golang:1.23 AS builder

WORKDIR /app
COPY . .

RUN go mod tidy

RUN GOOS=linux GOARCH=amd64 go build -o server .

FROM debian:latest

WORKDIR /root/
COPY --from=builder /app/server .

EXPOSE 5000

CMD ["./server"]
