package main

import (
	"context"
	"fmt"
	"log"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/uptrace/bun"

	router "goat/app/controllers"
	"goat/services/config"
	model "goat/services/models"
)

func demoUsersOps(db *bun.DB, ctx context.Context) {
	// Retrieve a user by ID
	fmt.Println("\n--- Retrieving User ---")
	testUser, err := model.GetUserByID(db, ctx, 1)

	if err != nil {
		log.Printf("Error retrieving user: %v\n", err)
	} else {
		fmt.Printf("Retrieved User:\n")
		fmt.Printf("            id: %d\n "+
			"         name: %s\n "+
			"        email: %s\n "+
			"         hash: %s\n "+
			"         role: %s\n "+
			"      created: %s \n\n", testUser.ID,
			testUser.Name,
			testUser.Email,
			testUser.PasswordHash,
			testUser.Role,
			testUser.CreatedAt)
	}

	// Insert a new user
	fmt.Println("\n--- Creating New User ---")
	newUser := &model.User{
		Name:  gofakeit.Name(),
		Email: gofakeit.Email(),
		Role:  gofakeit.RandomString([]string{"Admin", "Agent", "Customer"}),
	}

	err = model.CreateUser(db, ctx, newUser)
	if err != nil {
		log.Printf("Error creating user: %v\n", err)
	} else {
		fmt.Printf("Created User:\n")
		fmt.Printf("          id: %d\n "+
			"       name: %s\n "+
			"      email: %s\n "+
			"       hash: %s\n "+
			"       role: %s\n "+
			"    created: %s \n\n", newUser.ID,
			newUser.Name,
			newUser.Email,
			newUser.PasswordHash,
			newUser.Role,
			newUser.CreatedAt)
	}
}

func demoTicketOps(db *bun.DB, ctx context.Context) {
	// Retrieve a ticket by ID
	fmt.Println("\n--- Retrieving Ticket ---")
	testTicket, err := model.GetTicketByID(db, ctx, 1)

	if err != nil {
		log.Printf("Error retrieving ticket: %v\n", err)
	} else {
		fmt.Printf("Retrieved Ticket:\n")
		fmt.Printf("              id: %d\n "+
			"          title: %s\n "+
			"    description: %s\n "+
			"         status: %s\n "+
			"       priority: %s\n "+
			"   requester_id: %d\n "+
			"    assignee_id: %v\n "+
			"        created: %s\n "+
			"        updated: %s\n "+
			"         closed: %v \n\n",
			testTicket.ID,
			testTicket.Title,
			testTicket.Description,
			testTicket.Status,
			testTicket.Priority,
			testTicket.RequesterID,
			testTicket.AssigneeID.Int64,
			testTicket.CreatedAt,
			testTicket.UpdatedAt,
			testTicket.ClosedAt.Time)
	}

	// Insert a new ticket
	fmt.Println("\n--- Creating New Ticket ---")
	newTicket := &model.Ticket{
		Title:       gofakeit.Sentence(6),
		Description: gofakeit.Paragraph(2, 5, 10, ""),
		Status:      gofakeit.RandomString([]string{"Open", "Closed", "Pending"}),
		Priority:    gofakeit.RandomString([]string{"Low", "Medium", "High", "Urgent"}),
		RequesterID: int64(gofakeit.Number(1, 10)),
	}

	err = model.CreateTicket(db, ctx, newTicket)
	if err != nil {
		log.Printf("Error creating ticket: %v\n", err)
	} else {
		fmt.Printf("Created Ticket:\n")
		fmt.Printf("            id: %d\n "+
			"        title: %s\n "+
			"  description: %s\n "+
			"       status: %s\n "+
			"     priority: %s\n "+
			" requester_id: %d\n "+
			"  assignee_id: %v\n "+
			"      created: %s\n "+
			"      updated: %s\n "+
			"       closed: %v \n\n",
			newTicket.ID,
			newTicket.Title,
			newTicket.Description,
			newTicket.Status,
			newTicket.Priority,
			newTicket.RequesterID,
			newTicket.AssigneeID.Int64,
			newTicket.CreatedAt,
			newTicket.UpdatedAt,
			newTicket.ClosedAt.Time)
	}
}

func demoCommentOps(db *bun.DB, ctx context.Context) {
	// Retrieve a comment by ID
	fmt.Println("\n--- Retrieving Comment ---")
	testComment, err := model.GetCommentByID(db, ctx, 1)

	if err != nil {
		log.Printf("Error retrieving comment: %v\n", err)
	} else {
		fmt.Printf("Retrieved Comment:\n")
		fmt.Printf("               id: %d\n "+
			"       ticket_id: %d\n "+
			"       author_id: %d\n "+
			"            body: %s\n "+
			"     is_internal: %t\n "+
			"         created: %s \n\n", testComment.ID,
			testComment.TicketID,
			testComment.AuthorID,
			testComment.Body,
			testComment.IsInternal,
			testComment.CreatedAt)
	}

	// Insert a new comment
	fmt.Println("\n--- Creating New Comment ---")
	newComment := &model.Comment{
		TicketID:   int64(gofakeit.Number(1, 10)),
		AuthorID:   int64(gofakeit.Number(1, 10)),
		Body:       gofakeit.Sentence(10),
		IsInternal: gofakeit.Bool(),
	}

	err = model.CreateComment(db, ctx, newComment)
	if err != nil {
		log.Printf("Error creating comment: %v\n", err)
	} else {
		fmt.Printf("Created Comment:\n")
		fmt.Printf("             id: %d\n "+
			"     ticket_id: %d\n "+
			"     author_id: %d\n "+
			"          body: %s\n "+
			"   is_internal: %t\n "+
			"       created: %s \n\n", newComment.ID,
			newComment.TicketID,
			newComment.AuthorID,
			newComment.Body,
			newComment.IsInternal,
			newComment.CreatedAt)
	}
}

func main() {
	ctx := context.Background()
	db := config.ConnectDB()

	demoUsersOps(db, ctx)
	demoTicketOps(db, ctx)
	demoCommentOps(db, ctx)

	// Set up and start the HTTP server.
	router.SetupServer()
}
