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

type CategoryRepository interface {
	Create(categoryInput *models.Category) (*models.Category, error)
	Update(categoryInput *models.Category) (*models.Category, error)
	Delete(id int) error
	FindAll() ([]models.Category, error)
	FindByID(id int) (*models.Category, error)
}

type categoryRepository struct {
	DB *sql.DB
}

// NewCategoryRepository creates a new instance of CategoryRepository
func NewCategoryRepository(db *sql.DB) CategoryRepository {
	return &categoryRepository{DB: db}
}

// Create implements CategoryRepository.
func (c *categoryRepository) Create(categoryInput *models.Category) (*models.Category, error) {
	tx, err := c.DB.Begin()
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

	sqlStatement := `INSERT INTO categories (name, description) VALUES ($1, $2) RETURNING id`
	err = tx.QueryRow(sqlStatement, categoryInput.Name, categoryInput.Description).Scan(&categoryInput.ID)
	if err != nil {
		log.Printf("Error inserting category: %v", err)
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v", err)
		return nil, err
	}

	log.Printf("Inserted category with ID: %d", categoryInput.ID) // Log success
	return categoryInput, nil
}

// Delete implements CategoryRepository.
func (c *categoryRepository) Delete(id int) error {
	tx, err := c.DB.Begin()
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

	sqlStatement := `UPDATE categories SET status = 'deleted' WHERE id = $1`
	_, err = tx.Exec(sqlStatement, id)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v", err)
		return err
	}

	return nil
}

// FindAll implements CategoryRepository.
func (c *categoryRepository) FindAll() ([]models.Category, error) {
	var categories []models.Category
	sqlStatement := `SELECT id, name, description FROM categories WHERE status = 'active'`
	rows, err := c.DB.Query(sqlStatement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var category models.Category
		err := rows.Scan(&category.ID, &category.Name, &category.Description)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

// FindByID implements CategoryRepository.
func (c *categoryRepository) FindByID(id int) (*models.Category, error) {
	var category models.Category
	sqlStatement := `SELECT id, name, description FROM categories WHERE id = $1 AND status = 'active'`
	err := c.DB.QueryRow(sqlStatement, id).Scan(&category.ID, &category.Name, &category.Description)
	if err == sql.ErrNoRows {
		return nil, errors.New("category does not exist")
	} else if err != nil {
		return nil, err
	}

	return &category, nil
}

// Update implements CategoryRepository.
func (c *categoryRepository) Update(categoryInput *models.Category) (*models.Category, error) {
	tx, err := c.DB.Begin()
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

	if categoryInput.Name != "" {
		fields["name"] = categoryInput.Name
	}

	if categoryInput.Description != "" {
		fields["description"] = categoryInput.Description
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

	sqlStatement := fmt.Sprintf("UPDATE categories SET %s WHERE id = $%d AND status = 'active' RETURNING id, name, description",
		strings.Join(setClauses, ", "), index)
	values = append(values, categoryInput.ID)

	// Execute the update query and scan the result
	var updatedCategory models.Category
	err = tx.QueryRow(sqlStatement, values...).Scan(&updatedCategory.ID, &updatedCategory.Name, &updatedCategory.Description)
	if err != nil {
		log.Printf("Error updating category: %v", err)
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v", err)
		return nil, err
	}

	// Return the updated category
	return &updatedCategory, nil
}
