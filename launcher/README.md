# Minerux Apocalypse — estado del launcher

Este directorio está reservado para el código fuente del launcher de Windows.

**No hay un ejecutable final funcional.** El prototipo anterior solo registra el servidor en la lista multijugador y abre el launcher oficial; NO instala Minecraft/NeoForge/mods, NO realiza autenticación Microsoft propia ni conecta automáticamente. Por tanto, no debe publicarse como versión final.

## Criterios obligatorios antes de una Release

1. Autenticación Microsoft oficial con OAuth, sin almacenar contraseñas.
2. Instalación verificable de Minecraft Java y NeoForge con versiones compatibles.
3. Descarga de mods con licencia adecuada, hashes y actualización segura.
4. Inicio real del juego con conexión al servidor y mensajes de error claros.
5. UI fluida, progreso, cancelación y recuperación de descargas.
6. Pruebas en Windows 10 y 11 x64 en un equipo limpio.

El workflow de GitHub Actions está preparado para compilar código Go cuando se añada `launcher/main.go`. Un artefacto DEV no implica que estas pruebas hayan pasado.
