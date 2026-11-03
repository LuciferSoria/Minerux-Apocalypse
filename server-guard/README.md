# Integración servidor — fase inicial

Este directorio contiene un **adaptador Java de servidor** que consulta la API privada del backend para decidir si un UUID verificado puede entrar.

- Solo se debe invocar con el UUID obtenido de la autenticación online-mode del servidor Minecraft.
- Configurar URL HTTPS y MINERUX_SERVER_TOKEN **en el servidor**, nunca en el cliente.
- Si la API falla, responde false y bloquea el acceso (fail closed).
- Ejecutar en un hilo de trabajo, nunca bloquear el tick principal del servidor.
- La respuesta HTTP esperada es el JSON del backend: `{"allowed":true,"status":"beta"}`.

**IMPORTANTE:** todavía NO es un mod NeoForge instalable. Falta integrarlo con el evento de login/validación de NeoForge 26.1.2.125, compilarlo contra las APIs correctas y probarlo en el servidor. Añadir un archivo Java no protege las conexiones por sí solo.

Prueba local con Java 21+:

```sh
mkdir -p out
javac -d out src/main/java/dev/minerux/guard/BetaAccessClient.java src/test/java/dev/minerux/guard/BetaAccessClientTest.java
java -cp out dev.minerux.guard.BetaAccessClientTest
```

## Nuevo límite de integración asíncrona

`LoginAccessGate.check(UUID)` devuelve `CompletableFuture<Boolean>` y ejecuta la consulta HTTP en el `Executor` proporcionado. Si no hay UUID, el ejecutor falla o la consulta genera un error, el resultado es `false`. **No usar `join()` ni `get()` en el hilo principal de Minecraft.** El futuro debe resolverse durante una fase de login que permita retrasar la aceptación; si NeoForge no ofrece esa fase, se necesita una estrategia segura para mantener al jugador desconectado hasta la verificación. Todavía falta registrar el evento NeoForge concreto y verificarlo con la versión 26.1.2.125.

La autenticación online-mode debe estar activa. Una UUID suministrada por el cliente no es prueba de identidad.

## Diagnóstico desde el servidor (sin instalar NeoForge)

El comando `ServerGuardPreflight` permite verificar desde la máquina del servidor que el backend HTTPS reconoce un UUID **ya verificado por Minecraft online-mode**. Requiere Java 21+ y las variables de entorno `MINERUX_API_URL` y `MINERUX_SERVER_TOKEN`. No se debe ejecutar con una UUID suministrada por un cliente sin autenticar.

Compilar las clases Java de `src/main/java/dev/minerux/guard/` y ejecutar:

```sh
java -cp out dev.minerux.guard.ServerGuardPreflight 01234567-89ab-cdef-0123-456789abcdef
```

Códigos de salida: `0` = autorizado; `2` = denegado o backend inaccesible; `3` = uso o configuración inválidos. El comando no imprime el token ni la URL del backend.

**Este diagnóstico no protege el servidor por sí solo**. La integración con la fase de login de NeoForge sigue pendiente.
