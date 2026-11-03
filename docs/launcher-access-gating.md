# Control de acceso en el launcher DEV

El botón `JUGAR (EN DESARROLLO)` consulta `MINERUX_API_URL` (HTTPS).

- `maintenance`: muestra el aviso y no inicia el juego.
- `beta`: informa de que la autenticación y la validación beta todavía no están integradas.
- `online`: informa de disponibilidad, pero **no** inicia el juego porque aún falta implementar el lanzamiento directo.
- Sin API, URL inválida, error de red o respuesta desconocida: se bloquea el intento (fail-closed).

El botón `ABRIR MINECRAFT OFICIAL` es una utilidad independiente y **no** proporciona acceso al servidor Minerux. El launcher no puede impedir que un usuario abra Minecraft por otros medios.

**La protección real debe aplicarse en el servidor** mediante una lista de acceso o integración NeoForge con el backend. El servicio de identidad verificada de Minecraft sigue pendiente. No distribuir esta compilación DEV como launcher final.
