package models

import (
	"context"
	"net/http"

	"github.com/go-chi/render"
	"github.com/uptrace/bun"
	"goat/app/renderer"
	"goat/services/models"
)

type TicketHandler struct {
	db *bun.DB
}

func NewTicketHandler(db *bun.DB) *TicketHandler {
	return &TicketHandler{db: db}
}

// ListTickets handles the request to list all tickets.
func (h *TicketHandler) ListTickets(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()

	tickets, err := models.ListTickets(h.db, ctx)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		renderer.PrettyJSON(w, r, err.Error())
		return
	}
	render.Status(r, http.StatusOK)
	renderer.PrettyJSON(w, r, tickets)
}
