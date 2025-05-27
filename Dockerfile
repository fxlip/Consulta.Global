# Estágio de construção
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /busca-cpf

# Estágio de produção (imagem final leve)
FROM alpine:latest
RUN mkdir -p /app
WORKDIR /app
COPY --from=builder /busca-cpf .
COPY .env .
EXPOSE 8080
CMD ["./busca-cpf"]