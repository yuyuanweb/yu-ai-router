FROM golang:1.24-alpine AS builder
WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o server ./cmd/server

FROM alpine:3.22
WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata wget

COPY --from=builder /app/server /app/server

EXPOSE 8123

HEALTHCHECK --interval=30s --timeout=3s --start-period=20s --retries=5 \
  CMD wget --no-verbose --tries=1 -O /dev/null http://localhost:8123/api/health/ || exit 1

ENTRYPOINT ["/app/server", "--env", "prod"]
