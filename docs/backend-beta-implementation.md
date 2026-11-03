# Backend beta — implementación inicial

Se ha preparado un prototipo de backend en Go 1.23 (módulo `backend`) con:

- `GET /api/v1/status`
- `PATCH /api/v1/admin/status`
- `POST /api/v1/admin/beta-codes` (generación por lotes, códigos aleatorios de 128 bits)
- `GET /api/v1/admin/beta-codes` (solo metadatos)
- `POST /api/v1/admin/beta-codes/{id}/revoke`
- `POST /api/v1/beta/redeem` (canje con token de jugador firmado)
- `GET /api/v1/me/access`

**Estado:** el código del backend, sus pruebas y una base de datos local con 25 códigos beta de un solo uso se han entregado al propietario en un paquete privado. **Los secretos, códigos originales y base de datos no deben añadirse a este repositorio público.**

El backend usa un token administrador de entorno y exige tokens HMAC de identidad de jugador emitidos por un servicio de autenticación de Minecraft de confianza. Dicho servicio emisor **aún no está implementado**. Los tokens nunca deben ser generados por el launcher. La lista de acceso al servidor NeoForge tampoco está integrada, así que los códigos aún no restringen conexiones reales.

Para despliegue de producción: TLS mediante proxy inverso, rate limiting por IP en el proxy, auditoría, gestión de claves, copias de seguridad, autenticación oficial de jugadores y autorización en servidor. El almacenamiento JSON está pensado para una sola instancia; no se debe ejecutar en múltiples réplicas.

Los 25 códigos entregados caducan a los 30 días desde su creación, y solo serán canjeables tras desplegar el backend con su estado persistente y sus secretos originales.
