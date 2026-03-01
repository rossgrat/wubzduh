FROM golang:1.24-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o wubzduh .

FROM alpine:3.21

RUN adduser -D -u 10001 appuser

WORKDIR /app

RUN mkdir -p /var/log/wubzduh && chown appuser:appuser /var/log/wubzduh

COPY --from=builder /build/wubzduh .

USER appuser

CMD ["./wubzduh", "serve", "--config", "/app/config.yaml"]
