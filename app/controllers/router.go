package controllers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"goat/app/models"
	"goat/services/config"
)

func SetupServer() {
	fmt.Printf("Server starting on :8420\n")
	d := config.ConnectDB()
	r := chi.NewRouter()
	r.Use(middleware.Logger)
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
		r.Get("/", userHandler.ListUsers)
		r.Post("/", userHandler.CreateUser)
		r.Get("/{id}", userHandler.GetUsers)
		r.Put("/update/{id}", userHandler.UpdateUser)
		r.Delete("/delete/{id}", userHandler.DeleteUser)
		r.Get("/role/{role}", userHandler.ListUsersByRole)
		r.Get("/tickets", ticketHandler.ListTickets)
		r.Post("/tickets", ticketHandler.CreateTicket)
		r.Get("/tickets/{}", ticketHandler.ListTickets)
		r.Get("/comments", commentHandler.ListComments)
	})

	r.Route("/agent", func(r chi.Router) {

	})

	r.Route("/customer", func(r chi.Router) {

	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	http.ListenAndServe(":8420", r)
}
