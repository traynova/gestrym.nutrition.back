package config

import (
	"fmt"
	"gestrym-nutrition/src/common/models"
	"gestrym-nutrition/src/common/utils"
)

var logger = utils.NewLogger()

func MigrateDB() (IDatabaseConnection, error) {
	connection := NewPostgresConnection()
	db := connection.GetDB()

	// Register all models for AutoMigrate
	err := db.AutoMigrate(
		&models.FoodCategory{},
		&models.Food{},
		&models.MealPlan{},
		&models.MealDay{},
		&models.MealItem{},
		&models.NutritionLog{},
		&models.UserCalorieGoal{},
	)

	if err != nil {
		logger.Error(fmt.Sprintf("[ERROR] Error al migrar las entidades: %s", err.Error()))
		return nil, err
	}

	logger.Info("[OK] Todas las migraciones completadas exitosamente")
	return connection, nil
}
