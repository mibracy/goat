package models

import (
	"context"
	"net/http"

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

	ctx := context.Background()

	comments, err := models.ListComments(h.db, ctx)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		renderer.PrettyJSON(w, r, err.Error())
		return
	}

	render.Status(r, http.StatusOK)
	renderer.PrettyJSON(w, r, comments)

}
