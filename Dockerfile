FROM golang:latest AS builder

WORKDIR /app

COPY cmd/indexer/main.go .

RUN go mod init indexer && \
    go mod tidy && \
    go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/indexer ./cmd/indexer/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/bin/indexer .
COPY --from=builder /app/web ./web
COPY --from=builder /app/uploads ./uploads
EXPOSE ${INDEXER_PORT:-8080}
CMD ["./indexer"]
