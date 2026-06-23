# Guía de Despliegue con Docker - Golden Ticket

Este documento detalla cómo levantar e interactuar con la infraestructura dockerizada de **Golden Ticket** y justifica las decisiones arquitectónicas tomadas para esta configuración.

---

## Decisiones de Arquitectura

Para garantizar una arquitectura escalable, segura y portable, se optó por una estrategia multicontenedor utilizando **Docker Compose** con las siguientes características:

### 1. Compilación Multi-etapa (Multi-stage Builds)
Tanto en el `frontend/` como en el `backend/` se utiliza compilación multi-etapa. Esto separa las herramientas de desarrollo y construcción del entorno final de ejecución:
- **En el Backend (Go)**: La etapa de compilación instala las dependencias y construye el binario estático utilizando el compilador oficial de Go. La etapa final utiliza una imagen de **Alpine Linux** mínima que solo contiene el binario compilado. Esto reduce el tamaño de la imagen final a unos pocos megabytes (evitando incluir el compilador y código fuente) y minimiza la superficie de ataque de seguridad.
- **En el Frontend (React/Vite)**: Node.js se utiliza únicamente en la primera etapa para resolver dependencias de npm y compilar los archivos estáticos de producción (`dist/`). La segunda etapa descarta Node.js y copia estos archivos estáticos a una imagen ultra liviana de **Nginx**.

### 2. Frontend servido con Nginx
En entornos de producción, es una mala práctica levantar la aplicación frontend utilizando servidores de desarrollo como `npm run dev` o `vite`. Estos servidores no están optimizados para servir archivos estáticos a gran escala ni manejan concurrencia de forma eficiente. 
Nginx actúa como un servidor de producción de alto rendimiento. Además, en el archivo `nginx.conf` se configuró la directiva `try_files` para redirigir las peticiones que no corresponden a recursos estáticos hacia `index.html`, permitiendo que el enrutamiento del lado del cliente (React Router) funcione sin fallos de 404 al recargar el navegador.

### 3. Persistencia de Datos con Volúmenes de Docker
Los contenedores Docker son efímeros por naturaleza: al destruirse un contenedor, todos sus datos internos se pierden. Para evitar la pérdida de los eventos, compras e información de usuarios creados durante la ejecución del sistema, se configuró un volumen persistente de Docker llamado `db_data`, mapeado al directorio `/var/lib/mysql` dentro de la base de datos. De esta forma, los datos sobreviven a los comandos `docker compose down` y actualizaciones de la imagen.

### 4. Orquestación y Sincronización
La comunicación entre los componentes se realiza mediante una red virtual interna gestionada por Docker Compose. Para evitar fallos en el backend al intentar conectarse a la base de datos antes de que esta termine de inicializarse:
- Se implementó un **Healthcheck** en la base de datos usando `mysqladmin ping`.
- El servicio `backend` declara una dependencia estricta (`depends_on`) que espera a que el servicio `db` esté marcado como saludable (`condition: service_healthy`) antes de iniciar su ejecución.

---

## Instrucciones de Uso

### Requisitos Previos
- Tener instalado **Docker** y **Docker Compose**.

### 1. Levantar la Aplicación
Ejecutá el siguiente comando desde la raíz del proyecto para descargar las imágenes base, construir las imágenes locales de frontend y backend, y encender todos los servicios en segundo plano:

```bash
docker compose up --build -d
```

### 2. Verificar el Estado de los Servicios
Podés validar que los tres contenedores estén corriendo normalmente ejecutando:

```bash
docker compose ps
```

Deberías ver una salida similar a la siguiente:
```txt
NAME                      IMAGE               COMMAND                  SERVICE             CREATED             STATUS                    PORTS
golden_ticket_backend     golden-ticket-backend   "./server"               backend             10 seconds ago      Up 9 seconds              0.0.0.0:8080->8080/tcp
golden_ticket_db          mysql:8.0           "docker-entrypoint.s…"   db                  10 seconds ago      Up 9 seconds (healthy)    0.0.0.0:3306->3306/tcp, 33060/tcp
golden_ticket_frontend    golden-ticket-frontend  "/docker-entrypoint.…"   frontend            10 seconds ago      Up 9 seconds              0.0.0.0:3000->80/tcp
```

### 3. Acceder al Sistema
Una vez levantados los servicios, podés ingresar a las siguientes URLs en tu navegador:
- **Frontend (React)**: [http://localhost:3000](http://localhost:3000)
- **API Backend (Go)**: [http://localhost:8080/events](http://localhost:8080/events)

---

## Solución de Problemas (Troubleshooting)

### Puertos Ocupados
Si al levantar Docker Compose recibís un error indicando que un puerto ya está asignado (por ejemplo, el `3306` o el `8080`), asegurate de apagar cualquier servicio de base de datos MySQL local o instancias de Go/Node que se estén ejecutando en segundo plano en tu host.

### Ver los Logs en Tiempo Real
Si querés auditar qué está pasando con la migración de base de datos o el tráfico en la API, podés ver los logs de los contenedores usando:

```bash
docker compose logs -f
```

O inspeccionar un contenedor en particular (ej. el backend):
```bash
docker compose logs -f backend
```

### Borrar y Reiniciar la Base de Datos desde Cero
Si por algún motivo querés borrar todas las tablas y datos cargados para recrear el entorno desde cero, podés apagar los contenedores y eliminar el volumen persistente de datos corriendo:

```bash
docker compose down -v
```
Al volver a correr `docker compose up --build -d`, se creará una base de datos MySQL vacía y el backend ejecutará nuevamente las automigraciones automáticamente.
