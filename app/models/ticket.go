package models

import (
	"context"
	"database/sql"
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

// CreateTicket handles the request to create a new ticket.
func (h *TicketHandler) CreateTicket(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	var req struct {
		Title       string `json:"Title"`
		Description string `json:"Description"`
		Status      string `json:"Status"`
		Priority    string `json:"Priority"`
		RequesterID int64  `json:"RequesterID"`
		AssigneeID  *int64 `json:"AssigneeID"` // Use pointer to int64 to handle null
	}

	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Status(r, http.StatusBadRequest)
		renderer.PrettyJSON(w, r, err.Error())
		return
	}

	ticket := models.Ticket{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		Priority:    req.Priority,
		RequesterID: req.RequesterID,
	}

	if req.AssigneeID != nil {
		ticket.AssigneeID = sql.NullInt64{Int64: *req.AssigneeID, Valid: true}
	} else {
		ticket.AssigneeID = sql.NullInt64{Valid: false}
	}

	if err := models.CreateTicket(h.db, ctx, &ticket); err != nil {
		render.Status(r, http.StatusInternalServerError)
		renderer.PrettyJSON(w, r, err.Error())
		return
	}

	render.Status(r, http.StatusCreated)
	renderer.PrettyJSON(w, r, ticket)
}
