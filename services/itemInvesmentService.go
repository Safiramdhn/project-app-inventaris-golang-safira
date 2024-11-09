package services

import (
	"github.com/Safiramdhn/project-app-inventaris-golang-safira/models"
	"github.com/Safiramdhn/project-app-inventaris-golang-safira/repositories"
)

type ItemInvestmentService struct {
	ItemInvestmentRepo repositories.ItemInvestmentRepository
}

func NewItemInvestmentService(repo repositories.ItemInvestmentRepository) *ItemInvestmentService {
	return &ItemInvestmentService{ItemInvestmentRepo: repo}
}

func (s *ItemInvestmentService) GetByItemID(itemId int) (models.ItemInvestment, error) {
	return s.ItemInvestmentRepo.FindByItemId(itemId)
}

func (s *ItemInvestmentService) CountAllItemInvestments() (*models.ItemInvestment, error) {
	return s.ItemInvestmentRepo.CountAll()
}
