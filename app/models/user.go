package models

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/uptrace/bun"

	"goat/app/renderer"
	"goat/services/models"
)

type UserHandler struct {
	db *bun.DB
}

func NewUserHandler(db *bun.DB) *UserHandler {
	return &UserHandler{db: db}
}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := models.GetUsers(h.db, context.Background())
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		renderer.PrettyJSON(w, r, err.Error())
		return
	}
	render.Status(r, http.StatusOK)
	renderer.PrettyJSON(w, r, users)
}

func (h *UserHandler) ListUsersByRole(w http.ResponseWriter, r *http.Request) {
	role := chi.URLParam(r, "role")
	users, err := models.GetUsersByRole(h.db, context.Background(), role)
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
	user, err := models.GetUserByID(h.db, context.Background(), idNum)
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
	data := &models.User{}

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		renderer.PrettyJSON(w, r, err.Error())
		return
	}

	// Generate a random password hash
	data.PasswordHash = gofakeit.Password(true, true, true, true, true, 10)

	err = models.CreateUser(h.db, context.Background(), data)
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

	existingUser, err := models.GetUserByID(h.db, context.Background(), idNum)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		renderer.PrettyJSON(w, r, err.Error())
		return
	}
	if existingUser == nil {
		render.Status(r, http.StatusNotFound)
		renderer.PrettyJSON(w, r, "User not found")
		return
	}

	var updateData struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Role  string `json:"role"`
	}
	err = json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		renderer.PrettyJSON(w, r, err.Error())
		return
	}

	existingUser.Name = updateData.Name
	existingUser.Email = updateData.Email
	existingUser.Role = updateData.Role

	err = models.UpdateUser(h.db, context.Background(), existingUser)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		renderer.PrettyJSON(w, r, err.Error())
		return
	}
	render.Status(r, http.StatusOK)
	renderer.PrettyJSON(w, r, existingUser)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	idNum, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		renderer.PrettyJSON(w, r, "Invalid user ID")
		return
	}

	err = models.DeleteUser(h.db, context.Background(), idNum)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		renderer.PrettyJSON(w, r, err.Error())
		return
	}
	render.Status(r, http.StatusAccepted)
	renderer.PrettyJSON(w, r, map[string]string{"message": "User deleted successfully"})
}
