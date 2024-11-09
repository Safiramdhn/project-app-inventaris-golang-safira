package repositories

import (
	"database/sql"

	"github.com/Safiramdhn/project-app-inventaris-golang-safira/models"
)

type ItemInvestmentRepository interface {
	CountAll() (*models.ItemInvestment, error)
	FindByItemId(id int) (models.ItemInvestment, error)
}

type itemInvestmentRepository struct {
	DB *sql.DB
}

func NewItemInvestmentRepository(db *sql.DB) ItemInvestmentRepository {
	return &itemInvestmentRepository{DB: db}
}

// FindAll implements ItemInvestmentRepository.
func (i *itemInvestmentRepository) CountAll() (*models.ItemInvestment, error) {
	var itemInvestment models.ItemInvestment
	sqlStatement := `SELECT SUM(initial_price) AS total_investment, SUM(current_value) AS depreciated_value FROM item_investments`
	err := i.DB.QueryRow(sqlStatement).Scan(&itemInvestment.TotalInvestment, &itemInvestment.DepricatedValue)
	if err != nil {
		return nil, err
	}
	return &itemInvestment, nil
}

// FindByItemId implements ItemInvestmentRepository.
func (i *itemInvestmentRepository) FindByItemId(itemId int) (models.ItemInvestment, error) {
	sqlStatement := `SELECT i.id, i.name, i.depreciated_rate, inv.initial_price, inv.current_value FROM item_investments inv
				JOIN items i ON inv.item_id = i.id WHERE inv.item_id = $1`
	var itemInvestment models.ItemInvestment
	err := i.DB.QueryRow(sqlStatement, itemId).Scan(&itemInvestment.ItemID, &itemInvestment.ItemName, &itemInvestment.DepreciationRate, &itemInvestment.InitialPrice, &itemInvestment.CurrentValue)
	if err == sql.ErrNoRows {
		return itemInvestment, nil
	} else if err != nil {
		return itemInvestment, err
	}

	return itemInvestment, nil
}
