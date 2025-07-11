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

type CommentHandler struct {
	db *bun.DB
}

func NewCommentHandler(db *bun.DB) *CommentHandler {
	return &CommentHandler{db: db}
}

// ListComments handles the request to list all comments.
func (h *CommentHandler) ListComments(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	userRole, ok := ctx.Value(middleware.UserRoleKey).(string)
	if !ok {
		render.Status(r, http.StatusUnauthorized)
		renderer.PrettyJSON(w, r, "Unauthorized")
		return
	}

	comments, err := models.ListComments(h.db, ctx)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		renderer.PrettyJSON(w, r, err.Error())
		return
	}

	// Filter comments based on user role
	if userRole == "Customer" {
		filteredComments := []models.Comment{}
		for _, comment := range comments {
			if !comment.IsInternal {
				filteredComments = append(filteredComments, comment)
			}
		}
		render.Status(r, http.StatusOK)
		renderer.PrettyJSON(w, r, filteredComments)
		return
	}

	render.Status(r, http.StatusOK)
	renderer.PrettyJSON(w, r, comments)

}

// CreateComment handles the request to create a new comment.
func (h *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	var req struct {
		TicketID   int64  `json:"TicketID"`
		AuthorID   int64  `json:"AuthorID"`
		Body       string `json:"Body"`
		IsInternal bool   `json:"IsInternal"`
	}

	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Status(r, http.StatusBadRequest)
		renderer.PrettyJSON(w, r, err.Error())
		return
	}

	comment := models.Comment{
		TicketID:   req.TicketID,
		Body:       req.Body,
		IsInternal: req.IsInternal,
	}

	// Check if the author exists
	author, err := models.GetUserByID(h.db, ctx, req.AuthorID)
	if err != nil {
		if err == sql.ErrNoRows {
			render.Status(r, http.StatusNotFound)
			renderer.PrettyJSON(w, r, "Author not found")
			return
		}
		render.Status(r, http.StatusInternalServerError)
		renderer.PrettyJSON(w, r, err.Error())
		return
	}
	comment.AuthorID = author.ID

	// Check if the ticket exists
	_, err = models.GetTicketByID(h.db, ctx, req.TicketID)
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

	if err := models.CreateComment(h.db, ctx, &comment); err != nil {
		render.Status(r, http.StatusInternalServerError)
		renderer.PrettyJSON(w, r, err.Error())
		return
	}

	render.Status(r, http.StatusCreated)
	renderer.PrettyJSON(w, r, comment)
}

// ListCommentsByTicketID handles the request to list comments for a specific ticket.
func (h *CommentHandler) ListCommentsByTicketID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		renderer.PrettyJSON(w, r, "Invalid ticket ID")
		return
	}

	userRole, ok := ctx.Value(middleware.UserRoleKey).(string)
	if !ok {
		render.Status(r, http.StatusUnauthorized)
		renderer.PrettyJSON(w, r, "Unauthorized")
		return
	}

	comments, err := models.ListCommentsByTicketID(h.db, ctx, id)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		renderer.PrettyJSON(w, r, err.Error())
		return
	}

	// Filter comments based on user role
	if userRole == "Customer" {
		filteredComments := []models.Comment{}
		for _, comment := range comments {
			if !comment.IsInternal {
				filteredComments = append(filteredComments, comment)
			}
		}
		render.Status(r, http.StatusOK)
		renderer.PrettyJSON(w, r, filteredComments)
		return
	}

	render.Status(r, http.StatusOK)
	renderer.PrettyJSON(w, r, comments)
}

func (h *CommentHandler) CreateAgentComment(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		render.Status(r, http.StatusUnauthorized)
		renderer.PrettyJSON(w, r, "Unauthorized")
		return
	}

	authorID, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		renderer.PrettyJSON(w, r, "Invalid user ID")
		return
	}

	ticketIDParam := chi.URLParam(r, "id")
	ticketID, err := strconv.ParseInt(ticketIDParam, 10, 64)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		renderer.PrettyJSON(w, r, "Invalid ticket ID")
		return
	}

	// Check if the ticket exists
	_, err = models.GetTicketByID(h.db, r.Context(), ticketID)
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
		Body       string `json:"body"`
		IsInternal bool   `json:"is_internal"`
	}

	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Status(r, http.StatusBadRequest)
		renderer.PrettyJSON(w, r, err.Error())
		return
	}

	comment := models.Comment{
		TicketID:   ticketID,
		AuthorID:   authorID,
		Body:       req.Body,
		IsInternal: req.IsInternal,
	}

	if err := models.CreateComment(h.db, r.Context(), &comment); err != nil {
		render.Status(r, http.StatusInternalServerError)
		renderer.PrettyJSON(w, r, err.Error())
		return
	}

	render.Status(r, http.StatusCreated)
	renderer.PrettyJSON(w, r, comment)
}

func (h *CommentHandler) CreateCustomerComment(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		render.Status(r, http.StatusUnauthorized)
		renderer.PrettyJSON(w, r, "Unauthorized")
		return
	}

	authorID, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		renderer.PrettyJSON(w, r, "Invalid user ID")
		return
	}

	ticketIDParam := chi.URLParam(r, "id")
	ticketID, err := strconv.ParseInt(ticketIDParam, 10, 64)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		renderer.PrettyJSON(w, r, "Invalid ticket ID")
		return
	}

	ticket, err := models.GetTicketByID(h.db, r.Context(), ticketID)
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

	if ticket.RequesterID != authorID {
		render.Status(r, http.StatusForbidden)
		renderer.PrettyJSON(w, r, "You are not authorized to comment on this ticket")
		return
	}

	var req struct {
		Body string `json:"body"`
	}

	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Status(r, http.StatusBadRequest)
		renderer.PrettyJSON(w, r, err.Error())
		return
	}

	comment := models.Comment{
		TicketID:   ticketID,
		AuthorID:   authorID,
		Body:       req.Body,
		IsInternal: false,
	}

	if err := models.CreateComment(h.db, r.Context(), &comment); err != nil {
		render.Status(r, http.StatusInternalServerError)
		renderer.PrettyJSON(w, r, err.Error())
		return
	}

	render.Status(r, http.StatusCreated)
	renderer.PrettyJSON(w, r, comment)
}
