package models

import (
	"context"
	dbconn "goat/db"
	"net/http"

	"github.com/uptrace/bun"
	"goat/app/renderer"
	"goat/services/models"

	"github.com/go-chi/render"
)

type CustomerHandler struct {
	db *bun.DB
}

func NewCustomerHandler(db *bun.DB) *CustomerHandler {
	return &CustomerHandler{db: dbconn.ConnectDB()}
}

func (h *CustomerHandler) ListCustomers(w http.ResponseWriter, r *http.Request) {
	customers, err := models.ListCustomers(h.db, context.Background())
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		renderer.PrettyJSON(w, r, err.Error())
		return
	}
	render.Status(r, http.StatusOK)
	renderer.PrettyJSON(w, r, customers)
}
