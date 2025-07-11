package models

import (
	"context"
	"database/sql"
	"goat/app/middleware"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

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

func (h *TicketHandler) ListAgentTickets(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		render.Status(r, http.StatusUnauthorized)
		renderer.PrettyJSON(w, r, "Unauthorized")
		return
	}

	assigneeID, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		renderer.PrettyJSON(w, r, "Invalid user ID")
		return
	}

	tickets, err := models.ListTicketsByAssigneeID(h.db, r.Context(), assigneeID)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		renderer.PrettyJSON(w, r, err.Error())
		return
	}

	render.Status(r, http.StatusOK)
	renderer.PrettyJSON(w, r, tickets)
}

func (h *TicketHandler) ListOpenTickets(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	tickets, err := models.ListOpenTickets(h.db, ctx)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		renderer.PrettyJSON(w, r, err.Error())
		return
	}
	render.Status(r, http.StatusOK)
	renderer.PrettyJSON(w, r, tickets)
}

func (h *TicketHandler) GetAgentTicket(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		render.Status(r, http.StatusUnauthorized)
		renderer.PrettyJSON(w, r, "Unauthorized")
		return
	}

	assigneeID, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		renderer.PrettyJSON(w, r, "Invalid user ID")
		return
	}

	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		renderer.PrettyJSON(w, r, "Invalid ticket ID")
		return
	}

	ticket, err := models.GetTicketByID(h.db, r.Context(), id)
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

	if !ticket.AssigneeID.Valid || ticket.AssigneeID.Int64 != assigneeID {
		render.Status(r, http.StatusForbidden)
		renderer.PrettyJSON(w, r, "You are not authorized to view this ticket")
		return
	}

	// Agent can see all comments
	render.Status(r, http.StatusOK)
	renderer.PrettyJSON(w, r, ticket)
}

func (h *TicketHandler) UpdateAgentTicket(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		render.Status(r, http.StatusUnauthorized)
		renderer.PrettyJSON(w, r, "Unauthorized")
		return
	}

	assigneeID, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		renderer.PrettyJSON(w, r, "Invalid user ID")
		return
	}

	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		renderer.PrettyJSON(w, r, "Invalid ticket ID")
		return
	}

	existingTicket, err := models.GetTicketByID(h.db, r.Context(), id)
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

	var req struct {
		Status     *string `json:"status"`
		Priority   *string `json:"priority"`
		AssigneeID *int64  `json:"AssigneeID"`
	}

	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Status(r, http.StatusBadRequest)
		renderer.PrettyJSON(w, r, err.Error())
		return
	}

	// If the ticket is currently unassigned, or assigned to the current agent, allow updates.
	// If assigned to another agent, prevent updates unless the current agent is assigning it to themselves.
	canUpdate := true
	if existingTicket.AssigneeID.Valid && existingTicket.AssigneeID.Int64 != assigneeID {
		// Ticket is assigned to someone else
		if req.AssigneeID == nil || *req.AssigneeID != assigneeID {
			// If not trying to assign to self, forbid update
			canUpdate = false
		}
	}

	if !canUpdate {
		render.Status(r, http.StatusForbidden)
		renderer.PrettyJSON(w, r, "You are not authorized to update this ticket")
		return
	}

	if req.Status != nil && *req.Status != "" {
		existingTicket.Status = *req.Status
	}
	if req.Priority != nil && *req.Priority != "" {
		existingTicket.Priority = *req.Priority
	}

	// Explicitly set AssigneeID from the JWT-derived assigneeID if req.AssigneeID is provided
	if req.AssigneeID != nil {
		existingTicket.AssigneeID = sql.NullInt64{Int64: assigneeID, Valid: true}
	} else if !existingTicket.AssigneeID.Valid {
		// If it was unassigned and not provided in request, keep it unassigned
		existingTicket.AssigneeID = sql.NullInt64{Valid: false}
	}
	// If req.AssigneeID is nil and existingTicket.AssigneeID is valid, keep existing AssigneeID

	if err := models.UpdateTicket(h.db, r.Context(), existingTicket); err != nil {
		render.Status(r, http.StatusInternalServerError)
		renderer.PrettyJSON(w, r, err.Error())
		return
	}

	render.Status(r, http.StatusOK)
	renderer.PrettyJSON(w, r, existingTicket)
}

func (h *TicketHandler) CreateCustomerTicket(w http.ResponseWriter, r *http.Request) {

	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Priority    string `json:"priority"`
	}

	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Status(r, http.StatusBadRequest)
		renderer.PrettyJSON(w, r, err.Error())
		return
	}

	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		render.Status(r, http.StatusUnauthorized)
		renderer.PrettyJSON(w, r, "Unauthorized")
		return
	}

	requesterID, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		renderer.PrettyJSON(w, r, "Invalid user ID")
		return
	}

	ticket := models.Ticket{
		Title:       req.Title,
		Description: req.Description,
		RequesterID: requesterID,
		Status:      "Open",
		Priority:    req.Priority,
	}

	if err := models.CreateTicket(h.db, r.Context(), &ticket); err != nil {
		render.Status(r, http.StatusInternalServerError)
		renderer.PrettyJSON(w, r, err.Error())
		return
	}

	render.Status(r, http.StatusCreated)
	renderer.PrettyJSON(w, r, ticket)
}

func (h *TicketHandler) ListCustomerTickets(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		render.Status(r, http.StatusUnauthorized)
		renderer.PrettyJSON(w, r, "Unauthorized")
		return
	}

	requesterID, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		renderer.PrettyJSON(w, r, "Invalid user ID")
		return
	}

	tickets, err := models.ListTicketsByRequesterID(h.db, r.Context(), requesterID)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		renderer.PrettyJSON(w, r, err.Error())
		return
	}

	render.Status(r, http.StatusOK)
	renderer.PrettyJSON(w, r, tickets)
}

func (h *TicketHandler) GetCustomerTicket(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		render.Status(r, http.StatusUnauthorized)
		renderer.PrettyJSON(w, r, "Unauthorized")
		return
	}

	requesterID, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		renderer.PrettyJSON(w, r, "Invalid user ID")
		return
	}

	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		renderer.PrettyJSON(w, r, "Invalid ticket ID")
		return
	}

	ticket, err := models.GetTicketByID(h.db, r.Context(), id)
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

	if ticket.RequesterID != requesterID {
		render.Status(r, http.StatusForbidden)
		renderer.PrettyJSON(w, r, "You are not authorized to view this ticket")
		return
	}

	// Filter internal comments for customers
	filteredComments := []models.Comment{}
	for _, comment := range ticket.Comments {
		if !comment.IsInternal {
			filteredComments = append(filteredComments, comment)
		}
	}
	ticket.Comments = filteredComments

	render.Status(r, http.StatusOK)
	renderer.PrettyJSON(w, r, ticket)
}

func (h *TicketHandler) CloseCustomerTicket(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		render.Status(r, http.StatusUnauthorized)
		renderer.PrettyJSON(w, r, "Unauthorized")
		return
	}

	requesterID, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		renderer.PrettyJSON(w, r, "Invalid user ID")
		return
	}

	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		renderer.PrettyJSON(w, r, "Invalid ticket ID")
		return
	}

	existingTicket, err := models.GetTicketByID(h.db, r.Context(), id)
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

	if existingTicket.RequesterID != requesterID {
		render.Status(r, http.StatusForbidden)
		renderer.PrettyJSON(w, r, "You are not authorized to close this ticket")
		return
	}

	existingTicket.Status = "Closed"

	if err := models.UpdateTicket(h.db, r.Context(), existingTicket); err != nil {
		render.Status(r, http.StatusInternalServerError)
		renderer.PrettyJSON(w, r, err.Error())
		return
	}

	render.Status(r, http.StatusOK)
	renderer.PrettyJSON(w, r, existingTicket)
}
