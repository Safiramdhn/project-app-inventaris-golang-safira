package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Safiramdhn/project-app-inventaris-golang-safira/models"
	"github.com/Safiramdhn/project-app-inventaris-golang-safira/services"
	"github.com/go-chi/chi/v5"
)

type CategoryHandler struct {
	CategoryService *services.CategoryService
}

func NewCategoryHandler(categoryService *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{CategoryService: categoryService}
}

func (hc *CategoryHandler) CreateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		JsonResp.SendError(w, http.StatusMethodNotAllowed, "Method not allowed", r.Method)
		return
	}

	var categoryInput models.Category
	if err := json.NewDecoder(r.Body).Decode(&categoryInput); err != nil {
		JsonResp.SendError(w, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	category, err := hc.CategoryService.CreateCategory(categoryInput)
	if err != nil {
		JsonResp.SendError(w, http.StatusInternalServerError, "Failed to create category", err.Error())
		return
	}
	JsonResp.SendCreated(w, category, "Category created successfully")
}

func (hc *CategoryHandler) UpdateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		log.Println("Method not supported")
		JsonResp.SendError(w, http.StatusMethodNotAllowed, "Method not allowed", r.Method)
		return
	}

	id := chi.URLParam(r, "id")
	var categoryInput models.Category
	if err := json.NewDecoder(r.Body).Decode(&categoryInput); err != nil {
		JsonResp.SendError(w, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}
	categoryID, err := strconv.Atoi(id)
	if err != nil {
		JsonResp.SendError(w, http.StatusBadRequest, "Invalid category ID", err.Error())
		return
	}
	categoryInput.ID = categoryID

	category, err := hc.CategoryService.UpdateCategory(categoryInput)
	if err != nil {
		JsonResp.SendError(w, http.StatusInternalServerError, "Failed to update category", err.Error())
		return
	}
	JsonResp.SendSuccess(w, category, "Category updated successfully")
}

func (hc *CategoryHandler) DeleteCategoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		JsonResp.SendError(w, http.StatusMethodNotAllowed, "Method not allowed", r.Method)
	}

	id := chi.URLParam(r, "id")
	categoryID, err := strconv.Atoi(id)
	if err != nil {
		JsonResp.SendError(w, http.StatusBadRequest, "Invalid category ID", err.Error())
		return
	}

	err = hc.CategoryService.DeleteCategory(categoryID)
	if err != nil {
		JsonResp.SendError(w, http.StatusInternalServerError, "Failed to delete category", err.Error())
		return
	}
	JsonResp.SendSuccess(w, nil, "Category deleted successfully")
}

func (hc *CategoryHandler) GetCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		JsonResp.SendError(w, http.StatusMethodNotAllowed, "Method not allowed", r.Method)
		return
	}

	categories, err := hc.CategoryService.GetAllCategories()
	if err != nil {
		JsonResp.SendError(w, http.StatusInternalServerError, "Failed to get categories", err.Error())
		return
	}

	JsonResp.SendSuccess(w, categories, "Categories retrieved successfully")
}

func (hc *CategoryHandler) GetCategoryByIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		JsonResp.SendError(w, http.StatusMethodNotAllowed, "Method not allowed", r.Method)
	}

	id := chi.URLParam(r, "id")
	categoryID, err := strconv.Atoi(id)
	if err != nil {
		JsonResp.SendError(w, http.StatusBadRequest, "Invalid category ID", err.Error())
		return
	}

	category, err := hc.CategoryService.GetCategoryByID(categoryID)
	if err != nil {
		JsonResp.SendError(w, http.StatusInternalServerError, "Failed to get category", err.Error())
		return
	}
	JsonResp.SendSuccess(w, category, "Category retrieved successfully")
}
