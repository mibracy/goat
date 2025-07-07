package models

import (
	"context"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/uptrace/bun"
)

// Comment represents the Comment model in the database.
type Comment struct {
	bun.BaseModel `bun:"table:comments,alias:comment"`
	ID            int64     `bun:"id,pk,autoincrement,type:integer"`
	TicketID      int64     `bun:"ticket_id,notnull"`
	AuthorID      int64     `bun:"author_id,notnull"`
	Body          string    `bun:"body,notnull"`
	IsInternal    bool      `bun:"is_internal,default:false"`
	CreatedAt     time.Time `bun:"created_at,notnull,default:current_timestamp"`
}

// GetCommentByID retrieves a comment from the database by its ID.
func GetCommentByID(db *bun.DB, ctx context.Context, commentID int64) (*Comment, error) {
	comment := new(Comment)
	err := db.NewSelect().Model(comment).Where("id = ?", commentID).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return comment, nil
}

// ListComments retrieves all comments from the database.
func ListComments(db *bun.DB, ctx context.Context) ([]Comment, error) {
	var comments []Comment
	err := db.NewSelect().Model(&comments).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return comments, nil
}

// ListCommentsByTicketID retrieves comments for a specific ticket from the database.
func ListCommentsByTicketID(db *bun.DB, ctx context.Context, ticketID int64) ([]Comment, error) {
	var comments []Comment
	err := db.NewSelect().Model(&comments).Where("ticket_id = ?", ticketID).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return comments, nil
}

// CreateComment inserts a new comment into the database.
func CreateComment(db *bun.DB, ctx context.Context, comment *Comment) error {
	_, err := db.NewInsert().Model(comment).Exec(ctx)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			return fmt.Errorf("duplicate entry for comment: %w", err)
		}
		return err
	}
	return nil
}

// UpdateComment updates an existing comment in the database.
func UpdateComment(db *bun.DB, ctx context.Context, comment *Comment) error {
	_, err := db.NewUpdate().Model(comment).WherePK().Exec(ctx)
	return err
}

// DeleteComment deletes a comment from the database by its ID.
func DeleteComment(db *bun.DB, ctx context.Context, commentID int64) error {
	_, err := db.NewDelete().Model(&Comment{}).Where("id = ?", commentID).Exec(ctx)
	return err
}
