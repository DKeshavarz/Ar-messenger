
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git
ENV GOPROXY=https://goproxy.cn,direct

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
