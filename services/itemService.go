package services

import (
	"errors"

	"github.com/Safiramdhn/project-app-inventaris-golang-safira/models"
	"github.com/Safiramdhn/project-app-inventaris-golang-safira/repositories"
	"github.com/Safiramdhn/project-app-inventaris-golang-safira/validations"
)

type ItemService struct {
	ItemRepo repositories.ItemRepository
}

func NewItemService(repo repositories.ItemRepository) *ItemService {
	return &ItemService{ItemRepo: repo}
}

func (s *ItemService) CreateItem(itemInput models.Item) (*models.Item, error) {
	err := validations.ValidateItemInput(itemInput)
	if err != nil {
		return nil, err
	}

	return s.ItemRepo.Create(&itemInput)
}

func (s *ItemService) GetItemsByID(id int) (*models.Item, error) {
	if id == 0 {
		return nil, errors.New("invalid id")
	}
	return s.ItemRepo.FindByID(id)
}

func (s *ItemService) UpdateItem(itemInput models.Item) (*models.Item, error) {
	if itemInput.ID == 0 {
		return nil, errors.New("invalid id")
	}
	err := validations.ValidateItemInput(itemInput)
	if err != nil {
		return nil, err
	}

	return s.ItemRepo.Update(&itemInput)
}

func (s *ItemService) DeleteItem(id int) (string, error) {
	if id == 0 {
		return "", errors.New("invalid id")
	}
	return s.ItemRepo.Delete(id)
}

func (s *ItemService) GetAllItems() ([]models.Item, error) {
	return s.ItemRepo.FindAll()
}

func (s *ItemService) GetReplacementItems() ([]models.Item, error) {
	threshold := 100
	return s.ItemRepo.ReplaceReminder(threshold)
}
