package models

import (
	"context"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/uptrace/bun"
)

// User represents the User model in the database.
type User struct {
	bun.BaseModel `bun:"table:users,alias:user"`
	ID            int64     `bun:"id,pk,autoincrement,type:integer"`
	Name          string    `bun:"name,notnull"`
	Email         string    `bun:"email,notnull,unique"`
	PasswordHash  string    `bun:"password_hash,notnull"`
	Role          string    `bun:"role,notnull,default:'Agent'"`
	CreatedAt     time.Time `bun:"created_at,notnull,default:current_timestamp"`
}

// GetUserByID retrieves a user from the database by their ID.
func GetUserByID(db *bun.DB, ctx context.Context, userID int64) (*User, error) {
	user := new(User)
	err := db.NewSelect().Model(user).Where("id = ?", userID).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// CreateUser inserts a new user into the database.
func CreateUser(db *bun.DB, ctx context.Context, user *User) error {
	_, err := db.NewInsert().Model(user).Exec(ctx)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			return fmt.Errorf("duplicate entry for user: %w", err)
		}
		return err
	}
	return nil
}

// UpdateUser updates an existing user in the database.
func UpdateUser(db *bun.DB, ctx context.Context, user *User) error {
	_, err := db.NewUpdate().Model(user).Where("id = ?", user.ID).Exec(ctx)
	return err
}

// DeleteUser deletes a user from the database by their ID.
func DeleteUser(db *bun.DB, ctx context.Context, userID int64) error {
	_, err := db.NewDelete().Model(&User{}).Where("id = ?", userID).Exec(ctx)
	return err
}

// GetUsers retrieves a list of users from the database.
func GetUsers(db *bun.DB, ctx context.Context) ([]*User, error) {
	var users []*User
	err := db.NewSelect().Model(&users).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// GetCustomers retrieves a list of users with the role 'Customer' from the database.
func GetUsersByRole(db *bun.DB, ctx context.Context, role string) ([]*User, error) {
	var users []*User
	err := db.NewSelect().Model(&users).Where("role = ?", role).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}
