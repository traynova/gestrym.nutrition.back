package routes

import (
	"gestrym-nutrition/docs"
	"gestrym-nutrition/src/common/middleware"
	"gestrym-nutrition/src/common/utils"
	"net/http"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"gestrym-nutrition/src/common/config"
	nutritionUseCases "gestrym-nutrition/src/nutrition/application/usecases"
	"gestrym-nutrition/src/nutrition/infrastructure/adapters"
	nutritionAdapters "gestrym-nutrition/src/nutrition/infrastructure/adapters"
	nutritionRepos "gestrym-nutrition/src/nutrition/infrastructure/repositories"
	nutritionHandlers "gestrym-nutrition/src/nutrition/interfaces/http/handlers"
)

type routesDefinition struct {
	serverGroup    *gin.RouterGroup
	publicGroup    *gin.RouterGroup
	privateGroup   *gin.RouterGroup
	internalGroup  *gin.RouterGroup
	protectedGroup *gin.RouterGroup
	logger         utils.ILogger
}

var (
	routesInstance *routesDefinition
	routesOnce     sync.Once
)

func NewRoutesDefinition(serverInstance *gin.Engine) *routesDefinition {
	routesOnce.Do(func() {
		routesInstance = &routesDefinition{}
		routesInstance.logger = utils.NewLogger()
		docs.SwaggerInfo.Title = "Gestrym Nutrition API"
		docs.SwaggerInfo.Description = "API para el manejo de nutrición."
		docs.SwaggerInfo.Version = "1.0"
		docs.SwaggerInfo.BasePath = "/gestrym-nutrition"
		routesInstance.addCORSConfig(serverInstance)
		routesInstance.addRoutes(serverInstance)
	})
	return routesInstance
}

func (r *routesDefinition) addCORSConfig(serverInstance *gin.Engine) {
	corsMiddleware := cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-API-Key"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
	serverInstance.Use(corsMiddleware)
}

func (r *routesDefinition) addRoutes(serverInstance *gin.Engine) {
	r.addDefaultRoutes(serverInstance)

	// Instantiate DB
	dbConn := config.NewPostgresConnection()
	db := dbConn.GetDB()

	// ── Repositories ─────────────────────────────────────────────────────────
	foodRepo := nutritionRepos.NewFoodRepositoryImpl(db)
	mealPlanRepo := nutritionRepos.NewMealPlanRepositoryImpl(db)
	mealDayRepo := nutritionRepos.NewMealDayRepositoryImpl(db)
	mealItemRepo := nutritionRepos.NewMealItemRepositoryImpl(db)
	nutritionLogRepo := nutritionRepos.NewNutritionLogRepositoryImpl(db)
	calorieGoalRepo := nutritionRepos.NewUserCalorieGoalRepositoryImpl(db)

	// ── Adapters & Services ──────────────────────────────────────────────────
	usdaAdapter := nutritionAdapters.NewUSDAAdapterImpl("", viper.GetString("USDA_API_KEY"))
	pexelsAdapter := nutritionAdapters.NewPexelsAdapterImpl(viper.GetString("PEXELS_API_KEY"))
	storageAdapter := adapters.NewFileStorageAdapterImpl(viper.GetString("STORAGE_SERVICE_URL"), viper.GetString("STORAGE_SERVICE_API_KEY"))
	storageService := nutritionAdapters.NewStorageServiceAdapterImpl(storageAdapter)
	progressAdapter := nutritionAdapters.NewProgressServiceAdapterImpl(viper.GetString("PROGRESS_SERVICE_URL"), viper.GetString("PROGRESS_SERVICE_API_KEY"))

	// ── Use Cases: Food ───────────────────────────────────────────────────────
	searchFoodsUC := nutritionUseCases.NewSearchFoodsUseCase(foodRepo)
	getFoodByIDUC := nutritionUseCases.NewGetFoodByIDUseCase(foodRepo)
	importFoodsUC := nutritionUseCases.NewImportFoodsWithImagesUseCase(foodRepo, usdaAdapter, pexelsAdapter, storageService)

	// ── Use Cases: Meal Plans ─────────────────────────────────────────────────
	createMealPlanUC := nutritionUseCases.NewCreateMealPlanUseCase(mealPlanRepo)
	getMealPlanUC := nutritionUseCases.NewGetMealPlanUseCase(mealPlanRepo)
	getUserMealPlansUC := nutritionUseCases.NewGetUserMealPlansUseCase(mealPlanRepo)
	addMealDayUC := nutritionUseCases.NewAddMealDayUseCase(mealPlanRepo, mealDayRepo)
	addMealItemUC := nutritionUseCases.NewAddMealItemUseCase(mealDayRepo, mealItemRepo, foodRepo)

	// ── Use Cases: Nutrition Tracking ─────────────────────────────────────────
	logFoodIntakeUC := nutritionUseCases.NewLogFoodIntakeUseCase(nutritionLogRepo, foodRepo)
	getDailyNutritionUC := nutritionUseCases.NewGetDailyNutritionUseCase(nutritionLogRepo, mealPlanRepo)
	getNutritionHistoryUC := nutritionUseCases.NewGetNutritionHistoryUseCase(nutritionLogRepo)

	// ── Use Cases: Calorie Goals ──────────────────────────────────────────────
	setCalorieGoalUC := nutritionUseCases.NewSetUserCalorieGoalUseCase(calorieGoalRepo)
	getCalorieGoalUC := nutritionUseCases.NewGetUserCalorieGoalUseCase(calorieGoalRepo)
	adjustCaloriesUC := nutritionUseCases.NewAdjustCaloriesFromProgressUseCase(calorieGoalRepo, progressAdapter)

	// ── Handlers ──────────────────────────────────────────────────────────────
	foodHandler := nutritionHandlers.NewFoodHandler(searchFoodsUC, getFoodByIDUC, importFoodsUC)
	mealPlanHandler := nutritionHandlers.NewMealPlanHandler(
		createMealPlanUC,
		getMealPlanUC,
		getUserMealPlansUC,
		addMealDayUC,
		addMealItemUC,
	)
	nutritionLogHandler := nutritionHandlers.NewNutritionLogHandler(
		logFoodIntakeUC,
		getDailyNutritionUC,
		getNutritionHistoryUC,
	)
	calorieGoalHandler := nutritionHandlers.NewCalorieGoalHandler(
		setCalorieGoalUC,
		getCalorieGoalUC,
		adjustCaloriesUC,
	)

	// ── Router Groups ─────────────────────────────────────────────────────────
	r.serverGroup = serverInstance.Group(docs.SwaggerInfo.BasePath)
	r.serverGroup.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.publicGroup = r.serverGroup.Group("/public")
	r.privateGroup = r.serverGroup.Group("/private")
	r.protectedGroup = r.serverGroup.Group("/protected")

	// Middleware
	r.privateGroup.Use(middleware.SetupJWTMiddleware())
	r.protectedGroup.Use(middleware.SetupApiKeyMiddleware())

	// ── Public Routes ─────────────────────────────────────────────────────────
	foodsGroup := r.publicGroup.Group("/foods")
	{
		foodsGroup.GET("", foodHandler.SearchFoods)
		foodsGroup.GET("/:id", foodHandler.GetFoodByID)
		foodsGroup.POST("/import", foodHandler.ImportFoods)
	}

	// ── Private Routes: Meal Plans ────────────────────────────────────────────
	mealPlansGroup := r.privateGroup.Group("/meal-plans")
	{
		mealPlansGroup.POST("", mealPlanHandler.CreateMealPlan)
		mealPlansGroup.GET("/:id", mealPlanHandler.GetMealPlan)
		mealPlansGroup.GET("/user/:userId", mealPlanHandler.GetUserMealPlans)
		mealPlansGroup.POST("/:id/days", mealPlanHandler.AddMealDay)
		mealPlansGroup.POST("/:id/items", mealPlanHandler.AddMealItem)
	}

	// ── Private Routes: Nutrition Tracking ───────────────────────────────────
	logsGroup := r.privateGroup.Group("/logs")
	{
		logsGroup.POST("", nutritionLogHandler.LogFood)
		logsGroup.GET("", nutritionLogHandler.GetDailyNutrition)
		logsGroup.GET("/history", nutritionLogHandler.GetNutritionHistory)
	}

	// ── Private Routes: Calorie Goals ────────────────────────────────────────
	goalsGroup := r.privateGroup.Group("/goals")
	{
		goalsGroup.POST("/calories", calorieGoalHandler.SetCalorieGoal)
		goalsGroup.GET("/calories", calorieGoalHandler.GetCalorieGoal)
		goalsGroup.POST("/calories/adjust", calorieGoalHandler.AdjustCaloriesFromProgress)
	}

	r.addPublicRoutes()
	r.addInternalRoutes()
	r.addProtectedRoutes()
}

func (r *routesDefinition) addDefaultRoutes(serverInstance *gin.Engine) {
	serverInstance.GET("/", func(cnx *gin.Context) {
		cnx.JSON(http.StatusOK, map[string]interface{}{
			"code":    "OK",
			"message": "gestrym-nutrition OK...",
			"date":    utils.GetCurrentTime(),
		})
	})

	serverInstance.NoRoute(func(cnx *gin.Context) {
		cnx.JSON(http.StatusNotFound, map[string]interface{}{
			"code":    "NOT_FOUND",
			"message": "Resource not found",
			"date":    utils.GetCurrentTime(),
		})
	})
}

func (r *routesDefinition) addPublicRoutes() {}

func (r *routesDefinition) addPrivateRoutes() {}

func (r *routesDefinition) addInternalRoutes() {}

func (r *routesDefinition) addProtectedRoutes() {}
