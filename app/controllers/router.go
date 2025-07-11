package controllers

import (
	"fmt"
	"goat/app/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"goat/app/models"
	"goat/services/config"
)

func SetupServer() {
	fmt.Printf("Server starting on :8420\n")
	d := config.ConnectDB()
	r := chi.NewRouter()
	r.Use(chiMiddleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	userHandler := models.NewUserHandler(d)
	ticketHandler := models.NewTicketHandler(d)
	commentHandler := models.NewCommentHandler(d)

	r.Route("/admin", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Use(middleware.RoleMiddleware("Admin"))
		r.Get("/", userHandler.ListUsers)
		r.Post("/", userHandler.CreateUser)
		r.Get("/{id}", userHandler.GetUsers)
		r.Put("/update/{id}", userHandler.UpdateUser)
		r.Delete("/delete/{id}", userHandler.DeleteUser)
		r.Get("/role/{role}", userHandler.ListUsersByRole)
		r.Get("/tickets", ticketHandler.ListTickets)
		r.Post("/tickets", ticketHandler.CreateTicket)
		r.Get("/tickets/{id}", ticketHandler.GetTicket)
		r.Put("/tickets/{id}", ticketHandler.UpdateTicket)
		r.Get("/comments", commentHandler.ListComments)
		r.Post("/comments", commentHandler.CreateComment)
		r.Get("/comments/ticket/{id}", commentHandler.ListCommentsByTicketID)
	})

	r.Post("/login", userHandler.Login)
	r.Post("/forgot-password", userHandler.ForgotPassword)
	r.Post("/reset-password", userHandler.ResetPassword)
	r.Post("/register", userHandler.RegisterCustomer)

	r.Route("/agent", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Use(middleware.RoleMiddleware("Admin", "Agent"))
		r.Get("/tickets", ticketHandler.ListAgentTickets)
		r.Get("/tickets/{id}", ticketHandler.GetAgentTicket)
		r.Put("/tickets/{id}", ticketHandler.UpdateAgentTicket)
		r.Post("/tickets/{id}/comments", commentHandler.CreateAgentComment)
	})

	r.Route("/customer", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Use(middleware.RoleMiddleware("Admin", "Agent", "Customer"))
		r.Post("/tickets", ticketHandler.CreateCustomerTicket)
		r.Get("/tickets", ticketHandler.ListCustomerTickets)
		r.Get("/tickets/{id}", ticketHandler.GetCustomerTicket)
		r.Post("/tickets/{id}/comments", commentHandler.CreateCustomerComment)
		r.Put("/tickets/{id}", ticketHandler.CloseCustomerTicket)
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	http.ListenAndServe(":8420", r)
}
