package handlers

import (
	"gestrym-nutrition/src/common/models"
	"gestrym-nutrition/src/nutrition/application/usecases"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type NutritionLogHandler struct {
	LogFoodUC    *usecases.LogFoodIntakeUseCase
	DailyUC      *usecases.GetDailyNutritionUseCase
	HistoryUC    *usecases.GetNutritionHistoryUseCase
}

func NewNutritionLogHandler(
	logFoodUC *usecases.LogFoodIntakeUseCase,
	dailyUC *usecases.GetDailyNutritionUseCase,
	historyUC *usecases.GetNutritionHistoryUseCase,
) *NutritionLogHandler {
	return &NutritionLogHandler{
		LogFoodUC: logFoodUC,
		DailyUC:   dailyUC,
		HistoryUC: historyUC,
	}
}

type logFoodRequest struct {
	FoodID   uint             `json:"foodId" binding:"required"`
	Quantity float64          `json:"quantity" binding:"required,gt=0"`
	MealType models.MealType  `json:"mealType"`
	Date     string           `json:"date"` // "YYYY-MM-DD", optional
	Notes    string           `json:"notes"`
}

// LogFood godoc
// @Summary      Log food intake
// @Description  Records a food item the user ate, calculates nutrition automatically.
// @Tags         NutritionTracking
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      logFoodRequest  true  "Food intake"
// @Success      201   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Router       /private/logs [post]
func (h *NutritionLogHandler) LogFood(c *gin.Context) {
	var req logFoodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(uint)

	// Parse date if provided
	var date time.Time
	if req.Date != "" {
		parsed, err := time.Parse("2006-01-02", req.Date)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format, use YYYY-MM-DD"})
			return
		}
		date = parsed
	}

	input := usecases.LogFoodIntakeInput{
		UserID:   userID,
		FoodID:   req.FoodID,
		Quantity: req.Quantity,
		MealType: req.MealType,
		Date:     date,
		Notes:    req.Notes,
	}

	log, err := h.LogFoodUC.Execute(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": log})
}

// GetDailyNutrition godoc
// @Summary      Get daily nutrition summary
// @Description  Returns totals, goals, and progress for a specific date.
// @Tags         NutritionTracking
// @Produce      json
// @Security     BearerAuth
// @Param        date  query  string  false  "Date (YYYY-MM-DD), defaults to today"
// @Success      200   {object}  map[string]interface{}
// @Router       /private/logs [get]
func (h *NutritionLogHandler) GetDailyNutrition(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	dateStr := c.DefaultQuery("date", time.Now().UTC().Format("2006-01-02"))
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format, use YYYY-MM-DD"})
		return
	}

	result, err := h.DailyUC.Execute(userID, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

// GetNutritionHistory godoc
// @Summary      Get nutrition history
// @Description  Returns paginated nutrition logs within a date range.
// @Tags         NutritionTracking
// @Produce      json
// @Security     BearerAuth
// @Param        start     query  string  true   "Start date (YYYY-MM-DD)"
// @Param        end       query  string  true   "End date (YYYY-MM-DD)"
// @Param        page      query  int     false  "Page number (default 1)"
// @Param        pageSize  query  int     false  "Page size (default 20, max 100)"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Router       /private/logs/history [get]
func (h *NutritionLogHandler) GetNutritionHistory(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	startStr := c.Query("start")
	endStr := c.Query("end")

	if startStr == "" || endStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start and end date parameters are required"})
		return
	}

	start, err := time.Parse("2006-01-02", startStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start date format, use YYYY-MM-DD"})
		return
	}

	end, err := time.Parse("2006-01-02", endStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end date format, use YYYY-MM-DD"})
		return
	}

	// Set end to end of day
	end = end.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	result, err := h.HistoryUC.Execute(userID, start, end, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":     result.Logs,
		"totals":   result.Totals,
		"total":    result.Total,
		"page":     result.Page,
		"pageSize": result.PageSize,
	})
}
