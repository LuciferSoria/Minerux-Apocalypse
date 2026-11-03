# Backend privado — 30 comprobaciones de seguridad

El backend privado incluye `security_contract_test.go`, una suite de 30 casos HTTP para:

- separación entre tokens de administración, servidor y jugador;
- rechazo de solicitudes no autenticadas;
- límites de creación de invitaciones y validación de JSON;
- denegación de canjes inválidos;
- consulta privada de permisos y consulta pública del estado.

Las pruebas se ejecutan desde el directorio `backend` con `go test ./...` y `go vet ./...`.

El código fuente del backend y los secretos se distribuyen **únicamente en el ZIP privado**. La documentación pública no contiene códigos beta ni tokens.

**Límite:** estas pruebas validan contratos HTTP del backend, no la autenticación oficial de Microsoft ni el rechazo de conexiones dentro de NeoForge. Esas integraciones continúan pendientes.
