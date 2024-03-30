FROM golang:1.22.1-alpine AS builder

WORKDIR /app

COPY . .

RUN go build -o weatherstation

FROM alpine:latest

COPY --from=builder /app/weatherstation /usr/local/bin/

EXPOSE 8080

CMD ["weatherstation"]
