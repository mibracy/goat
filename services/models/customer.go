package models

import (
	"context"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/uptrace/bun"
)

// Customer represents the Customer model in the database.
// It embeds bun.BaseModel for ORM functionalities.
type Customer struct {
	bun.BaseModel `bun:"table:customers,alias:customer"`
	ID            int64     `bun:"id,pk,autoincrement,type:integer"` // Unique identifier for the customer, primary key, auto-incrementing.
	Name          string    `bun:"name,notnull"`                     // Name of the customer, cannot be null.
	Email         string    `bun:"email,unique"`                     // Email address of the customer.
	Created       time.Time `bun:"created_at,notnull,default:current_timestamp"`
}

// GetCustomerById retrieves a customer from the database by their ID.
func GetCustomerById(db *bun.DB, ctx context.Context, id int64) (*Customer, error) {
	customer := new(Customer)
	err := db.NewSelect().Model(customer).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return customer, nil
}

// CreateCustomer inserts a new customer into the database.
func CreateCustomer(db *bun.DB, ctx context.Context, customer *Customer) error {
	_, err := db.NewInsert().Model(customer).Exec(ctx)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			return fmt.Errorf("duplicate entry for customer: %w", err)
		}
		return err
	}
	return nil
}

// UpdateCustomer updates an existing customer in the database.
func UpdateCustomer(db *bun.DB, ctx context.Context, customer *Customer) error {
	_, err := db.NewUpdate().Model(customer).WherePK().Exec(ctx)
	return err
}

// DeleteCustomer deletes a customer from the database by their ID.
func DeleteCustomer(db *bun.DB, ctx context.Context, id int64) error {
	_, err := db.NewDelete().Model(&Customer{}).Where("id = ?", id).Exec(ctx)
	return err
}

// ListCustomers retrieves a list of customers from the database.
func ListCustomers(db *bun.DB, ctx context.Context) ([]*Customer, error) {
	var customers []*Customer
	err := db.NewSelect().Model(&customers).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return customers, nil
}
