package validations

import (
	"errors"

	"github.com/Safiramdhn/project-app-inventaris-golang-safira/models"
)

func ValidateCategoryInput(categoryInput *models.Category) error {
	if categoryInput.Name == "" {
		return errors.New("name is required")
	}
	if categoryInput.Description == "" {
		return errors.New("description is required")
	}
	return nil
}
