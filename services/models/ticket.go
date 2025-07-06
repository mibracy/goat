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
	CreatedAt     time.Time     `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt     time.Time     `bun:"updated_at,notnull,default:current_timestamp"`
	ClosedAt      sql.NullTime  `bun:"closed_at"` // Use sql.NullTime for nullable timestamp
}

// GetTicketByID retrieves a ticket from the database by its ID.
func GetTicketByID(db *bun.DB, ctx context.Context, ticketID int64) (*Ticket, error) {
	ticket := new(Ticket)
	err := db.NewSelect().Model(ticket).Where("id = ?", ticketID).Scan(ctx)
	if err != nil {
		return nil, err
	}
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
	_, err := db.NewUpdate().Model(ticket).Where("id = ? ", ticket.ID).Exec(ctx)
	return err
}

// DeleteTicket deletes a ticket from the database by its ID.
func DeleteTicket(db *bun.DB, ctx context.Context, ticketID int64) error {
	_, err := db.NewDelete().Model(&Ticket{}).Where("id = ?", ticketID).Exec(ctx)
	return err
}
