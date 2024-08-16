FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN go build -o binary server.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/binary .

ENTRYPOINT ["./binary"]