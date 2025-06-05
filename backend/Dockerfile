# Etapa de build
FROM golang:1.23 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app .

# Etapa final
FROM debian:bullseye-slim
WORKDIR /app
COPY --from=builder /app/app .
COPY wait-for-db.sh .
RUN apt-get update && \
    apt-get install -y ca-certificates postgresql-client && \
    rm -rf /var/lib/apt/lists/* && \
    chmod +x wait-for-db.sh
EXPOSE 8080
CMD ["./wait-for-db.sh", "db", "./app"]

# Erasmo Cardoso da Silva
# Desenvolvedor Full Stack
