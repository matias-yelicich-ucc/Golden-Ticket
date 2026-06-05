# Tareas: Cancelación de Compra de Entrada (Cliente E)

Este plan de tareas guiará la implementación incremental y el proceso de pruebas.

## Fase 1: Backend (TDD - Test Driven Development)

- [ ] **Tarea 1.1: Modificar Interfaces y Structs**
  - Agregar `CancelTicket(userID uint, ticketID uint) error` en `dao.TicketDAO` y `services.TicketService`.
  - Implementar la firma de los métodos en `dao/ticket_dao.go` y `services/ticket_service.go` (retornando `nil` provisionalmente).
- [ ] **Tarea 1.2: Agregar Tests Unitarios y de Integración en el Backend**
  - En `backend/controllers/ticket_controller_test.go`, agregar un nuevo bloque de pruebas que verifique la llamada a `POST /my-tickets/:id/cancel` con diferentes escenarios:
    - Cancelación exitosa.
    - Error: entrada no encontrada.
    - Error: no eres el propietario de esta entrada.
    - Error: entrada ya cancelada.
    - Error: el evento ya ocurrió o está en curso.
- [ ] **Tarea 1.3: Implementar la Lógica en DAO**
  - Desarrollar la lógica completa de `CancelTicket` en `dao/ticket_dao.go`.
- [ ] **Tarea 1.4: Implementar la Lógica en Service y Controller**
  - Implementar la delegación en `services/ticket_service.go`.
  - Implementar el método `Cancel` en `controllers/ticket_controller.go`.
  - Registrar la ruta en `backend/main.go`.
- [ ] **Tarea 1.5: Ejecutar Pruebas de Backend**
  - Correr `go test ./...` y verificar que todas las pruebas pasen (incluyendo las nuevas).

## Fase 2: Frontend

- [ ] **Tarea 2.1: Modificar API Client**
  - Agregar la función `cancelTicket(ticketId)` en `frontend/src/services/api/client.js`.
- [ ] **Tarea 2.2: Modificar Componente "Mis Entradas" (`MyTickets.jsx`)**
  - Agregar estado `selectedTicketForCancel` and `showCancelModal` (o similar).
  - En el mapeo de tickets, agregar el botón "Cancelar Entrada" si `ticket.estado == "activo"`.
  - Mostrar una etiqueta `"Cancelada"` si `ticket.estado == "cancelado"`.
- [ ] **Tarea 2.3: Implementar Modal de Confirmación**
  - Diseñar y renderizar el modal de devolución de dinero con el mensaje solicitado.
  - Implementar la lógica para realizar el llamado HTTP, mostrar el estado de carga y manejar el éxito/error.
- [ ] **Tarea 2.4: Agregar Estilos Premium**
  - Asegurar la consistencia estética usando los colores y clases CSS definidos en el sistema.

## Fase 3: Verificación Final

- [ ] **Tarea 3.1: Probar Flujo Extremo a Extremo**
  - Levantar base de datos y servidores, comprar una entrada, verla en "Mis Entradas", cancelarla y comprobar que el cupo disponible del evento aumente en 1 unidad.
