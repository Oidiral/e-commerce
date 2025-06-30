FROM golang:1.24.2 AS builder
WORKDIR /app
COPY go.mod go.sum ./
COPY migrations ./
RUN go mod download
ENV CGO_ENABLED=0
RUN go install github.com/pressly/goose/v3/cmd/goose@latest
COPY . .
RUN go build -ldflags="-s -w" -o cart ./cmd

FROM alpine:3.18 AS runtime
RUN apk add --no-cache ca-certificates postgresql-client
COPY --from=builder /go/bin/goose /usr/local/bin/goose
WORKDIR /app
COPY --from=builder /app/cart ./cart
COPY --from=builder /app/migrations ./db/migrations
COPY entrypoint.sh ./entrypoint.sh
RUN chmod +x entrypoint.sh
ENTRYPOINT ["./entrypoint.sh"]
EXPOSE 8081