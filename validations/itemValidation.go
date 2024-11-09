package validations

import (
	"errors"
	"log"

	"github.com/Safiramdhn/project-app-inventaris-golang-safira/models"
)

func ValidateItemInput(item models.Item) error {
	if item.CategoryID == 0 {
		log.Printf("invalid category ID %d", item.CategoryID)
		return errors.New("invalid category id")
	}
	if item.Name == "" {
		log.Println("item name is required")
		return errors.New("item name is required")
	}
	if item.Price == 0 {
		log.Println("item price is required")
		return errors.New("item price is required")
	}
	if item.PurchaseDate.IsZero() {
		log.Println("item purchase date is required")
		return errors.New("item purchase date is required")
	}
	return nil
}
