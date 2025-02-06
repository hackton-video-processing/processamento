# Etapa 1: Build
FROM golang:1.23.4 AS builder

WORKDIR /app

# Copiar e instalar dependências
COPY go.mod go.sum ./
RUN go mod download

# Copiar o código fonte
COPY . .

# Compilar o binário de forma estática
WORKDIR /app/cmd/api
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /main .

# Etapa 2: Imagem final
FROM scratch

WORKDIR /
COPY --from=builder /main .

# Definir o ponto de entrada
ENTRYPOINT ["/main"]
