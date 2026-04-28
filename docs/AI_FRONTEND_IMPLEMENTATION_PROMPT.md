# 🤖 MASTER PROMPT: Implementación del Módulo de Nutrición (Gestrym)

**Instrucciones para la IA:** Actúa como un experto Senior Frontend Engineer especializado en React, TypeScript y Tailwind CSS. Tu objetivo es implementar el nuevo módulo de **Nutrición** para la plataforma Gestrym.

---

## contexto del Proyecto
- **Stack:** React (Vite), TypeScript, Tailwind CSS.
- **Estado:** TanStack Query (fetching), Zustand (global auth), Framer Motion (animaciones).
- **UI:** Recharts para gráficas, Lucide React para iconos.
- **Backend:** Microservicio de Nutrición (Hexagonal Go) con endpoints REST.

---

## 🎯 Objetivo de la Tarea
Implementar una interfaz premium y funcional para que los usuarios gestionen su nutrición. El módulo debe incluir:
1. **Dashboard Diario:** Resumen de macros (Donut charts) y progreso vs objetivos.
2. **Registro de Alimentos (Log):** Buscador con autocompletado y modal para registrar consumo en gramos.
3. **Planes de Comida:** Visualizador de planes semanales/mensuales asignados por el coach.
4. **Perfil de Objetivos IA:** Configuración de TDEE y botón de "Ajuste Adaptativo" que se integra con el servicio de progreso.

---

## 🛠️ Especificaciones Técnicas (API Contracts)

### Base URL: `/gestrym-nutrition`

#### 1. Nutrición Diaria (`GET /private/logs?date=YYYY-MM-DD`)
Devuelve totales del día, objetivos y porcentaje de progreso.
```typescript
interface DailyNutritionResponse {
  data: {
    date: string;
    totals: { calories: number; protein: number; carbs: number; fats: number };
    goals: { calories: number; protein: number; carbs: number; fats: number };
    progress: { caloriesPct: number; proteinPct: number; carbsPct: number; fatsPct: number };
    foods: Array<{ id: number; food: { name: string; imageUrl: string }; quantity: number; calories: number; mealType: string }>;
  }
}
```

#### 2. Registro (`POST /private/logs`)
Cuerpo: `{ "foodId": number, "quantity": number, "mealType": string, "date": string }`

#### 3. Objetivos Adaptativos (`POST /private/goals/calories/adjust`)
Endpoint clave de IA. No requiere cuerpo. Al llamarlo, el backend recalibra los macros del usuario basándose en su progreso real de peso. Devuelve un `adjustmentNote` que debe mostrarse como un toast o notificación de éxito.

---

## 🎨 Requerimientos de UI/UX (Premium Feel)
- **Visualización de Macros:** Usa `recharts` para mostrar un gráfico de dona central con las calorías y 3 barras de progreso menores para Proteína, Carbos y Grasas.
- **Experiencia de Búsqueda:** El buscador de alimentos debe ser instantáneo (useQuery con debounce). Muestra miniaturas de los alimentos.
- **Feed de Actividad:** Lista los alimentos del día agrupados por tipo (Desayuno, Almuerzo, etc.) con iconos descriptivos.
- **Feedback de IA:** Al presionar "Ajustar con IA", usa `framer-motion` para mostrar un estado de "analizando" antes de presentar el nuevo objetivo y la nota de ajuste del coach/sistema.

---

## 📂 Estructura de Archivos Sugerida
1. `src/services/nutritionService.ts`: Definición de llamadas Axios/TanStack Query.
2. `src/components/nutrition/NutritionDashboard.tsx`: Componente principal.
3. `src/components/nutrition/MacroCard.tsx`: Visualización de círculos de progreso.
4. `src/components/nutrition/FoodSearchModal.tsx`: Buscador e inserción de logs.
5. `src/components/nutrition/AIAdjustmentPanel.tsx`: Control de objetivos adaptativos.

---

## 🚀 Instrucción Final
Genera el código completo para estos componentes, asegurando que los tipos de TypeScript sean estrictos y que el diseño sea responsive, moderno (modo oscuro/claro compatible) y use las utilidades de Tailwind configuradas en el proyecto. Prioriza la legibilidad y el manejo de estados de carga/error.
