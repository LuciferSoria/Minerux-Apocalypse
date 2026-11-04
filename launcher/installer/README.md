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

## Importación de mods desde ZIP local

La función Go `installer.ImportModsZIP(zipPath, root, version)` acepta un ZIP local con archivos `mods/*.jar` y prepara los mods en `versions/<version>/mods/`. Rechaza rutas anidadas, nombres de dispositivo reservados de Windows, duplicados, enlaces simbólicos y entradas demasiado grandes. Una segunda importación reutiliza los JAR cuyo contenido coincide y rechaza los modificados.

Se ha inspeccionado el archivo aportado `mods.zip`: contiene 12 JAR y un directorio `mods/`. Este ZIP no se publica automáticamente en el repositorio. La importación todavía no está conectada al botón de la interfaz, y no verifica por sí misma compatibilidad con NeoForge ni instala Minecraft.

## Botón de Windows: IMPORTAR MODS.ZIP

La interfaz de desarrollo incluye el botón **IMPORTAR MODS.ZIP**. Busca un archivo llamado `mods.zip` junto al `.exe` y llama a `ImportModsZIP` para preparar sus archivos `mods/*.jar` en `%LOCALAPPDATA%\\MineruxApocalypse\\versions\\26.1.2\\mods\\`.

No descarga el ZIP, no modifica la carpeta oficial de Minecraft y no garantiza que los mods sean compatibles con NeoForge. Para que el botón funcione hay que proporcionar el ZIP localmente; el repositorio no incluye los mods. El botón tampoco instala Minecraft ni ejecuta el juego.
