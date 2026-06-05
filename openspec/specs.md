# Especificaciones: Cancelación de Compra de Entrada (Cliente E)

## 1. Especificaciones de la API

El nuevo endpoint permitirá a los clientes autenticados dar de baja una entrada específica de forma lógica.

### 1.1. Detalle del Endpoint
* **Ruta**: `POST /my-tickets/:id/cancel`
* **Método**: `POST`
* **Cabeceras**:
  * `Authorization: Bearer <JWT_TOKEN>`

### 1.2. Códigos de Estado y Respuestas

#### Respuesta Exitosa (200 OK)
Retornado cuando la entrada es cancelada con éxito.
* **Cuerpo (JSON)**:
  ```json
  {
    "message": "Entrada cancelada con éxito y cupo liberado"
  }
  ```

#### Errores de Cliente (400 / 401 / 403 / 404)
* **400 Bad Request**: El ID provisto no es un número válido, el ticket ya está cancelado, o el evento ya ocurrió/está en curso.
  ```json
  {
    "error": "la entrada ya se encuentra cancelada"
  }
  ```
* **401 Unauthorized**: Token JWT faltante, inválido o expirado.
  ```json
  {
    "error": "Usuario no autenticado"
  }
  ```
* **403 Forbidden**: El usuario autenticado intenta cancelar una entrada que no le pertenece.
  ```json
  {
    "error": "no eres el propietario de esta entrada"
  }
  ```
* **404 Not Found**: La entrada con el ID especificado no existe.
  ```json
  {
    "error": "entrada no encontrada"
  }
  ```

---

## 2. Reglas de Validación
1. **Verificación de Propiedad**: `ticket.UserID` debe ser igual al `userID` extraído del token JWT.
2. **Verificación de Estado del Ticket**: `ticket.Estado` debe ser `"activo"`.
3. **Verificación Temporal del Evento**: No se puede cancelar una entrada si el evento ya ocurrió o está en curso (es decir, la fecha/hora de inicio del evento es anterior a la hora actual).

---

## 3. Especificaciones del Frontend

### 3.1. Listado de Entradas ("Mis Entradas")
* Se renderizará un botón de "Cancelar Entrada" al lado de cada entrada en estado `activo`.
* Si el estado es `cancelado`, se mostrará una etiqueta con estilo deshabilitado/alerta que diga `"Cancelada"`.

### 3.2. Modal de Cancelación y Devolución
* El modal tendrá la misma estructura visual y estética que el de compra de entradas.
* Mostrará un mensaje claro de reembolso:
  > **Confirmación de Devolución**
  >
  > ¿Estás seguro de que deseas cancelar tu entrada para el evento **[Nombre del Evento]**?
  >
  > Se reembolsará el importe de **$[Precio]** a tu medio de pago original de forma automática.
* Botón **"Confirmar Cancelación"** (color rojo/alerta).
* Botón **"Volver"** (color secundario/gris).
* Manejo de estado de carga: deshabilitar botones y mostrar loader mientras la API responde.
* Mensaje de éxito/error: mostrar notificación o feedback visual tras completar la operación.
