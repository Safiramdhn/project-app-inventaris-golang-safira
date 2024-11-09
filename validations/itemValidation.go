package validations

import (
	"errors"

	"github.com/Safiramdhn/project-app-inventaris-golang-safira/models"
)

func ValidateItemInput(item models.Item) error {
	if item.CategoryID == 0 {
		return errors.New("invalid category id")
	}
	if item.Name == "" {
		return errors.New("item name is required")
	}
	if item.Price == 0 {
		return errors.New("item price is required")
	}
	if item.PurchaseDate.IsZero() {
		return errors.New("item purchase date is required")
	}
	return nil
}
