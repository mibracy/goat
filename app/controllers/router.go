package controllers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"goat/app/models"
	"goat/db"
)

func SetupServer() {
	fmt.Printf("Server starting on :8420\n")
	d := db.ConnectDB()
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
	customerHandler := models.NewCustomerHandler(d)
	ticketHandler := models.NewTicketHandler(d)
	commentHandler := models.NewCommentHandler(d)

	r.Route("/admin", func(r chi.Router) {

		r.Get("/", userHandler.ListUsers)
		r.Post("/", userHandler.CreateUser)
		r.Get("/{id}", userHandler.GetUsers)
		r.Put("/update/{id}", userHandler.UpdateUser)
		r.Delete("/delete/{id}", userHandler.DeleteUser)

		r.Get("/customers", customerHandler.ListCustomers)
		r.Get("/tickets", ticketHandler.ListTickets)
		r.Get("/comments", commentHandler.ListComments)

	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	http.ListenAndServe(":8420", r)
}
