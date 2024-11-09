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
	Delete(id int) (string, error)
	ReplaceReminder(threshold int) ([]models.Item, error)
	CreateItemInvestment(item *models.Item) error
}

type itemRepository struct {
	DB *sql.DB
}

func NewItemRepository(db *sql.DB) ItemRepository {
	return &itemRepository{DB: db}
}

// CreateItemInvestment implements ItemRepository.
func (i *itemRepository) CreateItemInvestment(item *models.Item) error {
	if item == nil {
		return errors.New("item cannot be nil")
	}

	const dateLayout = "2006-01-02"
	lastDepreciationDate := time.Now()
	formattedLastDepreciationDate := lastDepreciationDate.Format(dateLayout)
	currentValue := item.Price - (item.Price * (float64(item.DepreciatedRate) / 100.0))

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

	sqlStatement := `INSERT INTO item_investments (item_id, initial_price, current_value, last_depreciation_date) VALUES ($1, $2, $3, $4)`
	_, err = tx.Exec(sqlStatement, item.ID, item.Price, currentValue, formattedLastDepreciationDate)
	if err != nil {
		log.Printf("Error inserting item investment: %v", err)
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v", err)
		return err
	}
	return nil
}

func (i *itemRepository) Create(itemInput *models.Item) (*models.Item, error) {
	if itemInput == nil {
		return nil, fmt.Errorf("itemInput cannot be nil")
	}

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

	sqlStatement := `INSERT INTO items (name, category_id, photo_url, price, purchase_date, depreciated_rate) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	err = tx.QueryRow(sqlStatement, itemInput.Name, itemInput.CategoryID, itemInput.PhotoURL, itemInput.Price, itemInput.PurchaseDate, itemInput.DepreciatedRate).Scan(&itemInput.ID)
	if err != nil {
		log.Printf("Error inserting item: %v", err)
		return nil, err
	}

	if err = i.CreateItemInvestment(itemInput); err != nil {
		log.Printf("Error creating item investment: %v", err)
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v", err)
		return nil, err
	}

	item, err := i.FindByID(itemInput.ID)
	if err != nil {
		log.Printf("Error finding item by ID after insert: %v", err)
		return nil, err
	}
	return item, nil
}

// Delete implements ItemRepository.
func (i *itemRepository) Delete(id int) (string, error) {
	tx, err := i.DB.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %v", err.Error())
		return "", err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // Re-panic after rollback
		} else if err != nil {
			log.Printf("Rolling back transaction due to error: %v", err.Error())
			tx.Rollback()
		}
	}()

	var photoUrl string
	sqlStatement := `UPDATE items SET status = 'deleted' WHERE id = $1 RETURNING photo_url`
	err = tx.QueryRow(sqlStatement, id).Scan(&photoUrl)
	if err != nil {
		return "", err
	}

	if err := tx.Commit(); err != nil {
		return "", err
	}
	return photoUrl, err
}

// FindAll implements ItemRepository.
func (i *itemRepository) FindAll() ([]models.Item, error) {
	sqlStatement := `SELECT i.id, i.name, c.name, i.photo_url, i.price, i.purchase_date, i.total_usage_days, i.is_replacement_needed, i.depreciated_rate FROM items i 
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
		err = rows.Scan(&item.ID, &item.Name, &item.CategoryName, &item.PhotoURL, &item.Price, &item.PurchaseDate, &item.TotalUsageDays, &item.IsReplacementNeeded, &item.DepreciatedRate)
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
		log.Printf("Error starting transaction: %v", err.Error())
		return nil, err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // Re-panic after rollback
		} else if err != nil {
			log.Printf("Rolling back transaction due to error: %v", err.Error())
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
	if itemInput.DepreciatedRate != 0 {
		fields["depreciated_rate"] = itemInput.DepreciatedRate
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

	sqlStatement := fmt.Sprintf("UPDATE items SET %s WHERE id = $%d AND status = 'active' RETURNING id", strings.Join(setClauses, ", "), index)
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

// ReplaceReminder implements ItemRepository.
func (i *itemRepository) ReplaceReminder(threshold int) ([]models.Item, error) {
	tx, err := i.DB.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %v", err.Error())
		return nil, err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // Re-panic after rollback
		} else if err != nil {
			log.Printf("Rolling back transaction due to error: %v", err.Error())
			tx.Rollback()
		}
	}()

	items, err := i.FindAll()
	if err != nil {
		return nil, err
	}

	var needReplaceItems []models.Item
	currentDate := time.Now()
	for _, item := range items {
		if !item.PurchaseDate.IsZero() {
			daysSincePurchase := int(currentDate.Sub(item.PurchaseDate).Hours() / 24)
			item.TotalUsageDays = daysSincePurchase
		}

		if item.TotalUsageDays < threshold && item.IsReplacementNeeded {
			item.IsReplacementNeeded = false
		} else if item.TotalUsageDays > threshold {
			item.IsReplacementNeeded = true
		}
		updateStatement := `UPDATE items SET total_usage_days = $1, is_replacement_needed = $2 WHERE id = $3`
		_, err = tx.Exec(updateStatement, item.TotalUsageDays, item.IsReplacementNeeded, item.ID)
		if err != nil {
			return nil, err
		}
		needReplaceItems = append(needReplaceItems, item)
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return needReplaceItems, nil
}
