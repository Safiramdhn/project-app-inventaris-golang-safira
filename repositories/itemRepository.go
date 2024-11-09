package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Safiramdhn/project-app-inventaris-golang-safira/models"
)

type ItemRepository interface {
	FindAll() ([]models.Item, error)
	FindByID(id int) (*models.Item, error)
	Create(itemInput *models.Item) (*models.Item, error)
	Update(itemInput *models.Item) (*models.Item, error)
	Delete(id int) error
}

type itemRepository struct {
	DB *sql.DB
}

func NewItemRepository(db *sql.DB) ItemRepository {
	return &itemRepository{DB: db}
}

// Create implements ItemRepository.
func (i *itemRepository) Create(itemInput *models.Item) (*models.Item, error) {
	tx, err := i.DB.Begin()
	if err != nil {
		log.Printf("Error starting create transaction: %v\n", err)
		return nil, err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // Re-panic after rollback
		} else if err != nil {
			log.Printf("Rolling back transaction due to error: %v", err)
			tx.Rollback()
		}
	}()

	sqlStatement := `INSERT INTO items (name, category_id, photo_url, price, purchase_date) VALUES ($1, $2, $3, $4, $5)`
	result, err := tx.Exec(sqlStatement, itemInput.Name, itemInput.CategoryID, itemInput.PhotoURL, itemInput.Price, itemInput.PurchaseDate)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	itemId := int(id)
	item, err := i.FindByID(itemId)
	if err != nil {
		return nil, err
	}
	return item, nil
}

// Delete implements ItemRepository.
func (i *itemRepository) Delete(id int) error {
	tx, err := i.DB.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // Re-panic after rollback
		} else if err != nil {
			log.Printf("Rolling back transaction due to error: %v", err)
			tx.Rollback()
		}
	}()

	sqlStatement := `UPDATE items SET status = 'deleted' WHERE id = $1`
	_, err = tx.Exec(sqlStatement, id)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return err
}

// FindAll implements ItemRepository.
func (i *itemRepository) FindAll() ([]models.Item, error) {
	sqlStatement := `SELECT i.id, i.name, c.name, i.photo_url, i.price, i.purchase_date, i.total_usage_days FROM items i 
				JOIN categories c ON i.category_id = c.id 
				WHERE i.status = 'active'`
	rows, err := i.DB.Query(sqlStatement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.Item
	for rows.Next() {
		var item models.Item
		err = rows.Scan(&item.ID, &item.Name, &item.CategoryName, &item.PhotoURL, &item.Price, &item.PurchaseDate, &item.TotalUsageDays)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

// FindByID implements ItemRepository.
func (i *itemRepository) FindByID(id int) (*models.Item, error) {
	var item models.Item
	sqlStatement := `SELECT i.id, i.name, c.name, i.photo_url, i.price, i.purchase_date, i.total_usage_days FROM items i 
					JOIN categories c ON i.category_id = c.id 
					WHERE i.id = $1 AND i.status = 'active'`
	err := i.DB.QueryRow(sqlStatement, id).Scan(&item.ID, &item.Name, &item.CategoryName, &item.PhotoURL, &item.Price, &item.PurchaseDate, &item.TotalUsageDays)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &item, nil
}

// Update implements ItemRepository.
func (i *itemRepository) Update(itemInput *models.Item) (*models.Item, error) {
	tx, err := i.DB.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		return nil, err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // Re-panic after rollback
		} else if err != nil {
			log.Printf("Rolling back transaction due to error: %v", err)
			tx.Rollback()
		}
	}()

	fields := make(map[string]interface{})

	if itemInput.Name != "" {
		fields["name"] = itemInput.Name
	}
	if itemInput.CategoryID != 0 {
		fields["category_id"] = itemInput.CategoryID
	}
	if itemInput.PhotoURL != "" {
		fields["photo_url"] = itemInput.PhotoURL
	}
	if itemInput.Price != 0 {
		fields["price"] = itemInput.Price
	}
	if !itemInput.PurchaseDate.IsZero() {
		fields["purchase_date"] = itemInput.PurchaseDate
	}
	if itemInput.TotalUsageDays != 0 {
		fields["total_usage_days"] = itemInput.TotalUsageDays
	}

	fields["updated_at"] = time.Now()
	setClauses := []string{}
	values := []interface{}{}
	index := 1
	for field, value := range fields {
		setClauses = append(setClauses, field+"=$"+strconv.Itoa(index))
		values = append(values, value)
		index++
	}

	if len(setClauses) == 0 {
		return nil, errors.New("no fields to update")
	}

	sqlStatement := fmt.Sprintf("UPDATE items SET %s WHERE id = $%d AND status = 'active RETURNING id", strings.Join(setClauses, ", "), index)
	values = append(values, itemInput.ID)

	var id int
	err = tx.QueryRow(sqlStatement, values...).Scan(&id)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	updatedItem, err := i.FindByID(id)
	if err != nil {
		return nil, err
	}
	return updatedItem, nil
}
