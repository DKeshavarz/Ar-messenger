
FROM golang:1.17-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o messenger ./cmd/messenger

FROM alpine:latest

RUN adduser -D appuser
WORKDIR /app
COPY --from=builder /app/messenger .
USER appuser

EXPOSE 8080
CMD ["./messenger"]
