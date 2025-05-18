FROM golang:1.24 AS builder
WORKDIR /app

COPY go.mod go.sum* ./

RUN go mod download

COPY ..

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/myapp .

FROM debian:bookworm-slim
RUN apk add --no-cache redis

COPY --from=builder /app/myapp /app/myapp

COPY startServices.sh /app/startServices.sh
RUN chmod +x /app/startServices.sh

EXPOSE 8080

CMD ["/app/startServices.sh"]
