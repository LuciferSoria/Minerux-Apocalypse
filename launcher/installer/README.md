# Paquete local del launcher (experimental)

El botón **INSTALAR PAQUETE VERIFICADO** busca `minerux-package.json` en el mismo directorio que el ejecutable. Si falta, no descarga nada.

El archivo JSON debe tener la forma:

```json
{
  "version": "26.1.2-beta1",
  "files": [
    {
      "name": "ejemplo.jar",
      "url": "https://downloads.example.org/ejemplo.jar",
      "sha256": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
    }
  ]
}
```

**El ejemplo es ilustrativo y no contiene un mod real ni un hash válido para descargar.** No distribuirlo como manifiesto de producción. Los SHA-256 deben calcularse de los archivos oficiales o de los mods autorizados y verificarse antes de publicar el manifiesto.

El launcher instala los archivos en `%LOCALAPPDATA%\\MineruxApocalypse\\versions\\<version>\\`. Rechaza rutas peligrosas, nombres duplicados, redirecciones y descargas cuyo hash no coincida.

## Limitaciones actuales

- Este botón **no** instala Minecraft ni NeoForge, ni inicia el juego.
- No hay manifiesto de producción firmado o incorporado al launcher.
- El manifiesto local se considera confiable: no aceptar archivos JSON de terceros desconocidos.
- La instalación por versión todavía no es transaccional: un fallo puede dejar archivos ya descargados en la carpeta de esa versión. No se activa automáticamente ninguna versión.
- Los archivos existentes se reutilizan únicamente si su SHA-256 coincide con el manifiesto. Si un archivo está dañado o ha sido modificado, la instalación se detiene sin sobrescribirlo; para repararlo habrá que eliminarlo manualmente o implementar una reparación controlada.
- La funcionalidad debe probarse en Windows antes de publicar un ejecutable.

La compilación de desarrollo no constituye una versión jugable.
