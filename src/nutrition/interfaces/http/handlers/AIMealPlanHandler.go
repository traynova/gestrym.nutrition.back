package handlers

import (
	"gestrym-nutrition/src/nutrition/application/usecases"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AIMealPlanHandler struct {
	CreateFromAIUC *usecases.CreateMealPlanFromAIUseCase
}

func NewAIMealPlanHandler(createFromAIUC *usecases.CreateMealPlanFromAIUseCase) *AIMealPlanHandler {
	return &AIMealPlanHandler{
		CreateFromAIUC: createFromAIUC,
	}
}

// CreateMealPlanFromAI godoc
// @Summary      Create a meal plan from AI
// @Description  Receives a structured meal plan from the AI service and stores it.
// @Tags         Internal
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        body  body  usecases.CreateMealPlanFromAIInput  true  "AI Meal Plan Data"
// @Success      201   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Router       /internal/meal-plans/ai [post]
func (h *AIMealPlanHandler) CreateMealPlanFromAI(c *gin.Context) {
	var input usecases.CreateMealPlanFromAIInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	plan, err := h.CreateFromAIUC.Execute(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "AI meal plan created successfully",
		"data":    plan,
	})
}
