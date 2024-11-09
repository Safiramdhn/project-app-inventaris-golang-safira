package handlers

import (
	"net/http"
	"strconv"

	"github.com/Safiramdhn/project-app-inventaris-golang-safira/services"
	"github.com/go-chi/chi/v5"
)

type ItemInvestmentHandler struct {
	ItemInvestmentService services.ItemInvestmentService
}

func NewItemInvestmentHandler(service services.ItemInvestmentService) *ItemInvestmentHandler {
	return &ItemInvestmentHandler{ItemInvestmentService: service}
}

func (inh *ItemInvestmentHandler) CountAllItemInvestmentsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		JsonResp.SendError(w, http.StatusMethodNotAllowed, "Method not allowed", r.Method)
		return
	}

	itemInvesment, err := inh.ItemInvestmentService.CountAllItemInvestments()
	if err != nil {
		JsonResp.SendError(w, http.StatusInternalServerError, "Internal Server Error", err)
		return
	}
	JsonResp.SendSuccess(w, itemInvesment, "")
}

func (inh *ItemInvestmentHandler) GetItemInvesmentByItemIdHandler(w http.ResponseWriter, r *http.Request) {
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

	itemInvesment, err := inh.ItemInvestmentService.GetByItemID(itemId)
	if err != nil {
		JsonResp.SendError(w, http.StatusInternalServerError, "Internal Server Error", err)
		return
	}
	JsonResp.SendSuccess(w, itemInvesment, "")
}
