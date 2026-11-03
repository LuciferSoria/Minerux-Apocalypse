# Panel de administración y acceso beta

## Objetivo
El propietario de Minerux Apocalypse podrá gestionar el estado del servicio y el acceso a la beta sin recompilar el launcher.

## Roles
- **Administrador:** inicia sesión mediante autenticación segura; puede cambiar el estado del servicio, crear, listar y revocar invitaciones beta.
- **Jugador:** consulta el estado público, introduce un código beta y juega únicamente cuando su cuenta tiene acceso.

## Estados públicos
- `online`: juego disponible.
- `maintenance`: el launcher deshabilita JUGAR y muestra el mensaje del administrador.
- `beta`: acceso sujeto a una autorización beta vigente.

Campos públicos: `status`, `message`, `estimated_return` (opcional), `updated_at`, `minimum_launcher_version` (opcional). **Nunca** publicar secretos de administración ni códigos beta completos.

## Modelo de invitaciones beta
- Código aleatorio criptográficamente seguro, al menos 128 bits de entropía, mostrado una sola vez.
- Almacenar solo un hash seguro del código; nunca guardar códigos en texto plano.
- `max_redemptions`, `redeemed_count`, `expires_at`, `revoked_at` y notas internas.
- Canje atómico y con límite de intentos para impedir carreras y ataques por fuerza bruta.
- Vincular el acceso canjeado a una identidad de Minecraft verificada, no al nombre escrito manualmente ni solo al dispositivo.
- Registrar acciones administrativas sin guardar códigos completos en los logs.

## API propuesta
- `GET /api/v1/status` pública, con caché corta.
- `POST /api/v1/beta/redeem` requiere identidad de jugador verificada y código.
- `GET /api/v1/me/access` requiere identidad verificada.
- `PATCH /api/v1/admin/status` requiere sesión de administrador.
- `POST /api/v1/admin/beta-codes` requiere sesión de administrador.
- `GET /api/v1/admin/beta-codes` devuelve metadatos, nunca códigos originales.
- `POST /api/v1/admin/beta-codes/{id}/revoke` requiere sesión de administrador.

## Seguridad y aplicación real
El launcher es un cliente no confiable: el bloqueo del botón JUGAR es solo UX. La **lista de acceso/whitelist o plugin/mod de servidor** debe verificar las autorizaciones para impedir conexiones directas sin código. El backend debe usar HTTPS, control de permisos del lado servidor, protección contra CSRF si hay cookies, auditoría, rate limiting y almacenamiento persistente. No incorporar una contraseña de administrador dentro del ejecutable ni confiar en una bandera local para conceder acceso.

## Orden de implementación
1. Definir contrato de estados y clientes públicos.
2. Implementar backend persistente con cuentas administrativas.
3. Añadir generador y canje seguro de códigos beta.
4. Integrar el panel privado y las pantallas de estado/canje en el launcher.
5. Integrar el control de acceso en el servidor NeoForge y hacer pruebas extremo a extremo.

**Estado:** especificación de arquitectura; no representa un backend o panel ya implementado.
