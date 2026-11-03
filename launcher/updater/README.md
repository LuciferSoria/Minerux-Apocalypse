# Actualizador seguro (componente)

El paquete `updater` valida manifiestos JSON, rutas relativas, URLs HTTPS y SHA-256. Descarga a un archivo temporal, comprueba el hash y solo entonces reemplaza el archivo de destino.

**No está integrado todavía con el botón Instalar/Jugar.** Los manifiestos deben provenir de una fuente controlada y fiable; HTTPS y SHA-256 por sí solos no autentican un manifiesto manipulado. La descarga tiene un límite de tiempo y un límite de tamaño por archivo de aproximadamente 1 GiB; la protección frente a enlaces simbólicos concurrentes requiere endurecimiento adicional antes de un lanzamiento público.

No incluir mods sin autorización de distribución. No usar el manifiesto de ejemplo como un modpack real.
