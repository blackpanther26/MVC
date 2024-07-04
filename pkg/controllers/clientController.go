package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
	"github.com/blackpanther26/mvc/pkg/utils"
	"github.com/blackpanther26/mvc/pkg/config"
	"github.com/blackpanther26/mvc/pkg/models"
	"github.com/gorilla/mux"
)

func ListBooks(w http.ResponseWriter, r *http.Request) {
	var books []models.Book
	result := config.DB.Find(&books)
	if result.Error != nil {
		http.Error(w, "Failed to fetch books", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Books": books,
	}

	utils.RenderTemplate(w, "clientPortal", data)
}

func CheckoutBook(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    bookID, err := strconv.Atoi(vars["id"])
    if err != nil {
        utils.RenderTemplateWithMessage(w, "clientPortal", "Invalid book ID", "error")
        return
    }

    var book models.Book
    result := config.DB.First(&book, bookID)
    if result.Error != nil {
        utils.RenderTemplateWithMessage(w, "clientPortal", "Book not found", "error")
        return
    }

    if book.TotalCopies <= 0 {
        utils.RenderTemplateWithMessage(w, "clientPortal", "No copies available for checkout", "error")
        return
    }

    var transaction models.Transaction
    result = config.DB.Where("user_id =? AND book_id =? AND status =?", getUserIDFromContext(r.Context()), uint(bookID), "pending").First(&transaction)
    if result.Error == nil && transaction.UserID == getUserIDFromContext(r.Context()) {
        utils.RenderTemplateWithMessage(w, "clientPortal", "You already have a pending checkout request for this book.", "error")
        return
    }

    tx := config.DB.Begin()

    transaction = models.Transaction{
        UserID:          getUserIDFromContext(r.Context()),
        BookID:          uint(bookID),
        TransactionType: "checkout",
        TransactionDate: time.Now(),
        DueDate:         calculateDueDate(time.Now()),
        Status:          "pending",
    }

    if err := tx.Save(&book).Error; err != nil {
        tx.Rollback()
        utils.RenderTemplateWithMessage(w, "clientPortal", "Failed to update book availability", "error")
        return
    }

    if err := tx.Create(&transaction).Error; err != nil {
        tx.Rollback()
        utils.RenderTemplateWithMessage(w, "clientPortal", "Failed to log transaction", "error")
        return
    }

    tx.Commit()

    utils.RenderTemplateWithMessage(w, "clientPortal", "Book checkout request sent successfully.", "success")
}

func CheckinBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookID, err := strconv.Atoi(vars["id"])
	if err!= nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	var book models.Book
	result := config.DB.First(&book, bookID)
	if result.Error!= nil {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	// Simulate check-in process (update book availability, update transaction status, etc.)
	tx := config.DB.Begin()

	// Find the transaction for this book that is approved
	var transaction models.Transaction
	result = config.DB.Where("book_id =? AND status = 'approved'", bookID).First(&transaction)
	if result.Error!= nil {
		tx.Rollback()
		http.Error(w, "Transaction not found or already checked in", http.StatusBadRequest)
		return
	}

	transaction.TransactionType = "checkedin"
	transaction.Status = "pending"
	transaction.ReturnDate = func() *time.Time { t := time.Now(); return &t }()

	if err := tx.Save(&book).Error; err!= nil {
		tx.Rollback()
		http.Error(w, "Failed to update book availability", http.StatusInternalServerError)
		return
	}

	if err := tx.Save(&transaction).Error; err!= nil {
		tx.Rollback()
		http.Error(w, "Failed to update transaction status", http.StatusInternalServerError)
		return
	}

	tx.Commit()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Book checked in successfully"}`))
}


func UserHistory(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r.Context()) 

	var transactions []models.Transaction
	result := config.DB.Where("user_id = ?", userID).Find(&transactions)
	if result.Error != nil {
		http.Error(w, "Failed to fetch user history", http.StatusInternalServerError)
		return
	}

	jsonTransactions, err := json.Marshal(transactions)
	if err != nil {
		http.Error(w, "Failed to marshal user history", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonTransactions)
}

func getUserIDFromContext(ctx context.Context) uint {
	user := ctx.Value("user").(models.User)
	return user.ID
}

func calculateDueDate(startDate time.Time) *time.Time {
	dueDate := startDate.AddDate(0, 0, 14)
	return &dueDate
}
