FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY backend/go.mod backend/go.sum* ./
RUN go mod download

COPY backend/ .
RUN CGO_ENABLED=0 GOOS=linux go build -o shortfyurl .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/shortfyurl .
COPY frontend/ ./frontend/

EXPOSE 8080

CMD ["./shortfyurl"]
