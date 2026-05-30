# Excel to MySQL — API REST Backend

<p align="center">
  <strong>API REST en Go para almacenar archivos Excel en MySQL</strong>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go&logoColor=white" alt="Go">
  <img src="https://img.shields.io/badge/MySQL-8+-4479A1?style=flat&logo=mysql&logoColor=white" alt="MySQL">
  <img src="https://img.shields.io/badge/GORM-1.31-00ADD8?style=flat" alt="GORM">
  <img src="https://img.shields.io/badge/License-MIT-green" alt="MIT">
</p>

---

## Descripción

API REST construida en **Go** que permite recibir la estructura de archivos Excel (hojas, encabezados, filas) vía JSON y almacenarla en una base de datos **MySQL**. Incluye CRUD de registros genéricos y endpoints para subir/consultar archivos Excel procesados.

Diseñada como backend para aplicaciones **móviles** (Android, Flutter) y **web** (React, Vue, HTML vanilla).

---

## Arquitectura

```
┌─────────────┐     HTTP/JSON     ┌──────────────┐     SQL      ┌─────────┐
│  Frontend   │ ────────────────> │  API Go      │ ──────────> │  MySQL  │
│  (Web/App)  │ <──────────────── │  (Gorilla)   │ <────────── │   DB    │
└─────────────┘    JSON Response  └──────────────┘   GORM ORM  └─────────┘
                                                              │
                                                         ┌────┴────┐
                                                         │ uploads │
                                                         │ sheets  │
                                                         │  rows   │
                                                         │ records │
                                                         └─────────┘
```

---

## Stack Tecnológico

| Componente       | Tecnología       | Versión | Propósito                        |
| ---------------- | ---------------- | ------- | -------------------------------- |
| Lenguaje         | Go               | 1.25+   | Backend                          |
| Router           | Gorilla Mux      | 1.8.1   | Enrutamiento HTTP                |
| CORS             | Gorilla Handlers | 1.5.2   | Middleware CORS                   |
| ORM              | GORM             | 1.31.1  | Mapeo objeto-relacional          |
| Base de datos    | MySQL            | 8+      | Almacenamiento relacional        |
| Live-reload      | Air              | -       | Recarga en desarrollo            |
| Contenedor       | Docker           | -       | Despliegue portátil              |

---

## Estructura del Proyecto

```
Excel_to_db_backend-main/
├── main.go                    # Punto de entrada, rutas, CORS, servidor
├── db/
│   └── conection.go           # Conexión a MySQL vía GORM
├── models/
│   ├── califications.go       # Modelo Record (registros genéricos)
│   ├── upload.go              # Modelos Upload, Sheet, Row (Excel)
│   └── models_test.go         # Tests unitarios de modelos
├── routes/
│   ├── index.routes.go        # Health check (GET /)
│   ├── people.route.go        # CRUD de registros (/names)
│   ├── excel.route.go         # Subida y consulta de Excel
│   └── routes_test.go         # Tests unitarios de rutas
├── docker-compose.yml         # Orquestación MySQL + App
├── Dockerfile                 # Imagen Docker multi-stage
├── .env                       # Variables de entorno (no commitear)
├── .env.example               # Plantilla de variables de entorno
├── go.mod / go.sum            # Dependencias Go
├── .air.toml                  # Configuración de Air (live-reload)
├── railway.json               # Configuración para Railway
└── README.md                  # Este archivo
```

---

## Diagrama de Base de Datos

```
┌──────────────────┐
│     uploads      │
├──────────────────┤
│ id          (PK) │
│ created_at      │
│ updated_at      │
│ deleted_at      │
│ file_name       │──── 1
└──────────────────┘
                    │
                    │ N
              ┌─────┴──────┐
              │   sheets   │
              ├────────────┤
              │ id     (PK)│
              │ upload_id  │──── FK → uploads.id
              │ sheet_name │
              │ has_headers│
              │ headers    │ (JSON: ["col1","col2",...])
              └─────┬──────┘
                    │ 1
                    │
                    │ N
              ┌─────┴──────┐
              │    rows    │
              ├────────────┤
              │ id     (PK)│
              │ sheet_id   │──── FK → sheets.id
              │ data       │ (JSON: {"col1":"val","col2":"val"})
              └────────────┘

┌──────────────────┐
│     records      │
├──────────────────┤
│ id          (PK) │
│ created_at      │
│ updated_at      │
│ deleted_at      │
│ nombre          │
│ data            │ (JSON flexible)
└──────────────────┘
```

---

## Requisitos Previos

### macOS

```bash
# 1. Instalar Go
brew install go

# 2. Instalar MySQL
brew install mysql

# 3. Iniciar MySQL
brew services start mysql

# 4. (Opcional) Instalar Air para live-reload
go install github.com/air-verse/air@latest
```

### Ubuntu/Debian

```bash
# 1. Instalar Go
sudo apt update
sudo apt install -y golang-go

# 2. Instalar MySQL
sudo apt install -y mysql-server
sudo systemctl start mysql

# 3. (Opcional) Instalar Air
go install github.com/air-verse/air@latest
```

### Windows

```bash
# 1. Instalar Go desde https://go.dev/dl/
# 2. Instalar MySQL desde https://dev.mysql.com/downloads/installer/
# 3. (Opcional) Instalar Air
go install github.com/air-verse/air@latest
```

---

## Instalación Paso a Paso

### Paso 1: Clonar el repositorio

### Paso 2: Configurar variables de entorno

```bash
# Copiar el archivo de ejemplo
cp .env.example .env
```

Editar el archivo `.env` con tus datos:

```env
# Base de datos
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASS=tu_contraseña    # Dejar vacío si no tienes contraseña
DB_NAME=test_web_excel1

# Servidor
PORT=3000
CORS_ORIGIN=*
```

### Paso 3: Crear la base de datos

**Opción A — Desde la línea de comandos:**

```bash
mysql -u root -e "CREATE DATABASE IF NOT EXISTS test_web_excel1;"
```

**Opción B — Desde Sequel Ace / MySQL Workbench:**

Ejecuta el siguiente SQL:

```sql
CREATE DATABASE IF NOT EXISTS test_web_excel1;
```

> Las tablas se crean automáticamente al iniciar la aplicación (AutoMigrate).

### Paso 4: Instalar dependencias Go

```bash
go mod download
```

### Paso 5: Ejecutar el servidor

**Desarrollo (con live-reload):**

```bash
air
```

**Producción:**

```bash
go run main.go
```

El servidor iniciará en `http://localhost:3000`.

### Paso 6: Verificar que funciona

```bash
curl http://localhost:3000/
# Respuesta: "funcionando correctamente api rest excel a mysql"
```

---

## Docker

Si prefieres no instalar MySQL localmente, puedes usar Docker:

```bash
# Levantar MySQL + App en contenedores
docker compose up -d

# Ver logs
docker compose logs -f

# Detener
docker compose down
```

Esto levanta:
- **MySQL** en puerto `3306`
- **App Go** en puerto `3000`

---

## Variables de Entorno

| Variable         | Default           | Descripción                                    |
| ---------------- | ----------------- | ---------------------------------------------- |
| `DB_HOST`        | `localhost`       | Host de MySQL                                  |
| `DB_PORT`        | `3306`            | Puerto de MySQL                                |
| `DB_USER`        | `root`            | Usuario de MySQL                               |
| `DB_PASS`        | *(vacío)*         | Contraseña de MySQL                            |
| `DB_NAME`        | `test_web_excel1` | Nombre de la base de datos                     |
| `MYSQL_HOST`     | → `DB_HOST`       | Variable automática de Railway MySQL           |
| `MYSQL_PORT`     | → `DB_PORT`       | Variable automática de Railway MySQL           |
| `MYSQL_USER`     | → `DB_USER`       | Variable automática de Railway MySQL           |
| `MYSQL_PASSWORD` | → `DB_PASS`       | Variable automática de Railway MySQL           |
| `MYSQL_DATABASE` | → `DB_NAME`       | Variable automática de Railway MySQL           |
| `PORT`           | `3000`            | Puerto del servidor HTTP                       |
| `CORS_ORIGIN`    | `*`               | Origen permitido para CORS                     |

> **Nota:** El código detecta primero las variables `MYSQL_*` (proporcionadas por Railway) y usa `DB_*` como fallback.

---

## API Endpoints

### Health Check

```
GET /
```

```bash
curl http://localhost:3000/
```

### Registros (CRUD)

| Método  | Ruta            | Descripción                  |
| ------- | --------------- | ---------------------------- |
| `GET`   | `/names`        | Listar todos los registros   |
| `GET`   | `/names/{id}`   | Obtener registro por ID      |
| `POST`  | `/names`        | Crear un registro            |
| `POST`  | `/names/all`    | Crear múltiples registros    |
| `DELETE` | `/names/{id}`  | Eliminar registro por ID     |
| `DELETE` | `/names`       | Eliminar todos los registros |

**Ejemplo — Crear registro:**

```bash
curl -X POST http://localhost:3000/names \
  -H "Content-Type: application/json" \
  -d '{
    "nombre": "Juan Pérez",
    "data": {"math": "90", "english": "85"}
  }'
```

**Ejemplo — Listar registros:**

```bash
curl http://localhost:3000/names
```

### Excel

| Método | Ruta             | Descripción                          |
| ------ | ---------------- | ------------------------------------ |
| `POST` | `/excel/upload`  | Subir estructura de archivo Excel    |
| `GET`  | `/excel/uploads` | Listar todas las subidas realizadas  |

**Ejemplo — Subir Excel:**

```bash
curl -X POST http://localhost:3000/excel/upload \
  -H "Content-Type: application/json" \
  -d '{
    "file": "notas_2024.xlsx",
    "sheets": [
      {
        "sheet": "1er Trimestre",
        "hasHeaders": true,
        "headers": ["nombre", "math", "english"],
        "rows": [
          {"nombre": "Juan", "math": "90", "english": "85"},
          {"nombre": "María", "math": "95", "english": "88"}
        ]
      }
    ]
  }'
```

**Ejemplo — Listar subidas:**

```bash
curl http://localhost:3000/excel/uploads
```

---

## Conectar el Frontend

### Web (JavaScript / HTML)

```html
<!DOCTYPE html>
<html>
<head><title>Excel to MySQL</title></head>
<body>
  <h1>Registros</h1>
  <div id="data"></div>

  <script>
    const API = 'http://localhost:3000';

    // Listar registros
    async function loadRecords() {
      const res = await fetch(`${API}/names`);
      const data = await res.json();
      document.getElementById('data').innerHTML =
        JSON.stringify(data, null, 2);
    }

    // Crear registro
    async function createRecord() {
      await fetch(`${API}/names`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          nombre: 'Juan Pérez',
          data: { math: '90', english: '85' }
        })
      });
      loadRecords();
    }

    loadRecords();
  </script>
</body>
</html>
```

### React

```jsx
const API = 'http://localhost:3000';

// Listar registros
const getRecords = async () => {
  const res = await fetch(`${API}/names`);
  return res.json();
};

// Crear registro
const createRecord = async (data) => {
  const res = await fetch(`${API}/names`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  });
  return res.json();
};

// Subir Excel
const uploadExcel = async (excelData) => {
  const res = await fetch(`${API}/excel/upload`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(excelData),
  });
  return res.json();
};
```

### Android (Kotlin + Retrofit)

```kotlin
// 1. Definir la interfaz de la API
interface ExcelApi {
    @GET("names")
    suspend fun getRecords(): List<Record>

    @POST("names")
    suspend fun createRecord(@Body record: Record): Record

    @POST("excel/upload")
    suspend fun uploadExcel(@Body data: UploadRequest): UploadResponse
}

// 2. Configurar Retrofit
val retrofit = Retrofit.Builder()
    .baseUrl("http://10.0.2.2:3000/") // Emulador Android
    .addConverterFactory(GsonConverterFactory.create())
    .build()

val api = retrofit.create(ExcelApi::class.java)

// 3. Usar en una corrutina
lifecycleScope.launch {
    val records = api.getRecords()
    // Actualizar UI con los registros
}
```

> **Nota:** Para dispositivos físicos, reemplaza `10.0.2.2` con la IP local de tu máquina (ej: `192.168.1.100`).

### Flutter (Dart)

```dart
import 'dart:convert';
import 'package:http/http.dart' as http;

class ApiService {
  final String baseUrl = 'http://localhost:3000';

  // Listar registros
  Future<List<dynamic>> getRecords() async {
    final response = await http.get(Uri.parse('$baseUrl/names'));
    if (response.statusCode == 200) {
      return json.decode(response.body);
    }
    throw Exception('Error al obtener registros');
  }

  // Crear registro
  Future<Map<String, dynamic>> createRecord(Map<String, dynamic> data) async {
    final response = await http.post(
      Uri.parse('$baseUrl/names'),
      headers: {'Content-Type': 'application/json'},
      body: json.encode(data),
    );
    if (response.statusCode == 200) {
      return json.decode(response.body);
    }
    throw Exception('Error al crear registro');
  }

  // Subir Excel
  Future<Map<String, dynamic>> uploadExcel(Map<String, dynamic> excelData) async {
    final response = await http.post(
      Uri.parse('$baseUrl/excel/upload'),
      headers: {'Content-Type': 'application/json'},
      body: json.encode(excelData),
    );
    if (response.statusCode == 200) {
      return json.decode(response.body);
    }
    throw Exception('Error al subir Excel');
  }
}
```

---

## Tests

```bash
# Ejecutar todos los tests
go test ./...

# Con verbose
go test -v ./...

# Con cobertura
go test -cover ./...
```

Los tests cubren:
- Validación de datos (`isJSONContent`)
- Handlers HTTP (códigos de respuesta, errores de validación)
- Modelos y sus relaciones
- Deserialización de JSON

---

## Troubleshooting

| Problema | Solución |
|----------|----------|
| `Error conectando a la DB` | Verifica que MySQL esté corriendo: `brew services list` |
| `port already in use` | Cambia el puerto en `.env` o mata el proceso: `kill $(lsof -ti:3000)` |
| `Access denied for user 'root'` | Verifica la contraseña en `.env` |
| `Table doesn't exist` | Las tablas se crean automáticamente; reinicia el servidor |
| CORS error en frontend | Verifica `CORS_ORIGIN=*` en `.env` |

---

## Licencia

MIT
