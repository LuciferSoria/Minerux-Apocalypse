# Control administrativo de permisos de jugadores

El backend privado ofrece dos operaciones de administración, ambas protegidas por el token de administrador:

- `POST /api/v1/admin/players/{uuid}/grant`: concede acceso beta a un UUID autenticado.
- `POST /api/v1/admin/players/{uuid}/revoke`: retira el permiso beta de ese UUID.

El UUID se expresa en 32 caracteres hexadecimales minúsculos, sin guiones. Las dos operaciones persisten el cambio en el archivo de estado.

El servidor consulta `GET /api/v1/internal/player-access/{uuid}` con su propio token. No se deben introducir tokens de administración o de servidor en el launcher ni en mods del cliente.

**Límite:** estas operaciones no verifican por sí solas la identidad del jugador. El UUID debe proceder del servidor Minecraft en online-mode o de un proveedor de identidad confiable. El mod NeoForge aún no está desplegado.
