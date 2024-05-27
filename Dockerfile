FROM golang:1.21.7-alpine AS builder

WORKDIR /app

COPY --chown=app:app . .

RUN go build -ldflags="-s -w" -o bin/webhook-listener

FROM alpine:3.20.0

COPY --from=builder /app/bin /

EXPOSE 3032

CMD ["/webhook-listener"]