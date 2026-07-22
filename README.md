# 🥗 Gestrym Nutrition Service Backend

> **Servicio de Gestión Nutricional, Planes de Alimentación y Tracking Adaptativo con IA para la plataforma fitness Gestrym.**

[![Go Version](https://img.shields.io/badge/Go-1.22%2B-00ADD8?style=for-the-badge&logo=go)](https://golang.org/)
[![Framework](https://img.shields.io/badge/Framework-Gin-008080?style=for-the-badge&logo=gin)](https://gin-gonic.com/)
[![ORM](https://img.shields.io/badge/ORM-GORM-blue?style=for-the-badge)](https://gorm.io/)
[![Database](https://img.shields.io/badge/Database-PostgreSQL-336791?style=for-the-badge&logo=postgresql)](https://www.postgresql.org/)
[![Architecture](https://img.shields.io/badge/Architecture-Hexagonal-orange?style=for-the-badge)](https://en.wikipedia.org/wiki/Hexagonal_architecture_(software))

---

## 📋 Tabla de Contenidos
- [🎯 Contexto del Proyecto](#-contexto-del-proyecto)
- [🛠️ Stack Tecnológico](#️-stack-tecnológico)
- [🏗️ Arquitectura Hexagonal](#️-arquitectura-hexagonal)
- [🧱 Modelos de Datos](#-modelos-de-datos)
- [🔐 Autenticación y Matriz de Permisos](#-autenticación-y-matriz-de-permisos)
- [🌐 Endpoints de la API](#-endpoints-de-la-api)
- [🧠 Lógica de Tracking e Integración con IA](#-lógica-de-tracking-e-integración-con-ia)
- [🚀 Configuración e Instalación](#-configuración-e-instalación)
- [📄 Documentación para Frontend & Swagger](#-documentación-para-frontend--swagger)

---

## 🎯 Contexto del Proyecto

**Gestrym Nutrition Service** (`gestrym-nutrition`) es el microservicio encargado del módulo de nutrición dentro de la plataforma fitness integral **Gestrym**. 

### 🌌 Ecosistema de Microservicios Gestrym

| Microservicio | Responsabilidad Principal | Estado |
| :--- | :--- | :---: |
| `auth-service` | Gestión de usuarios, autenticación y roles | 🟢 Activo |
| `training-service` | Ejercicios, rutinas y planes de entrenamiento | 🟢 Activo |
| `progress-service` | Seguimiento de peso, métricas corporales, fotos y notas | 🟢 Activo |
| `notification-service` | Envío de notificaciones (Email, Push, SMS, In-App) | 🟢 Activo |
| `storage-service` | Almacenamiento distribuido de multimedia (MinIO) | 🟢 Activo |
| **`nutrition-service`** | **Base de alimentos, planes nutricionales, registro diario e IA** | 🟡 **Este Repositorio** |

---

## 🛠️ Stack Tecnológico

- **Lenguaje:** Go (Golang) 1.22+
- **Framework Web HTTP:** [Gin Gonic](https://github.com/gin-gonic/gin)
- **Base de Datos:** PostgreSQL
- **ORM:** [GORM](https://gorm.io/)
- **Gestión de Configuración:** [Viper](https://github.com/spf13/viper)
- **Documentación API:** Swagger / OpenAPI 2.0 (Swag)
- **Validación de Datos:** Go-Playground Validator v10

---

## 🏗️ Arquitectura Hexagonal

El proyecto está diseñado bajo los principios de **Arquitectura Hexagonal (Puertos y Adaptadores)** para garantizar desacoplamiento, testabilidad y mantenibilidad.

```
gestrym.nutrition.back/
├── main.go                       # Punto de entrada principal
├── docs/                         # Guías de integración y especificación Swagger
└── src/
    ├── app.go                    # Inicialización de entorno, BD y servidor
    ├── common/                   # Transversal: Config, Middlewares, Modelos, Rutas
    │   ├── config/               # Base de datos, variables de entorno y migraciones
    │   ├── middleware/           # Autenticación JWT, API Key, BasicAuth y Logger
    │   ├── models/               # Entidades y Modelos GORM compartidos
    │   ├── routes/               # Registro global de rutas e Inyección de Dependencias
    │   ├── shared/               # DTOs comunes (PaginateResponse)
    │   └── utils/                # Utilidades generales y validadores
    └── nutrition/                # Módulo Core de Nutrición
        ├── domain/               # [Capa de Dominio] Puertos e Interfaces
        │   └── interfaces/       # Repositorios y adaptadores de servicios externos
        ├── application/          # [Capa de Aplicación] Casos de Uso
        │   └── usecases/         # Lógica de negocio (Búsqueda, Planes, Tracking)
        ├── infrastructure/       # [Capa de Infraestructura] Adaptadores
        │   ├── repositories/     # Implementación GORM de repositorios
        │   └── adapters/         # USDA, Pexels, StorageService
        └── interfaces/           # [Capa de Entrada] Handlers HTTP
            └── http/
                └── handlers/     # Controladores Gin HTTP
```

---

## 🧱 Modelos de Datos

Todos los modelos persistentes se ubican estrictamente en `src/common/models/` y son gestionados mediante auto-migraciones de GORM.

- **`Food`**: Representa los alimentos de la base de datos (con sus valores nutricionales por cada 100g: calorías, proteínas, carbohidratos, grasas e imágenes).
- **`FoodCategory`**: Categorización de alimentos (Proteínas, Carbohidratos, Frutas, Verduras, etc.).
- **`MealPlan`**: Plan nutricional asignado a un usuario (meta calórica, distribución de macros, creador y bandera de generación por IA).
- **`MealDay`**: Días que integran un plan nutricional (Día 1, Día 2, etc.).
- **`MealItem`**: Alimento específico asignado a un día con su tipo de comida (`breakfast`, `lunch`, `dinner`, `snack`) y cantidad en gramos.
- **`NutritionLog`**: Registro del consumo real del usuario en una fecha específica (guarda valores calóricos precalculados).
- **`UserCalorieGoal`**: Configuración de metas calóricas personalizadas (TDEE, peso, altura, edad, nivel de actividad y objetivo fitness).

---

## 🔐 Autenticación y Matriz de Permisos

El servicio valida solicitudes mediante **JWT (JSON Web Tokens)** generados por `auth-service`. Los claims leídos son `user_id`, `role_id` y `access_level_id`.

### Roles de Sistema

| Nombre de Rol | ID | Descripción |
| :--- | :---: | :--- |
| **`RoleAdmin`** | `1` | Administrador general de la plataforma |
| **`RoleGym`** | `2` | Gimnasio / Organización deportiva |
| **`RoleCoach`** | `3` | Entrenador o Nutricionista |
| **`RoleCliente`** | `4` | Usuario final / Atleta |

### Matriz de Operaciones

| Operación | Roles Permitidos |
| :--- | :--- |
| Ver / Crear sus propios logs y metas | `RoleCliente`, `RoleCoach` |
| Ver planes nutricionales de usuarios asignados | `RoleCoach`, `RoleGym` |
| Crear plantillas de planes de alimentación | `RoleCoach`, `RoleAdmin` |

---

## 🌐 Endpoints de la API

### 🔓 Rutas Públicas (`/gestrym-nutrition/public`)

| Método | Ruta | Descripción | Handler |
| :---: | :--- | :--- | :--- |
| `GET` | `/foods` | Buscar alimentos en la base de datos con paginación | `SearchFoods` |
| `GET` | `/foods/:id` | Obtener detalle nutricional de un alimento por ID | `GetFoodByID` |
| `POST` | `/foods/import` | Importar alimentos desde fuentes externas (USDA / Pexels) | `ImportFoods` |

### 🔒 Rutas Privadas (`/gestrym-nutrition/private`) — *Requieren JWT*

| Método | Ruta | Descripción |
| :---: | :--- | :--- |
| `POST` | `/meal-plans` | Crear un nuevo plan de alimentación |
| `GET` | `/meal-plans/:id` | Obtener el detalle de un plan específico |
| `GET` | `/meal-plans/user/:userId` | Listar todos los planes nutricionales de un usuario |
| `POST` | `/meal-plans/:id/days` | Agregar un nuevo día a un plan |
| `POST` | `/meal-plans/:id/items` | Agregar un alimento/ítem a un día del plan |
| `POST` | `/logs` | Registrar consumo real de un alimento por el usuario |
| `GET` | `/logs` | Obtener el resumen nutricional del día actual |
| `GET` | `/logs/history` | Consultar historial de registros nutricionales paginados |
| `POST` | `/goals/calories` | Configurar las metas calóricas y macros (TDEE) |
| `GET` | `/goals/calories` | Obtener las metas calóricas actuales del usuario |
| `POST` | `/goals/calories/adjust` | **IA:** Ajustar metas automáticamente según `progress-service` |
| `POST` | `/internal/meal-plans/ai` | **IA:** Crear/Actualizar planes mediante AI Service |

---

## 🧠 Lógica de Tracking e Integración con IA

### 📐 Lógica de Cálculo de Nutrientes
Al registrar el consumo de un alimento en `NutritionLog`, los valores nutricionales finales se calculan y persisten utilizando la fórmula:

$$\text{Nutriente Calculado} = \left( \frac{\text{Food.Nutriente}}{100} \right) \times \text{Cantidad (gramos)}$$

*Esto evita recalcular datos históricos si la información del alimento cambia en el futuro.*

### 🤖 Adaptación Nutricional Inteligente
El microservicio se conecta con el `progress-service` para realizar ajustes calóricos automáticos basados en la tendencia de peso semanal:

- **Sin pérdida en Déficit:** Ajuste de $-100 \text{ kcal}$.
- **Sin ganancia en Superávit:** Ajuste de $+100 \text{ kcal}$.
- **Pérdida excesivamente rápida:** Ajuste de $+200 \text{ kcal}$ para protección de masa muscular.

---

## 🚀 Configuración e Instalación

### Prerrequisitos
- **Go** 1.22 o superior.
- **PostgreSQL** 14+.

### Pasos de Ejecución

1. **Clonar el repositorio:**
   ```bash
   git clone https://github.com/tu-usuario/gestrym.nutrition.back.git
   cd gestrym.nutrition.back
   ```

2. **Instalar dependencias:**
   ```bash
   go mod download
   ```

3. **Configurar Variables de Entorno:**
   Crea o edita el archivo de configuración correspondiente en `./deployment/env_local.yaml`:

   ```yaml
   GIN_MODE: "debug"
   GORM_LOG_LEVEL: "info"
   POSTGRES_DB_HOST: "localhost"
   POSTGRES_DB_PORT: "5432"
   POSTGRES_DB_USER: "postgres"
   POSTGRES_DB_PASSWORD: "password"
   POSTGRES_DB_NAME: "gestrym_nutrition_db"
   POSTGRES_DB_SSLMODE: "disable"
   JWT_KEY: "your_jwt_secret_key"
   BASIC_AUTH_USERNAME: "admin"
   BASIC_AUTH_PASSWORD: "password"
   AUTH_API_KEY: "your_api_key"
   GESTRYM_TRAINNER_SERVER_ADDRESS: "http://localhost:8081"
   STORAGE_SERVICE_URL: "http://localhost:8085"
   STORAGE_SERVICE_API_KEY: "storage_api_key"
   RAPID_API_KEY: "rapid_api_key"
   RAPID_API_HOST: "usda-food-db.p.rapidapi.com"
   ```

4. **Ejecutar el microservicio:**
   ```bash
   # En entorno local
   go run main.go --env=local
   ```

---

## 📄 Documentación para Frontend & Swagger

El repositorio incluye guías completas y archivos Swagger listos para integración:

- 📘 **[Guía de Integración Frontend](docs/FRONTEND_GUIDE.md):** Especificación completa de payloads, respuestas JSON y flujos de integración para clientes Web / Móvil.
- 🤖 **[Prompt Maestro para UI con IA](docs/AI_FRONTEND_IMPLEMENTATION_PROMPT.md):** Contexto preconfigurado para generar interfaces de usuario de nutrición con modelos de lenguaje.
- 📜 **Swagger UI / Specification:** Los archivos `docs/swagger.json` y `docs/swagger.yaml` contienen la especificación OpenAPI detallada de los endpoints del servicio.

---

<p center="align">
  <b>Gestrym Platform</b> — Empowering Fitness Tech
</p>
