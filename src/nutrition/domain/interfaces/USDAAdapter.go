package interfaces

import "gestrym-nutrition/src/common/models"

type USDAAdapter interface {
	SearchFoods(query string) ([]models.Food, error)
}
