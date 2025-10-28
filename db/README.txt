===============================================================================
IMPORTANTE: SISTEMA DE MIGRACIONES AUTOMÁTICAS
===============================================================================

Los archivos SQL en este directorio son LEGACY (solo referencia histórica).

NO ejecutar manualmente estos archivos. El sistema ahora usa migraciones
automáticas que se ejecutan al iniciar la aplicación.

-------------------------------------------------------------------------------
PARA GESTIONAR LA BASE DE DATOS:
-------------------------------------------------------------------------------

1. Migraciones automáticas (desarrollo):
   - Configurar en config.json: "AUTO_MIGRATE": true
   - Ejecutar: ./main
   - Las migraciones se aplican automáticamente

2. Control manual de migraciones:
   - Compilar CLI: go build -o migrate cmd/migrate/main.go
   - Ver versión: ./migrate -version
   - Aplicar: ./migrate -up
   - Revertir: ./migrate -down

3. Crear nueva migración:
   - Crear: internal/migration/migrations/000015_nombre.up.sql
   - Crear: internal/migration/migrations/000015_nombre.down.sql
   - Recompilar: go build -o main cmd/main.go

-------------------------------------------------------------------------------
DOCUMENTACIÓN COMPLETA:
-------------------------------------------------------------------------------

Ver: docs/manual_programador.typ
     Sección: "Base de Datos > Sistema de Migraciones"

Incluye:
- Guía de migraciones
- Mejores prácticas SQL
- Convenciones de nomenclatura
- Ejemplos de uso

===============================================================================
