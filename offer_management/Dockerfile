# Dockerfile

# Usa la imagen base apropiada
FROM golang:1.20

# Configura el directorio de trabajo
WORKDIR /app

# Copia el código fuente y el script
COPY . .

# Instala dependencias necesarias si es necesario
RUN go mod download

# Compila la aplicación
RUN go build -o myapp .

CMD ["./myapp"]
