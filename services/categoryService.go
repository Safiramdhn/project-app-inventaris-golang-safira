package services

import (
	"errors"
	"log"

	"github.com/Safiramdhn/project-app-inventaris-golang-safira/models"
	"github.com/Safiramdhn/project-app-inventaris-golang-safira/repositories"
	"github.com/Safiramdhn/project-app-inventaris-golang-safira/validations"
)

type CategoryService struct {
	CategoryRepo repositories.CategoryRepository
}

func NewCategoryService(repo repositories.CategoryRepository) *CategoryService {
	return &CategoryService{CategoryRepo: repo}
}

func (cs *CategoryService) CreateCategory(categoryInput models.Category) (*models.Category, error) {
	// Validate input
	if err := validations.ValidateCategoryInput(&categoryInput); err != nil {
		return nil, err
	}

	// Attempt to create the category
	category, err := cs.CategoryRepo.Create(&categoryInput)
	if err != nil {
		log.Printf("Failed to create category: %v", err) // Log the error
		return nil, err
	}

	return category, nil
}

func (cs *CategoryService) UpdateCategory(categoryInput models.Category) (*models.Category, error) {
	// Validate input
	if categoryInput.ID == 0 {
		return nil, errors.New("category ID is required")
	}
	if err := validations.ValidateCategoryInput(&categoryInput); err != nil {
		return nil, err
	}

	// Attempt to update the category
	category, err := cs.CategoryRepo.Update(&categoryInput)
	if err != nil {
		log.Printf("Failed to update category: %v", err) // Log the error
		return nil, err
	}

	return category, nil
}

func (cs *CategoryService) DeleteCategory(id int) error {
	if id <= 0 {
		return errors.New("invalid category id")
	}

	// Attempt to delete the category
	err := cs.CategoryRepo.Delete(id)
	if err != nil {
		log.Printf("Failed to delete category: %v", err) // Log the error
		return err
	}

	return nil
}

func (cs *CategoryService) GetCategoryByID(id int) (*models.Category, error) {
	// Attempt to get the category by ID
	if id <= 0 {
		return nil, errors.New("invalid category id")
	}
	category, err := cs.CategoryRepo.FindById(id)
	if err != nil {
		log.Printf("Failed to get category by ID: %v", err) // Log the error
		return nil, err
	}
	return category, nil
}

func (cs *CategoryService) GetAllCategories() ([]models.Category, error) {
	// Attempt to get all categories
	return cs.CategoryRepo.FindAll()
}
