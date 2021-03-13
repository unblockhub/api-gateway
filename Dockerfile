FROM golang:alpine as builder

ADD . /app
WORKDIR /app

RUN go build -o /app/gateway

FROM alpine:latest

RUN addgroup -S app && adduser -S -G app app

WORKDIR /app
COPY --from=builder /app /app

EXPOSE 8080
USER app

ENTRYPOINT ["/app/gateway"]
