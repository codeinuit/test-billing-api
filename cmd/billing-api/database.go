package main

import (
	"database/sql"
	"fmt"

	logger "github.com/codeinuit/test-billing-api/pkg/log"
)

type Database struct {
	log  logger.Logger
	conn *sql.DB
}

type User struct {
	UserID    uint    `json:"id"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Balance   float64 `json:"balance"`
}

type Invoice struct {
	ID     uint          `json:"id"`
	UserID uint          `json:"user_id"`
	Amount float64       `json:"amount"`
	Label  string        `json:"label"`
	Status InvoiceStatus `json:"status"`
}

type InvoiceStatus string

const (
	Pending InvoiceStatus = "pending"
	Paid    InvoiceStatus = "paid"
)

//
// Users
//

func (db Database) getUserByUserID(userID uint) (user User, err error) {
	query := fmt.Sprintf("SELECT id, first_name, last_name, balance FROM users WHERE id = %d;", userID)
	row := db.conn.QueryRow(query)

	return user, row.Scan(&user.UserID, &user.FirstName, &user.LastName, &user.Balance)
}

// getUsers returns the list of users from database
func (db Database) getUsers() (users []User, err error) {
	rows, err := db.conn.Query("SELECT id, first_name, last_name, balance FROM users;")
	if err != nil {
		db.log.Error("could not get users from database : ", err.Error())
		return users, err
	}
	defer rows.Close()

	for rows.Next() {
		var user User

		rows.Scan(&user.UserID, &user.FirstName, &user.LastName, &user.Balance)
		users = append(users, user)
	}

	return users, err
}

type Operation string

const (
	OperationDecrement Operation = "-"
	OperationIncrement Operation = "+"
)

func (db Database) updateUserBalanceByID(ID uint, amount float64, op Operation) (err error) {
	query := fmt.Sprintf("UPDATE users SET balance = balance %s $1 WHERE id = $2;", op)
	_, err = db.conn.Exec(query, amount, ID)
	return err
}

//
// Invoices
//

func (db Database) getInvoiceByID(ID uint) (invoice Invoice, err error) {
	query := fmt.Sprintf("SELECT id, user_id, label, amount, status FROM invoices WHERE id = %d;", ID)
	row := db.conn.QueryRow(query)

	err = row.Scan(&invoice.ID, &invoice.UserID, &invoice.Label, &invoice.Amount, &invoice.Status)

	return invoice, err
}

func (db Database) insertNewInvoice(i Invoice) (err error) {
	query := "INSERT INTO invoices (user_id, amount, label) VALUES ($1, $2, $3);"
	_, err = db.conn.Exec(query, i.UserID, i.Amount, i.Label)
	return err
}

func (db Database) updateInvoiceStatusByID(ID uint, status InvoiceStatus) (err error) {
	query := "UPDATE invoices SET status = $1 WHERE id = $2;"
	_, err = db.conn.Exec(query, status, ID)
	return err
}
