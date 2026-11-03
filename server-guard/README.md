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
