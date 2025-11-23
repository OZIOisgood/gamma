FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o bin/api ./cmd/api
RUN go build -o bin/worker ./cmd/worker

FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache ffmpeg

COPY --from=builder /app/bin/api .
COPY --from=builder /app/bin/worker .
COPY --from=builder /app/assets ./assets
# Copy migrations if needed by the app (though app doesn't seem to run migrations itself, makefile does)

EXPOSE 8080

CMD ["./api"]
