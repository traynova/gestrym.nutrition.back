# 🧠 IA_MEMORY — Gestrym Nutrition Service

> **Última actualización:** 2026-04-27
> **Ingeniería:** Senior Golang Backend — Microservicios, Nutrición, Arquitectura Hexagonal

---

## 🎯 CONTEXTO DEL PROYECTO

Plataforma fitness **Gestrym** con los siguientes microservicios existentes:

| Servicio              | Descripción                             |
| --------------------- | --------------------------------------- |
| `auth-service`        | Usuarios, roles                         |
| `training-service`    | Ejercicios, workouts, planes de entreno |
| `progress-service`    | Métricas, fotos, notas                  |
| `notification-service`| Notificaciones (email, SMS, push, in-app)|
| `storage-service`     | MinIO para archivos                     |
| `nutrition-service`   | **← ESTE SERVICIO (en construcción)**   |

**Módulo Go:** `gestrym-nutrition`
**Base de datos:** PostgreSQL con GORM
**Framework HTTP:** Gin
**Config:** Viper

---

## ⚠️ REGLAS CRÍTICAS (NUNCA VIOLAR)

1. ❌ NO modificar la estructura existente del proyecto
2. ✅ TODOS los modelos DEBEN estar en `common/models`
3. ❌ NO duplicar modelos
4. ✅ Usar GORM
5. ✅ Seguir arquitectura hexagonal estrictamente:
   - `domain` → interfaces/ports
   - `application` → use cases
   - `infrastructure` → repositorios GORM, adapters
   - `interfaces/http` → handlers Gin
6. ✅ Inyección de dependencias
7. ❌ NO llamar APIs externas para alimentos (ya importados)

---

## 🏗️ ESTRUCTURA ACTUAL DEL PROYECTO

```
Gestrym.Nutrition.Back/
├── main.go
├── go.mod                    # module: gestrym-nutrition
├── src/
│   ├── app.go               # Bootstrap: setupEnvironment + initServer
│   ├── common/
│   │   ├── config/
│   │   │   ├── Database.go
│   │   │   ├── Enviroment.go
│   │   │   └── Migrations.go  ← AutoMigrate() aquí
│   │   ├── middleware/
│   │   │   ├── JWTModdleware.go   ← SetupJWTMiddleware()
│   │   │   ├── RoleMiddleware.go  ← RequireRoles(roleIDs...)
│   │   │   ├── ApiKeyMiddleware.go
│   │   │   ├── BasicAuthMiddleware.go
│   │   │   └── GinLoggerMiddleware.go
│   │   ├── models/
│   │   │   ├── Food.go           ← EXISTENTE
│   │   │   └── FoodCategory.go   ← EXISTENTE
│   │   ├── routes/
│   │   │   └── ServerRoutesDefinition.go  ← Registro de rutas y DI
│   │   ├── shared/
│   │   │   └── PaginateResponse.go
│   │   └── utils/
│   └── nutrition/
│       ├── domain/
│       │   └── interfaces/
│       │       ├── FoodRepository.go   ← EXISTENTE
│       │       ├── ImageProvider.go
│       │       ├── StorageService.go
│       │       ├── FileStorageAdapter.go
│       │       └── USDAAdapter.go
│       ├── application/
│       │   ├── usecases/
│       │   │   ├── GetFoodByIDUseCase.go
│       │   │   ├── SearchFoodsUseCase.go
│       │   │   └── ImportFoodsWithImagesUseCase.go
│       │   └── utils/
│       ├── infrastructure/
│       │   ├── repositories/
│       │   │   └── FoodRepositoryImpl.go ← EXISTENTE
│       │   └── adapters/
│       │       ├── FileStorageAdapterlmol.go
│       │       ├── PexelsAdapterImpl.go
│       │       ├── StorageServiceAdapterImpl.go
│       │       └── USDAAdapterImpl.go
│       └── interfaces/
│           └── http/
│               └── handlers/
│                   └── FoodHandler.go  ← EXISTENTE
```

---

## 🧱 MODELOS IMPLEMENTADOS (common/models)

### ✅ Existentes
- **Food** — ID, Name, CategoryID, Category, Calories, Protein, Carbs, Fats, ImageURL, CollectionID, CreatedAt, UpdatedAt
- **FoodCategory** — ID, Name

### 🆕 A Implementar
- **MealPlan** — ID, UserID, Name, DurationDays, CreatedBy, IsTemplate, GoalCalories, GoalProtein, GoalCarbs, GoalFats, CreatedAt
- **MealDay** — ID, MealPlanID, DayNumber
- **MealItem** — ID, MealDayID, FoodID, Quantity (gramos), MealType (breakfast/lunch/dinner/snack)
- **NutritionLog** — ID, UserID, Date, FoodID, Quantity, Calories, Protein, Carbs, Fats (valores pre-calculados)
- **UserCalorieGoal** — ID, UserID, WeightKg, HeightCm, AgeYears, ActivityLevel, FitnessGoal, TargetCalories/Macros, LastAdjustedAt, AdjustedByAI

---

## ⚙️ ROLES DEL SISTEMA

| Constante       | ID |
| --------------- | -- |
| `RoleAdmin`     | 1  |
| `RoleGym`       | 2  |
| `RoleCoach`     | 3  |
| `RoleCliente`   | 4  |

**JWT Claims disponibles:** `user_id` (uint), `role_id` (uint), `access_level_id` (uint)

---

## 🌐 ENDPOINTS IMPLEMENTADOS

### Públicos (`/gestrym-nutrition/public`)
| Método | Ruta              | Handler               |
| ------ | ----------------- | --------------------- |
| GET    | `/foods`          | SearchFoods           |
| GET    | `/foods/:id`      | GetFoodByID           |
| POST   | `/foods/import`   | ImportFoods (USDA)    |

### Privados (`/gestrym-nutrition/private`) — requieren JWT
| Método | Ruta                              | Descripción                               |
| ------ | --------------------------------- | ----------------------------------------- |
| POST   | `/meal-plans`                     | Crear plan de comida                      |
| GET    | `/meal-plans/:id`                 | Obtener detalle de plan                   |
| GET    | `/meal-plans/user/:userId`        | Listar planes por usuario                 |
| POST   | `/meal-plans/:id/days`            | Agregar día al plan                       |
| POST   | `/meal-plans/:id/items`           | Agregar alimento al día                   |
| POST   | `/logs`                           | Registrar consumo de alimento             |
| GET    | `/logs`                           | Resumen nutricional diario                |
| GET    | `/logs/history`                   | Historial de logs paginado                |
| POST   | `/goals/calories`                 | Configurar metas calóricas (TDEE)         |
| GET    | `/goals/calories`                 | Obtener metas actuales                    |
| POST   | `/goals/calories/adjust`          | **IA:** Ajustar metas con progress-service|

---

## 🧠 LÓGICA DE IA Y ADAPTACIÓN
El sistema se integra con el **progress-service** para:
1. Obtener peso y altura real del usuario.
2. Calcular el delta de peso semanal.
3. **Adaptar Calorías:** 
   - Si no baja de peso en déficit -> -100 kcal.
   - Si no sube en superávit -> +100 kcal.
   - Si baja demasiado rápido -> +200 kcal (protección muscular).

---

## 📄 DOCUMENTACIÓN PARA FRONTEND
- `docs/FRONTEND_GUIDE.md`: Guía técnica de integración y modelos.
- `docs/AI_FRONTEND_IMPLEMENTATION_PROMPT.md`: Prompt maestro para generar la UI con IA.

---

## 📋 PROMPT MAESTRO (CONTEXTO IA)

```
You are a senior Golang backend engineer specialized in microservices,
nutrition systems, and hexagonal architecture.

I already have a fitness platform with:
* auth-service (users, roles)
* training-service (exercises, workouts, training plans)
* progress-service (metrics, photos, notes)
* notification-service
* storage-service (MinIO)

Now I am building the nutrition-service.
```

### 🎯 OBJETIVO PRINCIPAL
Implementar un sistema COMPLETO de nutrición:
1. Planes de comida (semanal/mensual)
2. Estructura diaria de comidas
3. Asignación de alimentos (usando foods existentes)
4. Tracking de nutrición (qué come el usuario realmente)
5. Base para adaptación con IA (futuro)

### 🧠 LÓGICA DE TRACKING
Al registrar un alimento:
- Buscar datos nutricionales del food en BD
- Calcular: calories, protein, carbs, fats según la cantidad en gramos
- Fórmula: `nutriente = (food.Nutriente / 100) * quantity`
- Guardar valores calculados (NO recalcular cada vez)

### 📦 RESPUESTA DE NUTRICIÓN DIARIA (Frontend Friendly)
```json
{
  "date": "2024-01-15",
  "totals": {
    "calories": 2200,
    "protein": 150,
    "carbs": 250,
    "fats": 70
  },
  "goals": {
    "calories": 2500,
    "protein": 180,
    "carbs": 280,
    "fats": 80
  },
  "progress": {
    "calories_pct": 88,
    "protein_pct": 83,
    "carbs_pct": 89,
    "fats_pct": 87
  },
  "foods": [...]
}
```

---

## 🔐 AUTORIZACIÓN

| Operación                         | Roles permitidos        |
| --------------------------------- | ----------------------- |
| Ver/crear sus propios datos       | RoleCliente, RoleCoach  |
| Ver planes de usuarios asignados  | RoleCoach, RoleGym      |
| Crear planes template             | RoleCoach, RoleAdmin    |

---

## 🚀 PREPARACIÓN FUTURA (IA INTEGRATION)

El sistema debe estar diseñado para:
- [ ] Planes de comida generados por IA
- [ ] Nutrición adaptativa (basada en progress-service)
- [ ] Objetivos calóricos personalizados
- [ ] Metas de macros
- [ ] Helper: `CalculateNutritionTotals(logs []NutritionLog)` → totals para dashboards

---

## 💡 BONUS FEATURES (IMPLEMENTAR)

- [x] Daily calorie target (campo en MealPlan: GoalCalories)
- [x] Macro target (GoalProtein, GoalCarbs, GoalFats)
- [ ] Progress vs target (porcentaje) en respuesta diaria
- [ ] Paginación para logs
- [ ] Validación de meal types (enum: breakfast, lunch, dinner, snack)

---

## 📝 NOTAS TÉCNICAS

### Patrón de Repositorio
```go
// Interface en: nutrition/domain/interfaces/
type XRepository interface { ... }

// Implementación en: nutrition/infrastructure/repositories/
type XRepositoryImpl struct { DB *gorm.DB }
func NewXRepositoryImpl(db *gorm.DB) interfaces.XRepository { ... }
```

### Patrón de Use Case
```go
// En: nutrition/application/usecases/
type XUseCase struct { Repo interfaces.XRepository }
func NewXUseCase(repo interfaces.XRepository) *XUseCase { ... }
func (uc *XUseCase) Execute(...) (...) { ... }
```

### Patrón de Handler
```go
// En: nutrition/interfaces/http/handlers/
type XHandler struct { UC *usecases.XUseCase }
func NewXHandler(uc *usecases.XUseCase) *XHandler { ... }
func (h *XHandler) Method(c *gin.Context) { ... }
```

### Registro de Rutas
- Todo se registra en: `src/common/routes/ServerRoutesDefinition.go`
- DI (Dependency Injection) en `addRoutes()`
- AutoMigrate en: `src/common/config/Migrations.go`

---

## 📅 HISTORIAL DE CAMBIOS

| Fecha      | Cambio                                                    |
| ---------- | --------------------------------------------------------- |
| 2026-04-27 | Creación de IA_MEMORY.md con estado inicial del proyecto  |
| 2026-04-27 | Prompt maestro definido para nutrition-service completo   |
