# Excel to DB Backend

API REST en Go para gestionar calificaciones académicas y subir archivos Excel a una base de datos MySQL. Diseñada como backend para aplicaciones Flutter, ideal para portafolio y proyectos educativos.

## Stack Tecnológico

| Tecnología       | Versión | Propósito                        |
| ---------------- | ------- | -------------------------------- |
| Go               | 1.25+   | Lenguaje de programación         |
| Gorilla Mux      | 1.8     | Router HTTP                      |
| Gorilla Handlers | 1.5     | Middleware CORS                  |
| GORM             | 1.31    | ORM para base de datos           |
| MySQL            | 8+      | Base de datos relacional         |
| Air              | -       | Live-reload en desarrollo        |

## Funcionalidades

- **CRUD de calificaciones** — Crear, listar, obtener, eliminar registros de estudiantes con sus notas en múltiples materias.
- **Carga de Excel** — Subir estructura de archivos Excel (hojas, encabezados, filas) como JSON para almacenar en DB.
- **Listado de subidas** — Consultar el historial de archivos Excel subidos con sus hojas y datos.
- **CORS configurable** — Middleware CORS para permitir peticiones desde cualquier origen (configurable vía `CORS_ORIGIN`).

## Requisitos previos

- Go 1.25 o superior ([descargar](https://go.dev/dl/))
- MySQL 8+ corriendo localmente o en la nube
- (Opcional) [Air](https://github.com/air-verse/air) para live-reload

## Estructura del proyecto

```
Excel_to_db_backend-main/
├── db/
│   └── conection.go        # Conexión a MySQL vía GORM
├── models/
│   ├── califications.go    # Modelo Cal (calificaciones)
│   ├── upload.go           # Modelos Upload, Sheet, Row
│   └── models_test.go      # Tests unitarios de modelos
├── routes/
│   ├── index.routes.go     # Ruta raíz (health check)
│   ├── people.route.go     # CRUD de calificaciones
│   ├── excel.route.go      # Subida y consulta de Excel
│   ├── califications.route.go
│   └── routes_test.go      # Tests unitarios de rutas
├── main.go                 # Punto de entrada
├── go.mod / go.sum         # Dependencias
├── .air.toml               # Configuración de Air
├── railway.json            # Configuración para Railway
├── .env.example            # Variables de entorno de ejemplo
└── README.md
```

## Instalación y ejecución local

### 1. Clonar el repositorio

```bash
git clone https://github.com/tu-usuario/Excel_to_db_backend.git
cd "Excel_to_db_backend-main"
```

### 2. Configurar variables de entorno

Copia el archivo de ejemplo y ajústalo:

```bash
cp .env.example .env
```

Edita `.env` con los datos de tu base de datos MySQL:

```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASS=tu_contraseña
DB_NAME=test_web_excel1
```

### 3. Crear la base de datos

Conéctate a MySQL y crea la base de datos:

```sql
CREATE DATABASE IF NOT EXISTS test_web_excel1;
```

La aplicación creará las tablas automáticamente al iniciar (via `AutoMigrate`).

### 4. Ejecutar la aplicación

```bash
go run .
```

O con live-reload (si tienes Air instalado):

```bash
air
```

El servidor iniciará en `http://localhost:3000`.

## Variables de Entorno

| Variable       | Default           | Descripción                                    |
| -------------- | ----------------- | ---------------------------------------------- |
| `DB_HOST`      | `localhost`       | Host de MySQL                                  |
| `DB_PORT`      | `3306`            | Puerto de MySQL                                |
| `DB_USER`      | `root`            | Usuario de MySQL                               |
| `DB_PASS`      | (vacío)           | Contraseña de MySQL                            |
| `DB_NAME`      | `test_web_excel1` | Nombre de la base de datos                     |
| `MYSQL_HOST`   | → `DB_HOST`       | Usado automáticamente por Railway MySQL        |
| `MYSQL_PORT`   | → `DB_PORT`       | Usado automáticamente por Railway MySQL        |
| `MYSQL_USER`   | → `DB_USER`       | Usado automáticamente por Railway MySQL        |
| `MYSQL_PASSWORD` | → `DB_PASS`     | Usado automáticamente por Railway MySQL        |
| `MYSQL_DATABASE` | → `DB_NAME`     | Usado automáticamente por Railway MySQL        |
| `PORT`         | `3000`            | Puerto del servidor HTTP                       |
| `CORS_ORIGIN`  | `*`               | Origen permitido para CORS                     |

> **Nota para Railway:** cuando agregas el plugin MySQL de Railway, este provee automáticamente las variables `MYSQL_*`. El código las detecta primero y, si no existen, usa las `DB_*` como fallback.

## API Endpoints

### Health Check

```
GET /
```

Respuesta: `Hello Go, This is MY first API rest 40`

### Calificaciones

| Método | Ruta              | Descripción                        |
| ------ | ----------------- | ---------------------------------- |
| GET    | `/names`          | Lista todas las calificaciones     |
| GET    | `/names/{id}`     | Obtiene una calificación por ID    |
| POST   | `/names`          | Crea una calificación              |
| POST   | `/names/all`      | Crea múltiples calificaciones      |
| DELETE | `/names/{id}`     | Elimina una calificación por ID    |
| DELETE | `/names`          | Elimina todas las calificaciones   |

#### Ejemplo POST /names

```json
{
  "nombre": "Juan Pérez",
  "math": "90",
  "physical": "85",
  "chemistry": "88",
  "biologi": "92",
  "histori": "76",
  "geografi": "81",
  "literature": "95",
  "spanish": "87",
  "english": "93"
}
```

### Excel

| Método | Ruta              | Descripción                            |
| ------ | ----------------- | -------------------------------------- |
| POST   | `/excel/upload`   | Sube la estructura de un archivo Excel |
| GET    | `/excel/uploads`  | Lista todas las subidas realizadas     |

#### Ejemplo POST /excel/upload

```json
{
  "file": "notas_2024.xlsx",
  "sheets": [
    {
      "sheet": "1er Trimestre",
      "hasHeaders": true,
      "headers": ["nombre", "math", "physical"],
      "rows": [
        { "nombre": "Juan", "math": "90", "physical": "85" },
        { "nombre": "María", "math": "95", "physical": "88" }
      ]
    }
  ]
}
```

## Despliegue en Railway (Free Tier)

Railway ofrece una capa gratuita con $5 de crédito mensual, suficiente para correr esta aplicación junto con una base de datos MySQL pequeña.

### Paso 1: Crear cuenta en Railway

Ve a [railway.app](https://railway.app) e inicia sesión con GitHub.

### Paso 2: Subir el código a GitHub

```bash
git init
git add .
git commit -m "Initial commit"
git remote add origin https://github.com/tu-usuario/Excel_to_db_backend.git
git push -u origin main
```

### Paso 3: Crear proyecto en Railway

1. Haz clic en **New Project** → **Deploy from GitHub repo**.
2. Selecciona el repositorio `Excel_to_db_backend`.
3. Railway detectará automáticamente que es un proyecto Go y lo construirá con Nixpacks.

### Paso 4: Agregar base de datos MySQL

1. Dentro del proyecto en Railway, haz clic en **New** → **Database** → **Add MySQL**.
2. Railway creará un plugin MySQL y expondrá automáticamente las variables de entorno `MYSQL_HOST`, `MYSQL_PORT`, `MYSQL_USER`, `MYSQL_PASSWORD`, `MYSQL_DATABASE`.
3. El código ya está preparado para leer estas variables, por lo que no necesitas configurar nada adicional.

> 💡 **Optimización de costos:** El plugin MySQL de Railway cuesta aproximadamente $5/mes, lo que consume todo el crédito gratuito. La aplicación en sí no tiene costo adicional. Si prefieres usar una base de datos externa gratuita, puedes usar servicios como [TiDB Cloud](https://tidbcloud.com) (MySQL-compatible, free tier) o [Aiven](https://aiven.io) y configurar las variables `DB_*` manualmente.

### Paso 5: Configurar variable PORT

Railway asigna automáticamente el puerto. No necesitas hacer nada, el código usa `getEnv("PORT", "3000")`.

### Paso 6: Desplegar

Railway desplegará automáticamente. Cada vez que hagas push a la rama `main`, se redeployará automáticamente.

Tu aplicación quedará accesible en una URL como: `https://excel-to-db-backend.up.railway.app`

## Conectar el Frontend

Tu frontend (Flutter, React, etc.) solo necesita la URL base del servidor desplegado.

### En desarrollo local

```dart
final String baseUrl = 'http://localhost:3000';
```

### En producción (Railway)

```dart
final String baseUrl = 'https://tu-proyecto.up.railway.app';
```

Reemplaza `https://tu-proyecto.up.railway.app` con la URL que Railway te asigne (la encuentras en el dashboard de Railway, en la sección **Settings** → **Domains**).

### Ejemplo de uso con Flutter

```dart
import 'dart:convert';
import 'package:http/http.dart' as http;

class ApiService {
  final String baseUrl;

  ApiService(this.baseUrl);

  Future<List<dynamic>> getCalificaciones() async {
    final response = await http.get(Uri.parse('$baseUrl/names'));
    if (response.statusCode == 200) {
      return json.decode(response.body);
    }
    throw Exception('Error al obtener calificaciones');
  }

  Future<Map<String, dynamic>> createCalificacion(Map<String, dynamic> data) async {
    final response = await http.post(
      Uri.parse('$baseUrl/names'),
      headers: {'Content-Type': 'application/json'},
      body: json.encode(data),
    );
    if (response.statusCode == 200) {
      return json.decode(response.body);
    }
    throw Exception('Error al crear calificación');
  }
}
```

## Tests

### Ejecutar tests unitarios

```bash
go test ./...
```

### Ejecutar tests con verbose

```bash
go test -v ./...
```

### Ver cobertura

```bash
go test -cover ./...
```

Los tests cubren:

- Validación de datos (`validateCal`, `isJSONContent`)
- Handlers HTTP (códigos de respuesta, errores de validación)
- Modelos y sus relaciones
- Rutas de error (JSON inválido, campos faltantes, types incorrectos)

> Los handlers que requieren conexión a base de datos se prueban hasta el punto de validación de entrada. Para pruebas de integración completas, se recomienda usar una base de datos de prueba local.

## Licencia

MIT
