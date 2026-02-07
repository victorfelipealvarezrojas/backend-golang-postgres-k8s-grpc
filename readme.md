## Base de Datos

### Migraciones con golang-migrate

`golang-migrate` gestiona cambios en el esquema de la BD (agnóstico a lenguaje y BD).

#### Instalación
```bash
brew install golang-migrate
```

#### Crear una migración
```bash
migrate create -ext sql -dir ./db/migrations -seq nombre_migracion
```

Crea dos archivos: `.up.sql` (aplicar) y `.down.sql` (revertir).

#### Ejecutar migraciones (esto es informativo... ejecutar desde makefile)
```bash
migrate -path db/migrations -database "postgresql://admin:pass12345@localhost:5432/simple_bank?sslmode=disable" up
```

#### Revertir migraciones (esto es informativo... ejecutar desde makefile)
```bash
migrate -path db/migrations -database "postgresql://admin:pass12345@localhost:5432/simple_bank?sslmode=disable" down
```

### Generar código Go con sqlc

`sqlc` genera código Go tipado y seguro desde queries SQL.

#### Instalación
```bash
brew install kyleconroy/sqlc/sqlc
```

O con Go:
```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

**instalacion con Go, agrega al PATH:**
```bash
# Edita ~/.zshrc (o ~/.bash_profile si usas bash)
nano ~/.zshrc

# Agrega al final del archivo:
export PATH=$PATH:$(go env GOPATH)/bin

# Guarda (Ctrl+O, Enter, Ctrl+X) y recarga:
source ~/.zshrc

# Verifica la instalación:
sqlc version
```

#### Inicializar sqlc

Genera `sqlc.yaml` automáticamente:
```bash
sqlc init
```

Luego edita `sqlc.yaml`:
```yaml
version: "2"
sql:
  - engine: "postgresql"
    queries: "./db/query"
    schema: "./db/migrations"
    gen:
      go:
        package: "db"
        out: "db/sqlc"
        emit_json_tags: true
```

#### Generar código
```bash
sqlc generate
```

Lee queries en `./db/query` y genera código en `./db/sqlc`.

### Tareas Makefile
```bash
make createdb    # Crear y levantar BD
make dropdb      # Eliminar BD y volúmenes
make migrate-up  # Ejecutar migraciones
make migrate-down # Revertir migraciones
make sqlc        # Generar codigo de bd (convierte las querys a codigo go)
```

### Recursos

- [golang-migrate](https://github.com/golang-migrate/migrate)
- [sqlc](https://sqlc.dev)

### Inicialización del proyecto

Crear módulo Go:
```bash
go mod init github.com/valvarez/simplebank
```

Descargar dependencias necesarias usadas por el modelo:
```bash
go mod tidy
```
### intalacion Driver de postgres
```bash
go get github.com/lib/pq

go mod tidy 
```

### Herramienta de testing (testify)
```bash
go get github.com/stretchr/testify
```
