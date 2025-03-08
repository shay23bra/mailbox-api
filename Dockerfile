FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o mailbox-api ./main.go

FROM alpine:3.16

WORKDIR /app

RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /app/mailbox-api .
COPY --from=builder /app/.env .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/data ./data

RUN mkdir -p /var/log/mailbox-api

EXPOSE 8080

CMD ["./mailbox-api"]