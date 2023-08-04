package main

import (
	"fmt"
	"net/http"

	logger "github.com/codeinuit/test-billing-api/pkg/log"

	"github.com/gin-gonic/gin"
)

type handlers struct {
	log logger.Logger
	db  *Database
}

type postInvoice struct {
	UserID uint    `json:"user_id"`
	Amount float64 `json:"amount"`
	Label  string  `json:"label"`
}

type postTransaction struct {
	InvoiceID uint    `json:"invoice_id"`
	Amount    float64 `json:"amount"`
	Reference string  `json:"reference"`
}

// healthcheck works as a ping and returns a OK status
func (h handlers) healthcheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}

// getUsers handles /users route and returns the list of users
func (h handlers) getUsers(c *gin.Context) {
	u, err := h.db.getUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not retrieve users from database"})
		return
	}

	c.IndentedJSON(http.StatusOK, u)
}

// postInvoice handles POST /invoice and inserts a new invoice in DB
func (h handlers) postInvoice(c *gin.Context) {
	var invoice postInvoice

	if err := c.ShouldBindJSON(&invoice); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing parameter"})
		return
	}

	_, err := h.db.getUserByUserID(invoice.UserID)
	if err != nil {
		h.log.Warn("could not insert invoice : ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not insert invoice"})
		return
	}

	err = h.db.insertNewInvoice(Invoice{UserID: invoice.UserID, Label: invoice.Label, Amount: invoice.Amount})
	if err != nil {
		h.log.Warn("could not insert invoice : ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not insert invoice"})
		return
	}
	c.Status(http.StatusNoContent)
}

// postTransaction handles POST /transaction
func (h handlers) postTransaction(c *gin.Context) {
	var transaction postTransaction

	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing parameter"})
		return
	}

	invoice, err := h.db.getInvoiceByID(transaction.InvoiceID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("no invoice found with id %d", transaction.InvoiceID)})
		return
	}

	if invoice.Status == Paid {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "invoice is already paid"})
		return
	}

	if invoice.Amount != transaction.Amount {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("transaction's amount does not match the invoice (was awaiting for %f)", invoice.Amount)})
		return
	}

	user, err := h.db.getUserByUserID(invoice.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no user found related to the invoice"})
		return
	}

	if user.Balance < transaction.Amount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "not enough balance"})
		return
	}

	err = h.db.updateUserBalanceByID(invoice.UserID, invoice.Amount, OperationDecrement)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update user balance"})
		return
	}
	err = h.db.updateInvoiceStatusByID(invoice.ID, Paid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update invoice status"})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"status": "OK"})
}
