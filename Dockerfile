# syntax=docker/dockerfile:1

FROM golang:1.22.5 AS builder
WORKDIR /app

COPY go.mod ./
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /test-result

# --
FROM alpine
WORKDIR /app
COPY --from=builder /test-result /test-result
COPY static/* /app/static/
COPY templates/* /app/templates/
COPY favicon/* /app/favicon/

EXPOSE 8080
CMD ["/test-result"]