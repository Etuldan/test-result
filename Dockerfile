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
HEALTHCHECK --interval=30s --timeout=30s --start-period=5s --retries=3 CMD curl --fail http://localhost:8080 || exit 1
CMD ["/test-result"]