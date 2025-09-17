FROM golang:1.23-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git build-base
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o chatapp ./cmd/server

FROM alpine:3.20
WORKDIR /app
RUN adduser -D -g '' appuser
COPY --from=builder /app/chatapp /usr/local/bin/chatapp
COPY static ./static
USER appuser
EXPOSE 8000
ENV SERVER_PORT=8000
CMD ["/usr/local/bin/chatapp"]

