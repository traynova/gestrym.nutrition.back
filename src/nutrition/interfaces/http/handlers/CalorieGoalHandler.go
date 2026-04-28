package handlers

import (
	"gestrym-nutrition/src/common/models"
	"gestrym-nutrition/src/nutrition/application/usecases"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CalorieGoalHandler struct {
	SetGoalUC    *usecases.SetUserCalorieGoalUseCase
	GetGoalUC    *usecases.GetUserCalorieGoalUseCase
	AdjustFromProgressUC *usecases.AdjustCaloriesFromProgressUseCase
}

func NewCalorieGoalHandler(
	setGoalUC *usecases.SetUserCalorieGoalUseCase,
	getGoalUC *usecases.GetUserCalorieGoalUseCase,
	adjustFromProgressUC *usecases.AdjustCaloriesFromProgressUseCase,
) *CalorieGoalHandler {
	return &CalorieGoalHandler{
		SetGoalUC:            setGoalUC,
		GetGoalUC:            getGoalUC,
		AdjustFromProgressUC: adjustFromProgressUC,
	}
}

type setGoalRequest struct {
	WeightKg      float64                 `json:"weightKg" binding:"required,gt=0"`
	HeightCm      float64                 `json:"heightCm" binding:"required,gt=0"`
	AgeYears      int                     `json:"ageYears" binding:"required,gt=0"`
	IsMale        bool                    `json:"isMale"`
	ActivityLevel models.ActivityLevel    `json:"activityLevel"`
	FitnessGoal   models.FitnessGoalType  `json:"fitnessGoal"`
}

// SetCalorieGoal godoc
// @Summary      Set personalized calorie & macro targets
// @Description  Calculates TDEE using Mifflin-St Jeor + activity level and stores the result.
// @Description  activityLevel: sedentary | light | moderate | active | very_active
// @Description  fitnessGoal:   lose_weight | maintain | gain_mass
// @Tags         CalorieGoals
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      setGoalRequest  true  "Physical data + goal"
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Router       /private/goals/calories [post]
func (h *CalorieGoalHandler) SetCalorieGoal(c *gin.Context) {
	var req setGoalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(uint)

	result, err := h.SetGoalUC.Execute(usecases.SetUserCalorieGoalInput{
		UserID:        userID,
		WeightKg:      req.WeightKg,
		HeightCm:      req.HeightCm,
		AgeYears:      req.AgeYears,
		IsMale:        req.IsMale,
		ActivityLevel: req.ActivityLevel,
		FitnessGoal:   req.FitnessGoal,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

// GetCalorieGoal godoc
// @Summary      Get current calorie goal
// @Description  Returns the user's stored calorie and macro targets with AI adjustment metadata.
// @Tags         CalorieGoals
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /private/goals/calories [get]
func (h *CalorieGoalHandler) GetCalorieGoal(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	goal, err := h.GetGoalUC.Execute(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no calorie goal set for this user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": goal})
}

// AdjustCaloriesFromProgress godoc
// @Summary      Auto-adjust calories from progress-service data
// @Description  Fetches real weight metrics from progress-service and applies adaptive
// @Description  calorie/macro adjustments based on the user's fitness goal and actual progress.
// @Description  This is the AI integration endpoint — call it periodically (e.g. weekly).
// @Tags         CalorieGoals
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /private/goals/calories/adjust [post]
func (h *CalorieGoalHandler) AdjustCaloriesFromProgress(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	result, err := h.AdjustFromProgressUC.Execute(userID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error()[:13] == "no calorie go" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}
