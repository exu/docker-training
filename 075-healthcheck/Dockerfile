FROM golang:1.13-alpine as builder
WORKDIR /go/src/app
COPY main.go ./
RUN go build -o /app

FROM alpine:3.11.3
RUN apk add curl
COPY --from=builder /app /app
EXPOSE 8080
HEALTHCHECK --retries=3 --interval=5s --timeout=2s CMD curl -f localhost:8080/hello
CMD ["/app"]
