package models

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/uptrace/bun"

	"goat/app/renderer"
	userModel "goat/services/models"
)

type UserHandler struct {
	db *bun.DB
}

func NewUserHandler(db *bun.DB) *UserHandler {
	return &UserHandler{db: db}
}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := userModel.GetUsers(h.db, context.Background())
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		renderer.PrettyJSON(w, r, err.Error())
		return
	}
	render.Status(r, http.StatusOK)
	renderer.PrettyJSON(w, r, users)
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	idNum, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		renderer.PrettyJSON(w, r, "Invalid user ID")
		return
	}
	user, err := userModel.GetUserByID(h.db, context.Background(), idNum)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		renderer.PrettyJSON(w, r, err.Error())
		return
	}
	if user == nil {
		render.Status(r, http.StatusNotFound)
		renderer.PrettyJSON(w, r, "User not found")
		return
	}
	render.Status(r, http.StatusOK)
	renderer.PrettyJSON(w, r, user)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	data := &userModel.User{}

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		renderer.PrettyJSON(w, r, err.Error())
		return
	}

	err = userModel.CreateUser(h.db, context.Background(), data)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		renderer.PrettyJSON(w, r, err.Error())
		return
	}
	render.Status(r, http.StatusCreated)
	renderer.PrettyJSON(w, r, data)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	idNum, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		renderer.PrettyJSON(w, r, "Invalid user ID")
		return
	}
	user := &userModel.User{}
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		renderer.PrettyJSON(w, r, err.Error())
		return
	}
	user.ID = idNum

	err = userModel.UpdateUser(h.db, context.Background(), user)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		renderer.PrettyJSON(w, r, err.Error())
		return
	}
	render.Status(r, http.StatusOK)
	renderer.PrettyJSON(w, r, user)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	idNum, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		renderer.PrettyJSON(w, r, "Invalid user ID")
		return
	}

	err = userModel.DeleteUser(h.db, context.Background(), idNum)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		renderer.PrettyJSON(w, r, err.Error())
		return
	}
	render.Status(r, http.StatusAccepted)
	renderer.PrettyJSON(w, r, map[string]string{"message": "User deleted successfully"})
}
