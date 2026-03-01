FROM golang:1.24-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o wubzduh .

FROM alpine:3.21

RUN apk add --no-cache su-exec && adduser -D -u 10001 appuser

WORKDIR /app

RUN mkdir -p /var/log/wubzduh

COPY --from=builder /build/wubzduh .
COPY entrypoint.sh .
RUN chmod +x entrypoint.sh

ENTRYPOINT ["./entrypoint.sh"]
CMD ["./wubzduh", "serve", "--config", "/app/config.yaml"]
