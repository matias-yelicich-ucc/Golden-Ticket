# Propuesta: Cancelación de Compra de Entrada (Cliente E)

Esta propuesta detalla el diseño y plan de implementación para la funcionalidad de cancelación de compra de entradas, contemplando el impacto en el cupo del evento y la interfaz del usuario.

## 1. Diseño del Backend

### 1.1. Nuevo Endpoint de Cancelación
Se implementará un endpoint protegido que permita a un usuario autenticado cancelar su entrada.
* **Ruta**: `POST /my-tickets/:id/cancel` (y opcionalmente `DELETE /my-tickets/:id` para compatibilidad)
* **Método**: `POST`
* **Autenticación**: Sí, requiere token JWT (usando `middleware.AuthMiddleware()`).
* **Parámetros**: `id` de la entrada (en la URL).

### 1.2. Lógica de Negocio y Liberación de Cupo
* El estado de la entrada cambiará de `activo` a `cancelado`.
* **Impacto en el Cupo**: El sistema calcula la disponibilidad de entradas restando al cupo máximo (`Capacidad`) solo las entradas en estado `activo` (ver `ticket_dao.go` y `event_service.go`). Al cambiar el estado a `cancelado`, esta entrada dejará de computar, liberando automáticamente el cupo del evento para futuras compras.
* **Validaciones**:
  1. Que la entrada exista.
  2. Que pertenezca al usuario autenticado.
  3. Que la entrada esté en estado `activo`.
  4. Que el evento no haya ocurrido ni esté en curso (opcional pero recomendado por consistencia de negocio).

### 1.3. Cambios en el Backend
* **DAO (`dao/ticket_dao.go`)**: Agregar `CancelTicket(userID uint, ticketID uint) error`.
* **Servicio (`services/ticket_service.go`)**: Agregar `CancelTicket(userID uint, ticketID uint) error`.
* **Controlador (`controllers/ticket_controller.go`)**: Agregar método `Cancel(c *gin.Context)`.
* **Rutas (`main.go`)**: Registrar la nueva ruta protegida.

## 2. Diseño del Frontend

### 2.1. Interfaz de "Mis Entradas"
* En la pantalla `MyTickets.jsx`, para cada entrada activa (`activo`), se agregará un botón o acción de **"Cancelar Entrada"**.
* Si la entrada ya está cancelada (`cancelada` o `cancelado`), se mostrará una etiqueta indicando el estado sin opción de cancelación.

### 2.2. Modal de Confirmación
* Al presionar "Cancelar Entrada", se abrirá un modal de confirmación con estética premium similar al de compra.
* El modal contendrá:
  - Título: **"Confirmar Cancelación"**
  - Mensaje: **"¿Estás seguro de que deseas cancelar esta entrada? Se realizará la devolución del dinero a tu medio de pago original."** (Confirmación de devolución de dinero solicitada).
  - Botones de acción: **"Confirmar Cancelación"** (rojo/estética de alerta) y **"Volver"** (gris/cancelar).
* Al confirmar, se realizará la llamada a la API del backend. Si es exitosa, se mostrará un mensaje de éxito, se cerrará el modal y se recargará el listado de entradas.

---
¿Continuamos con las especificaciones detalladas?
