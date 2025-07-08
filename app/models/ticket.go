package models

import (
	"context"
	"database/sql"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"

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
	}

	// Check if the requester exists
	requester, err := models.GetUserByID(h.db, ctx, req.RequesterID)
	if err != nil {
		if err == sql.ErrNoRows {
			render.Status(r, http.StatusNotFound)
			renderer.PrettyJSON(w, r, "Requester not found")
			return
		}
		render.Status(r, http.StatusInternalServerError)
		renderer.PrettyJSON(w, r, err.Error())
		return
	}
	ticket.RequesterID = requester.ID

	if req.AssigneeID != nil {
		// Check if the assignee exists
		assignee, err := models.GetUserByID(h.db, ctx, *req.AssigneeID)
		if err != nil {
			if err == sql.ErrNoRows {
				render.Status(r, http.StatusNotFound)
				renderer.PrettyJSON(w, r, "Assignee not found")
				return
			}
			render.Status(r, http.StatusInternalServerError)
			renderer.PrettyJSON(w, r, err.Error())
			return
		}
		ticket.AssigneeID = sql.NullInt64{Int64: assignee.ID, Valid: true}
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

// GetTicket handles the request to get a ticket by ID.
func (h *TicketHandler) GetTicket(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		renderer.PrettyJSON(w, r, "Invalid ticket ID")
		return
	}

	ticket, err := models.GetTicketByID(h.db, ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			render.Status(r, http.StatusNotFound)
			renderer.PrettyJSON(w, r, "Ticket not found")
			return
		}
		render.Status(r, http.StatusInternalServerError)
		renderer.PrettyJSON(w, r, err.Error())
		return
	}

	render.Status(r, http.StatusOK)
	renderer.PrettyJSON(w, r, ticket)
}

// UpdateTicket handles the request to update an existing ticket.
func (h *TicketHandler) UpdateTicket(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		renderer.PrettyJSON(w, r, "Invalid ticket ID")
		return
	}

	var req struct {
		Title       string `json:"Title"`
		Description string `json:"Description"`
		Status      string `json:"Status"`
		Priority    string `json:"Priority"`
		RequesterID int64  `json:"RequesterID"`
		AssigneeID  *int64 `json:"AssigneeID"`
	}

	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Status(r, http.StatusBadRequest)
		renderer.PrettyJSON(w, r, err.Error())
		return
	}

	ticket := models.Ticket{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		Priority:    req.Priority,
	}

	// Check if the requester exists
	requester, err := models.GetUserByID(h.db, ctx, req.RequesterID)
	if err != nil {
		if err == sql.ErrNoRows {
			render.Status(r, http.StatusNotFound)
			renderer.PrettyJSON(w, r, "Requester not found")
			return
		}
		render.Status(r, http.StatusInternalServerError)
		renderer.PrettyJSON(w, r, err.Error())
		return
	}
	ticket.RequesterID = requester.ID

	if req.AssigneeID != nil {
		// Check if the assignee exists
		assignee, err := models.GetUserByID(h.db, ctx, *req.AssigneeID)
		if err != nil {
			if err == sql.ErrNoRows {
				render.Status(r, http.StatusNotFound)
				renderer.PrettyJSON(w, r, "Assignee not found")
				return
			}
			render.Status(r, http.StatusInternalServerError)
			renderer.PrettyJSON(w, r, err.Error())
			return
		}
		ticket.AssigneeID = sql.NullInt64{Int64: assignee.ID, Valid: true}
	} else {
		ticket.AssigneeID = sql.NullInt64{Valid: false}
	}

	if err := models.UpdateTicket(h.db, ctx, &ticket); err != nil {
		render.Status(r, http.StatusInternalServerError)
		renderer.PrettyJSON(w, r, err.Error())
		return
	}

	render.Status(r, http.StatusOK)
	renderer.PrettyJSON(w, r, ticket)
}
