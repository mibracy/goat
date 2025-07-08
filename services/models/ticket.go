package models

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/uptrace/bun"
)

// Ticket represents the Ticket model in the database.
type Ticket struct {
	bun.BaseModel `bun:"table:tickets,alias:ticket"`
	ID            int64         `bun:"id,pk,autoincrement,type:integer"`
	Title         string        `bun:"title,notnull"`
	Description   string        `bun:"description"`
	Status        string        `bun:"status,notnull,default:'Open'"`
	Priority      string        `bun:"priority,notnull,default:'Medium'"`
	RequesterID   int64         `bun:"requester_id,notnull"`
	AssigneeID    sql.NullInt64 `bun:"assignee_id"` // Use sql.NullInt64 for nullable foreign key
	CreatedAt     time.Time     `bun:"created_at,notnull,default:current_timestamp" json:"CreatedAt"`
	UpdatedAt     time.Time     `bun:"updated_at,notnull,default:current_timestamp" json:"UpdatedAt"`
	ClosedAt      sql.NullTime  `bun:"closed_at" json:"ClosedAt"`   // Use sql.NullTime for nullable timestamp
	Comments      []Comment     `bun:"-" json:"Comments,omitempty"` // This field is not stored in the database
}

// GetTicketByID retrieves a ticket from the database by its ID and also fetches related comments.
func GetTicketByID(db *bun.DB, ctx context.Context, ticketID int64) (*Ticket, error) {
	ticket := new(Ticket)
	err := db.NewSelect().Model(ticket).Where("id = ?", ticketID).Scan(ctx)
	if err != nil {
		return nil, err
	}

	comments, err := ListCommentsByTicketID(db, ctx, ticketID)
	if err != nil {
		// Log the error but don't fail the ticket retrieval if comments can't be fetched
		fmt.Printf("Error fetching comments for ticket %d: %v\n", ticketID, err)
	}
	ticket.Comments = comments

	return ticket, nil
}

// ListTickets retrieves all tickets from the database.
func ListTickets(db *bun.DB, ctx context.Context) ([]Ticket, error) {
	var tickets []Ticket
	err := db.NewSelect().Model(&tickets).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return tickets, nil
}

// CreateTicket inserts a new ticket into the database.
func CreateTicket(db *bun.DB, ctx context.Context, ticket *Ticket) error {
	_, err := db.NewInsert().Model(ticket).Exec(ctx)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			return fmt.Errorf("duplicate entry for ticket: %w", err)
		}
		return err
	}
	return nil
}

// UpdateTicket updates an existing ticket in the database.
func UpdateTicket(db *bun.DB, ctx context.Context, ticket *Ticket) error {

	// Preserve the original CreatedAt time.
	existingTicket, err := GetTicketByID(db, ctx, ticket.ID)
	if err != nil {
		return err
	}
	ticket.CreatedAt = existingTicket.CreatedAt

	ticket.UpdatedAt = time.Now()
	if ticket.Status == "Closed" {
		ticket.ClosedAt = sql.NullTime{Time: time.Now(), Valid: true}
	} else {
		ticket.ClosedAt = sql.NullTime{Valid: false}
	}

	_, err = db.NewUpdate().
		Model(ticket).
		Column("title", "description", "status", "priority", "requester_id", "assignee_id", "updated_at", "closed_at").
		Where("id = ?", ticket.ID).
		Exec(ctx)
	return err
}

// DeleteTicket deletes a ticket from the database by its ID.
func DeleteTicket(db *bun.DB, ctx context.Context, ticketID int64) error {
	_, err := db.NewDelete().Model(&Ticket{}).Where("id = ?", ticketID).Exec(ctx)
	return err
}
