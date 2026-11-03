# Revocación efectiva de acceso beta

**Advertencia operativa:** revocar un código ya canjeado no revoca automáticamente el permiso de un jugador. El backend guarda los permisos en `Grants[uuid]`, por separado del código original. La revocación del código solo impide futuros canjes.

Antes de distribuir la beta, el panel de administración necesita una acción explícita de revocar el acceso de un UUID verificado, persistir el cambio y comprobar que `GET /api/v1/internal/player-access/{uuid}` devuelve `allowed:false` en modo beta.

Además, el servidor NeoForge debe consultar el backend en cada entrada; si se desea expulsión inmediata de jugadores conectados, se necesita una notificación o comprobación periódica adicional. No confundir bloqueo de nuevos accesos con expulsión inmediata.

No publicar el archivo privado de códigos, el token de administración ni el token de servidor en GitHub.

El backend privado actualizado incorpora `POST /api/v1/admin/players/{uuid}/revoke` con token administrativo; revoca el permiso persistente de ese UUID. Este endpoint no expulsa automáticamente a jugadores ya conectados.
