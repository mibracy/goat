package models

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt/v5"
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
		if err == sql.ErrNoRows {
			render.Status(r, http.StatusNotFound)
			renderer.PrettyJSON(w, r, "User not found")
			return
		}
		render.Status(r, http.StatusInternalServerError)
		renderer.PrettyJSON(w, r, err.Error())
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

	if data.PasswordHash == "" {
		data.PasswordHash = "password"
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		renderer.PrettyJSON(w, r, "Failed to hash password")
		return
	}
	data.PasswordHash = string(hashedPassword)

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
		if err == sql.ErrNoRows {
			render.Status(r, http.StatusNotFound)
			renderer.PrettyJSON(w, r, "User not found")
			return
		}
		render.Status(r, http.StatusInternalServerError)
		renderer.PrettyJSON(w, r, err.Error())
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

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		renderer.PrettyJSON(w, r, "Invalid request body")
		return
	}

	user, err := models.GetUserByEmail(h.db, context.Background(), creds.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			render.Status(r, http.StatusUnauthorized)
			renderer.PrettyJSON(w, r, "Invalid credentials")
			return
		}
		render.Status(r, http.StatusInternalServerError)
		renderer.PrettyJSON(w, r, err.Error())
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(creds.Password))
	if err != nil {
		render.Status(r, http.StatusUnauthorized)
		renderer.PrettyJSON(w, r, "Invalid credentials")
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expirationTime),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Subject:   strconv.FormatInt(user.ID, 10),
		Issuer:    "goat",
		Audience:  []string{user.Role},
	}

	jwtKey := []byte(os.Getenv("JWT_SECRET"))
	if len(jwtKey) == 0 {
		jwtKey = []byte("default-secret-key")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		renderer.PrettyJSON(w, r, "Failed to create token")
		return
	}

	render.Status(r, http.StatusOK)
	renderer.PrettyJSON(w, r, map[string]string{"token": tokenString})
}

func (h *UserHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}

	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Status(r, http.StatusBadRequest)
		renderer.PrettyJSON(w, r, "Invalid request body")
		return
	}

	user, err := models.GetUserByEmail(h.db, context.Background(), req.Email)
	if err != nil {
		render.Status(r, http.StatusNotFound)
		renderer.PrettyJSON(w, r, "User with that email not found")
		return
	}

	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		render.Status(r, http.StatusInternalServerError)
		renderer.PrettyJSON(w, r, "Failed to generate token")
		return
	}
	token := hex.EncodeToString(tokenBytes)
	expires := time.Now().Add(time.Hour * 1).UTC()

	user.PasswordResetToken.String = token
	user.PasswordResetToken.Valid = true
	user.PasswordResetExpires.Time = expires
	user.PasswordResetExpires.Valid = true

	if err := models.UpdateUser(h.db, context.Background(), user); err != nil {
		render.Status(r, http.StatusInternalServerError)
		renderer.PrettyJSON(w, r, "Failed to save reset token")
		return
	}

	render.Status(r, http.StatusOK)
	renderer.PrettyJSON(w, r, map[string]string{"message": "Password reset link sent to your email (check console for token)", "token": token})
}

func (h *UserHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}

	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Status(r, http.StatusBadRequest)
		renderer.PrettyJSON(w, r, "Invalid request body")
		return
	}

	user := new(models.User)
	err := h.db.NewSelect().Model(user).Where("password_reset_token = ?", req.Token).Scan(context.Background())
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		renderer.PrettyJSON(w, r, "Invalid or expired token")
		return
	}

	if !user.PasswordResetExpires.Valid || user.PasswordResetExpires.Time.Before(time.Now().UTC()) {
		render.Status(r, http.StatusBadRequest)
		renderer.PrettyJSON(w, r, "Invalid or expired token")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		renderer.PrettyJSON(w, r, "Failed to hash new password")
		return
	}
	user.PasswordHash = string(hashedPassword)
	user.PasswordResetToken.Valid = false
	user.PasswordResetExpires.Valid = false

	if err := models.UpdateUser(h.db, context.Background(), user); err != nil {
		render.Status(r, http.StatusInternalServerError)
		renderer.PrettyJSON(w, r, err.Error())
		return
	}

	render.Status(r, http.StatusOK)
	renderer.PrettyJSON(w, r, map[string]string{"message": "Password has been reset successfully"})
}

func (h *UserHandler) RegisterCustomer(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Status(r, http.StatusBadRequest)
		renderer.PrettyJSON(w, r, "Invalid request body")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		renderer.PrettyJSON(w, r, "Failed to hash password")
		return
	}

	user := &models.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Role:         "Customer",
	}

	if err := models.CreateUser(h.db, context.Background(), user); err != nil {
		render.Status(r, http.StatusInternalServerError)
		renderer.PrettyJSON(w, r, err.Error())
		return
	}

	render.Status(r, http.StatusCreated)
	renderer.PrettyJSON(w, r, user)
}
