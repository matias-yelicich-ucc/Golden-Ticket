# Stage 1: Build the Go binary
FROM golang:alpine AS builder

# Instalar git para dependencias si son necesarias
RUN apk add --no-cache git

WORKDIR /app

# Descargar dependencias aprovechando la caché de Docker
COPY go.mod go.sum ./
RUN go mod download

# Copiar el resto del código fuente
COPY . .

# Compilar el binario optimizado
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server main.go

# Stage 2: Imagen mínima de ejecución
FROM alpine:latest

# Instalar ca-certificates y zona horaria por si acaso
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Copiar el binario compilado desde la etapa de construcción
COPY --from=builder /app/server .

# Exponer el puerto por el que escucha la API
EXPOSE 8080

# Arrancar la aplicación
CMD ["./server"]
