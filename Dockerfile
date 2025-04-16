FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
COPY . .
RUN go build -o /bin/pixellock

FROM alpine:latest
WORKDIR /app
COPY --from=builder /bin/pixellock /usr/local/bin/pixellock
ENTRYPOINT [ "/usr/local/bin/pixellock" ]
