FROM golang:1.24-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o wubzduh .

FROM alpine:latest

WORKDIR /app

RUN mkdir -p /var/log/wubzduh

COPY --from=builder /build/wubzduh .

CMD ["./wubzduh", "serve", "--config", "/app/config.yaml"]
