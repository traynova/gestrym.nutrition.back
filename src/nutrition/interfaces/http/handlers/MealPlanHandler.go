package handlers

import (
	"gestrym-nutrition/src/nutrition/application/usecases"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MealPlanHandler struct {
	CreateUC       *usecases.CreateMealPlanUseCase
	GetUC          *usecases.GetMealPlanUseCase
	GetUserPlansUC *usecases.GetUserMealPlansUseCase
	AddDayUC       *usecases.AddMealDayUseCase
	AddItemUC      *usecases.AddMealItemUseCase
}

func NewMealPlanHandler(
	createUC *usecases.CreateMealPlanUseCase,
	getUC *usecases.GetMealPlanUseCase,
	getUserPlansUC *usecases.GetUserMealPlansUseCase,
	addDayUC *usecases.AddMealDayUseCase,
	addItemUC *usecases.AddMealItemUseCase,
) *MealPlanHandler {
	return &MealPlanHandler{
		CreateUC:       createUC,
		GetUC:          getUC,
		GetUserPlansUC: getUserPlansUC,
		AddDayUC:       addDayUC,
		AddItemUC:      addItemUC,
	}
}

// CreateMealPlan godoc
// @Summary      Create a meal plan
// @Description  Creates a new meal plan for a user or as a template.
// @Tags         MealPlans
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body  usecases.CreateMealPlanInput  true  "Meal Plan Data"
// @Success      201   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Router       /private/meal-plans [post]
func (h *MealPlanHandler) CreateMealPlan(c *gin.Context) {
	var input usecases.CreateMealPlanInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Inject requester as creator if not set
	requesterID, _ := c.Get("user_id")
	if input.CreatedBy == 0 {
		input.CreatedBy = requesterID.(uint)
	}
	if input.UserID == 0 {
		input.UserID = requesterID.(uint)
	}

	plan, err := h.CreateUC.Execute(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": plan})
}

// GetMealPlan godoc
// @Summary      Get meal plan by ID
// @Description  Returns a meal plan with its days and items.
// @Tags         MealPlans
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Meal Plan ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /private/meal-plans/{id} [get]
func (h *MealPlanHandler) GetMealPlan(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	requesterID := c.MustGet("user_id").(uint)
	requesterRoleID := c.MustGet("role_id").(uint)

	plan, err := h.GetUC.Execute(uint(id), requesterID, requesterRoleID)
	if err != nil {
		if err.Error() == "access denied: plan belongs to another user" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "meal plan not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": plan})
}

// GetUserMealPlans godoc
// @Summary      Get meal plans for a user
// @Description  Returns all meal plans assigned to a user.
// @Tags         MealPlans
// @Produce      json
// @Security     BearerAuth
// @Param        userId  path  int  true  "User ID"
// @Success      200     {object}  map[string]interface{}
// @Router       /private/meal-plans/user/{userId} [get]
func (h *MealPlanHandler) GetUserMealPlans(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid userId"})
		return
	}

	plans, err := h.GetUserPlansUC.Execute(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": plans, "total": len(plans)})
}

// AddMealDay godoc
// @Summary      Add a day to a meal plan
// @Description  Adds a new meal day to an existing meal plan.
// @Tags         MealPlans
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      int                         true  "Meal Plan ID"
// @Param        body  body      usecases.AddMealDayInput    true  "Day data"
// @Success      201   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Router       /private/meal-plans/{id}/days [post]
func (h *MealPlanHandler) AddMealDay(c *gin.Context) {
	planID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid plan id"})
		return
	}

	var input usecases.AddMealDayInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	input.MealPlanID = uint(planID)

	day, err := h.AddDayUC.Execute(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": day})
}

// AddMealItem godoc
// @Summary      Add food item to a meal plan day
// @Description  Assigns a food item with quantity and meal type to a plan's day.
// @Tags         MealPlans
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      int                          true  "Meal Plan ID (used for context)"
// @Param        body  body      usecases.AddMealItemInput    true  "Item data"
// @Success      201   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Router       /private/meal-plans/{id}/items [post]
func (h *MealPlanHandler) AddMealItem(c *gin.Context) {
	var input usecases.AddMealItemInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := h.AddItemUC.Execute(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": item})
}
