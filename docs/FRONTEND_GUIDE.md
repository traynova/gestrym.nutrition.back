# 🥗 Guía de Integración Frontend: Nutrition Service

Esta guía detalla cómo interactuar con el microservicio de nutrición de **Gestrym**. El servicio sigue una arquitectura hexagonal y expone endpoints RESTful para la gestión de alimentos, planes de comida, seguimiento diario y objetivos adaptativos.

---

## 🔑 Autenticación y Configuración

- **Base URL:** `/gestrym-nutrition`
- **Headers Requeridos:**
  - `Authorization: Bearer <JWT_TOKEN>` (Para rutas `/private`)
  - `Content-Type: application/json`

---

## 📐 Modelos de Datos (TypeScript)

### Alimentos e Inventario
```typescript
interface Food {
  id: number;
  name: string;
  categoryId: number;
  category?: { id: number; name: string };
  calories: number; // Por 100g
  protein: number;  // Por 100g
  carbs: number;    // Por 100g
  fats: number;     // Por 100g
  imageUrl: string;
}
```

### Seguimiento (Logs)
```typescript
interface NutritionLog {
  id: number;
  userId: number;
  date: string; // ISO Date
  foodId: number;
  food?: Food;
  quantity: number; // Gramos
  calories: number; // Pre-calculado
  protein: number;
  carbs: number;
  fats: number;
  mealType: 'breakfast' | 'lunch' | 'dinner' | 'snack';
  notes?: string;
}
```

### Objetivos Calóricos (IA-Adaptive)
```typescript
interface UserCalorieGoal {
  userId: number;
  weightKg: number;
  heightCm: number;
  ageYears: number;
  isMale: boolean;
  activityLevel: 'sedentary' | 'light' | 'moderate' | 'active' | 'very_active';
  fitnessGoal: 'lose_weight' | 'maintain' | 'gain_mass';
  targetCalories: number;
  targetProtein: number;
  targetCarbs: number;
  targetFats: number;
  lastAdjustedAt?: string;
  adjustedByAI: boolean;
  adjustmentNote?: string;
}
```

---

## 🌐 Endpoints Principales

### 🍎 1. Catálogo de Alimentos (Público)

| Método | Ruta | Descripción |
| :--- | :--- | :--- |
| `GET` | `/public/foods?search=pollo&page=1&limit=10` | Buscar alimentos en la base de datos. |
| `GET` | `/public/foods/:id` | Obtener detalles de un alimento específico. |

---

### 📅 2. Planes de Comida (Privado)

| Método | Ruta | Descripción |
| :--- | :--- | :--- |
| `POST` | `/private/meal-plans` | Crear un nuevo plan (semanal/mensual). |
| `GET` | `/private/meal-plans/:id` | Obtener plan con sus días y alimentos asignados. |
| `GET` | `/private/meal-plans/user/:userId` | Listar todos los planes de un usuario. |
| `POST` | `/private/meal-plans/:id/days` | Agregar un día a un plan (ej: Día 1). |
| `POST` | `/private/meal-plans/:id/items` | Asignar un alimento a un día de comida. |

---

### 🍽️ 3. Seguimiento Diario (Log)

| Método | Ruta | Descripción |
| :--- | :--- | :--- |
| `POST` | `/private/logs` | Registrar un alimento consumido (calcula macros auto). |
| `GET` | `/private/logs?date=YYYY-MM-DD` | Resumen diario: Totales vs Objetivos + Progreso %. |
| `GET` | `/private/logs/history?start=...&end=...` | Historial paginado entre fechas. |

**Respuesta Resumen Diario (`GET /private/logs`):**
```json
{
  "data": {
    "date": "2024-04-27",
    "totals": { "calories": 1850, "protein": 120, "carbs": 180, "fats": 45 },
    "goals": { "calories": 2200, "protein": 150, "carbs": 250, "fats": 70 },
    "progress": {
      "caloriesPct": 84.1,
      "proteinPct": 80.0,
      "carbsPct": 72.0,
      "fatsPct": 64.3
    },
    "foods": [...] 
  }
}
```

---

### 🤖 4. Objetivos y Ajuste IA (Privado)

| Método | Ruta | Descripción |
| :--- | :--- | :--- |
| `POST` | `/private/goals/calories` | Configurar perfil (TDEE) y calcular objetivos iniciales. |
| `GET` | `/private/goals/calories` | Obtener objetivos actuales y estado de ajuste IA. |
| `POST` | `/private/goals/calories/adjust` | **IA:** Cruza datos con `progress-service` y ajusta calorías. |

---

## 💡 Recomendaciones para el Frontend

1. **Dashboard de Macros:** Usa el endpoint `GET /private/logs` para pintar gráficas de progreso con el `progressPct`.
2. **Botón IA:** Implementa un botón de "Optimizar con IA" que llame a `/adjust`. Esta función analiza si el usuario está cumpliendo sus metas de peso en el `progress-service` y recalibra su dieta automáticamente.
3. **Paginación:** Los logs de historia usan paginación estándar (`page`, `pageSize`).
