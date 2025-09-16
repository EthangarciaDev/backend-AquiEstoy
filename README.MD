# Backend AquiEstoy

Backend desarrollado en Go con Gin Framework, PostgreSQL (AWS RDS) y autenticación JWT.

## 🚀 Características

- **Framework**: Gin (Go)
- **Base de datos**: PostgreSQL (AWS RDS)
- **Autenticación**: JWT (JSON Web Tokens)
- **ORM**: GORM
- **Encriptación**: bcrypt para passwords
- **Variables de entorno**: godotenv

## 📁 Estructura del Proyecto

```
.
├── main.go                     # Punto de entrada de la aplicación
├── go.mod                      # Dependencias de Go
├── go.sum                      # Checksums de dependencias
├── .env.example               # Variables de entorno de ejemplo
├── README.md                  # Este archivo
└── internal/
    ├── config/
    │   └── database.go        # Configuración de base de datos
    ├── handlers/
    │   └── auth_handler.go    # Handlers de autenticación
    ├── models/
    │   └── user.go           # Modelos de datos
    └── services/
        ├── auth_service.go   # Servicio de autenticación
        └── user_service.go   # Servicio de usuarios
```

## 🛠️ Configuración

### 1. Clonar el repositorio

```bash
git clone https://github.com/EthangarciaDev/backend-AquiEstoy.git
cd backend-AquiEstoy
```

### 2. Instalar dependencias

```bash
go mod tidy
```

### 3. Configurar variables de entorno

Copia el archivo `.env.example` a `.env` y configura tus variables:

```bash
cp .env.example .env
```

Edita el archivo `.env` con tu configuración de AWS RDS:

```env
DB_HOST=tu-rds-endpoint.region.rds.amazonaws.com
DB_PORT=5432
DB_USER=tu_usuario
DB_PASSWORD=tu_password
DB_NAME=tu_base_de_datos
DB_SSLMODE=require
JWT_SECRET=tu-clave-secreta-jwt-super-segura
```

### 4. Ejecutar la aplicación

```bash
go run main.go
```

El servidor se iniciará en `http://localhost:8080`

## 📚 API Endpoints

### Rutas Públicas

#### Health Check

- **GET** `/`
  - Descripción: Verifica que el servidor esté funcionando
  - Respuesta: `{"message": "Servidor AquiEstoy Backend funcionando 🚀", "status": "ok"}`

#### Registro de Usuario

- **POST** `/auth/register`
  - Descripción: Registra un nuevo usuario
  - Body:
    ```json
    {
      "email": "usuario@ejemplo.com",
      "password": "password123",
      "name": "Nombre Usuario"
    }
    ```

#### Inicio de Sesión

- **POST** `/auth/login`
  - Descripción: Autentica a un usuario existente
  - Body:
    ```json
    {
      "email": "usuario@ejemplo.com",
      "password": "password123"
    }
    ```

### Rutas Protegidas

Estas rutas requieren el header `Authorization: Bearer <token>`

#### Perfil de Usuario

- **GET** `/api/profile`
  - Descripción: Obtiene el perfil del usuario autenticado
  - Headers: `Authorization: Bearer <token>`

## 🔐 Autenticación

El sistema utiliza JWT (JSON Web Tokens). Incluye el token en el header:

```
Authorization: Bearer <tu_jwt_token>
```

## 🗄️ Base de Datos

PostgreSQL hospedado en AWS RDS. Las migraciones se ejecutan automáticamente.
git clone https://github.com/EthangarciaDev/backend-AquiEstoy.git
cd backend-AquiEstoy
go mod tidy # Descarga las dependencias
go run . # Arranca el servidor en http://localhost:8080

```

```
