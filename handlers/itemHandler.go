package handlers

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Safiramdhn/project-app-inventaris-golang-safira/models"
	"github.com/Safiramdhn/project-app-inventaris-golang-safira/services"
	"github.com/go-chi/chi/v5"
)

type ItemHandler struct {
	ItemService *services.ItemService
}

func NewItemHandler(service *services.ItemService) *ItemHandler {
	return &ItemHandler{ItemService: service}
}

func (hi *ItemHandler) CreateItemHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure request is a POST method
	if r.Method != http.MethodPost {
		JsonResp.SendError(w, http.StatusMethodNotAllowed, "Method not allowed", r.Method)
		return
	}

	// Parse form with max memory limit for file uploads
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		JsonResp.SendError(w, http.StatusBadRequest, "Unable to parse form", err.Error())
		return
	}

	// Extract form values with validation
	itemName := r.FormValue("name")
	if itemName == "" {
		JsonResp.SendError(w, http.StatusBadRequest, "Item name is required", nil)
		return
	}

	// Parse category ID
	categoryId := r.FormValue("category_id")
	itemCategoryId, err := strconv.Atoi(categoryId)
	if err != nil {
		JsonResp.SendError(w, http.StatusBadRequest, "Invalid category ID", err.Error())
		return
	}

	// Parse item price
	itemPrice, err := strconv.ParseFloat(r.FormValue("price"), 64)
	if err != nil {
		JsonResp.SendError(w, http.StatusBadRequest, "Invalid price", err.Error())
		return
	}

	// Parse purchase date
	itemPurchaseDate := r.FormValue("purchase_date")
	const dateLayout = "2006-01-02"
	formattedItemPurchaseDate, err := time.Parse(dateLayout, itemPurchaseDate)
	if err != nil {
		JsonResp.SendError(w, http.StatusBadRequest, "Invalid date format. Please use YYYY-MM-DD.", err.Error())
		return
	}

	depreciatedRate := r.FormValue("depreciated_rate")
	itemDepreciatedRate, err := strconv.Atoi(depreciatedRate)
	if err != nil {
		JsonResp.SendError(w, http.StatusBadRequest, "Invalid depreciated rate", err.Error())
		return
	}

	// Handle file upload
	file, fileHeader, err := r.FormFile("photo")
	if err != nil {
		JsonResp.SendError(w, http.StatusBadRequest, "Unable to get file from form", err.Error())
		return
	}
	defer file.Close()

	// Define upload path and ensure directory exists
	uploadPath := "./uploads"
	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		JsonResp.SendError(w, http.StatusInternalServerError, "Failed to create upload directory", err.Error())
		return
	}

	// Save the uploaded file
	filePath := filepath.Join(uploadPath, fileHeader.Filename)
	out, err := os.Create(filePath)
	if err != nil {
		JsonResp.SendError(w, http.StatusInternalServerError, "Unable to save file", err.Error())
		return
	}
	defer out.Close()

	// Copy uploaded file content to destination file
	if _, err := io.Copy(out, file); err != nil {
		JsonResp.SendError(w, http.StatusInternalServerError, "Failed to copy file content", err.Error())
		return
	}
	filePathURL := strings.ReplaceAll(filePath, "\\", "/")

	// Initialize item data
	itemInput := models.Item{
		Name:            itemName,
		CategoryID:      itemCategoryId,
		Price:           itemPrice,
		PurchaseDate:    formattedItemPurchaseDate,
		PhotoURL:        filePathURL,
		DepreciatedRate: itemDepreciatedRate,
	}

	// Call service to create item
	item, err := hi.ItemService.CreateItem(itemInput)
	if err != nil {
		JsonResp.SendError(w, http.StatusInternalServerError, "Failed to create item", err.Error())
		return
	}

	// Send success response
	JsonResp.SendCreated(w, item, "Item created successfully")
}

func (hi *ItemHandler) GetItemByIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		JsonResp.SendError(w, http.StatusMethodNotAllowed, "Method not allowed", r.Method)
		return
	}

	id := chi.URLParam(r, "id")
	itemId, err := strconv.Atoi(id)
	if err != nil {
		JsonResp.SendError(w, http.StatusBadRequest, "Invalid item ID", err.Error())
		return
	}

	item, err := hi.ItemService.GetItemsByID(itemId)
	if err != nil {
		JsonResp.SendError(w, http.StatusInternalServerError, "Failed to get item", err.Error())
		return
	}
	JsonResp.SendSuccess(w, item, "Item retrieved successfully")
}

func (hi *ItemHandler) UpdateItemHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		JsonResp.SendError(w, http.StatusMethodNotAllowed, "Method not allowed", r.Method)
		return
	}

	id := chi.URLParam(r, "id")
	itemId, err := strconv.Atoi(id)
	if err != nil {
		JsonResp.SendError(w, http.StatusBadRequest, "Invalid item ID", err.Error())
		return
	}

	// Parse form with max memory limit for file uploads
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		JsonResp.SendError(w, http.StatusBadRequest, "Unable to parse form", err.Error())
		return
	}

	itemName := r.FormValue("name")
	if itemName == "" {
		JsonResp.SendError(w, http.StatusBadRequest, "Item name is required", nil)
		return
	}

	// Parse category ID
	categoryId := r.FormValue("category_id")
	itemCategoryId, err := strconv.Atoi(categoryId)
	if err != nil {
		JsonResp.SendError(w, http.StatusBadRequest, "Invalid category ID", err.Error())
		return
	}

	// Parse item price
	itemPrice, err := strconv.ParseFloat(r.FormValue("price"), 64)
	if err != nil {
		JsonResp.SendError(w, http.StatusBadRequest, "Invalid price", err.Error())
		return
	}

	// Parse purchase date
	itemPurchaseDate := r.FormValue("purchase_date")
	const dateLayout = "2006-01-02"
	formattedItemPurchaseDate, err := time.Parse(dateLayout, itemPurchaseDate)
	if err != nil {
		JsonResp.SendError(w, http.StatusBadRequest, "Invalid date format. Please use YYYY-MM-DD.", err.Error())
		return
	}

	depreciatedRate := r.FormValue("depreciated_rate")
	itemDepreciatedRate, err := strconv.Atoi(depreciatedRate)
	if err != nil {
		JsonResp.SendError(w, http.StatusBadRequest, "Invalid depreciated rate", err.Error())
		return
	}

	// Handle file upload
	file, fileHeader, err := r.FormFile("photo")
	if err != nil {
		JsonResp.SendError(w, http.StatusBadRequest, "Unable to get file from form", err.Error())
		return
	}
	defer file.Close()

	// Define upload path and ensure directory exists
	uploadPath := "./uploads"
	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		JsonResp.SendError(w, http.StatusInternalServerError, "Failed to create upload directory", err.Error())
		return
	}

	// Save the uploaded file
	filePath := filepath.Join(uploadPath, fileHeader.Filename)
	out, err := os.Create(filePath)
	if err != nil {
		JsonResp.SendError(w, http.StatusInternalServerError, "Unable to save file", err.Error())
		return
	}
	defer out.Close()

	// Copy uploaded file content to destination file
	if _, err := io.Copy(out, file); err != nil {
		JsonResp.SendError(w, http.StatusInternalServerError, "Failed to copy file content", err.Error())
		return
	}
	filePathURL := strings.ReplaceAll(filePath, "\\", "/")

	// Initialize item data
	itemInput := models.Item{
		ID:              itemId,
		Name:            itemName,
		CategoryID:      itemCategoryId,
		Price:           itemPrice,
		PurchaseDate:    formattedItemPurchaseDate,
		PhotoURL:        filePathURL,
		DepreciatedRate: itemDepreciatedRate,
	}

	// Call service to update item
	item, err := hi.ItemService.UpdateItem(itemInput)
	if err != nil {
		JsonResp.SendError(w, http.StatusInternalServerError, "Failed to update item", err.Error())
		return
	}
	JsonResp.SendSuccess(w, item, "Item updated successfully")
}

func (hi *ItemHandler) DeleteItemHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		JsonResp.SendError(w, http.StatusMethodNotAllowed, "Method not allowed", r.Method)
		return
	}

	// Get item ID from URL parameter
	id := chi.URLParam(r, "id")
	itemId, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("Error parsing item ID: %v", err.Error())
		JsonResp.SendError(w, http.StatusBadRequest, "Invalid item ID", err.Error())
		return
	}

	// Call service to delete item
	photoUrl, err := hi.ItemService.DeleteItem(itemId)
	if err != nil {
		JsonResp.SendError(w, http.StatusInternalServerError, "Failed to delete item", err.Error())
		return
	}

	err = os.Remove(photoUrl)
	if err != nil {
		JsonResp.SendError(w, http.StatusInternalServerError, "Failed to remove item", err.Error())
		return
	}

	JsonResp.SendSuccess(w, nil, "Item deleted successfully")
}

func (hi *ItemHandler) GetAllItemsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		JsonResp.SendError(w, http.StatusMethodNotAllowed, "Method not allowed", r.Method)
		return
	}
	items, err := hi.ItemService.GetAllItems()
	if err != nil {
		JsonResp.SendError(w, http.StatusInternalServerError, "Failed to get items", err.Error())
		return
	}

	JsonResp.SendSuccess(w, items, "Items retrieved successfully")
}

func (hi *ItemHandler) GetReplacementItemsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		JsonResp.SendError(w, http.StatusMethodNotAllowed, "Method not allowed", r.Method)
		return
	}

	items, err := hi.ItemService.GetReplacementItems()
	if err != nil {
		JsonResp.SendError(w, http.StatusInternalServerError, "Failed to get replacement items", err.Error())
		return
	}
	if len(items) == 0 {
		JsonResp.SendSuccess(w, nil, "No replacement items available")
	} else {
		JsonResp.SendSuccess(w, items, "Replacement items retrieved successfully")
	}
}
