package models

import (
	"errors"
	"time"

	"github.com/blackpanther26/mvc/pkg/config"
	"github.com/blackpanther26/mvc/pkg/types"
	"gorm.io/gorm"
)

func GetAllBooks() ([]types.Book, error) {
	var books []types.Book
	result := config.DB.Find(&books)
	if result.Error != nil {
		return nil, result.Error
	}
	return books, nil
}

func GetBookByID(id int) (*types.Book, error) {
	var book types.Book
	result := config.DB.First(&book, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &book, nil
}

func CheckoutBook(userID uint, bookID int) error {
	var book types.Book
	tx := config.DB.Begin()

	result := tx.First(&book, bookID)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	if book.TotalCopies <= 0 {
		tx.Rollback()
		return errors.New("no copies available for checkout")
	}

	var transaction types.Transaction
	result = tx.Where("user_id = ? AND book_id = ? AND status = ?", userID, bookID, "pending").First(&transaction)
	if result.Error == nil && transaction.UserID == userID {
		tx.Rollback()
		return errors.New("already have a pending checkout request for this book")
	}

	transaction = types.Transaction{
		UserID:          userID,
		BookID:          uint(bookID),
		TransactionType: "checkout",
		TransactionDate: time.Now(),
		DueDate:         calculateDueDate(time.Now()),
		Status:          "pending",
	}

	if err := tx.Save(&book).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func CheckinBook(bookID int) error {
	var book types.Book
	tx := config.DB.Begin()

	result := tx.First(&book, bookID)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	var transaction types.Transaction
	result = tx.Where("book_id = ? AND status = 'approved'", bookID).First(&transaction)
	if result.Error != nil {
		tx.Rollback()
		return errors.New("book not checked out earlier or record not found")
	}

	
	transaction.TransactionType = "checkin"
	transaction.Status = "pending"
	transaction.ReturnDate = func() *time.Time { t := time.Now(); return &t }()

	if err := tx.Save(&book).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Save(&transaction).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func GetUserTransactions(userID uint) ([]types.Transaction, error) {
	var transactions []types.Transaction
	result := config.DB.Preload("Book").Where("user_id = ?", userID).Find(&transactions)
	if result.Error != nil {
		return nil, result.Error
	}
	return transactions, nil
}

func calculateDueDate(startDate time.Time) *time.Time {
	dueDate := startDate.AddDate(0, 0, 14)
	return &dueDate
}

func SendAdminRequest(userID uint) error {
	var existingRequest types.AdminRequest
	result := config.DB.Where("user_id = ? AND status = ?", userID, "pending").First(&existingRequest)
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return errors.New("you already have a pending admin request")
	}

	newRequest := types.AdminRequest{
		UserID: userID,
		Status: "pending",
	}
	result = config.DB.Create(&newRequest)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func SearchBooks(query string) ([]types.Book, error) {
	var books []types.Book
	searchQuery := "%" + query + "%"
	result := config.DB.Where("title LIKE ? OR author LIKE ?", searchQuery, searchQuery).Find(&books)
	if result.Error != nil {
		return nil, result.Error
	}
	return books, nil
}