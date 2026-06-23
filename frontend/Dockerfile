# Stage 1: Construcción de la aplicación React con Vite
FROM node:20-alpine AS builder

WORKDIR /app

# Copiar archivos de dependencias e instalarlas usando npm ci
COPY package*.json ./
RUN npm ci

# Copiar el resto del código
COPY . .

# Argumento de construcción para la URL de la API
ARG VITE_API_URL
ENV VITE_API_URL=$VITE_API_URL

# Compilar la aplicación para producción (genera la carpeta dist/)
RUN npm run build

# Stage 2: Servir la aplicación con Nginx
FROM nginx:alpine

# Copiar la configuración personalizada de Nginx
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Copiar los archivos construidos de React al directorio de Nginx
COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
