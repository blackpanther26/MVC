package models

import (
	"errors"

	"github.com/blackpanther26/mvc/pkg/config"
	"github.com/blackpanther26/mvc/pkg/types"
)

func AddBook(book *types.Book) error {
	result := config.DB.Create(book)
	return result.Error
}

func UpdateBook(book *types.Book) error {
	result := config.DB.Save(book)
	return result.Error
}

func DeleteBook(id uint) error {
    db := config.GetDB()
    tx := db.Begin()
    if tx.Error != nil {
        return tx.Error
    }

    var book types.Book
    err := tx.First(&book, id).Error
    if err != nil {
        tx.Rollback()
        return errors.New("book not found")
    }

    if book.CheckedOutCopies > 0 {
        tx.Rollback()
        return errors.New("cannot delete book because some copies are checked out")
    }

    result := tx.Delete(&types.Book{}, id)
    if result.Error != nil {
        tx.Rollback()
        return result.Error
    }

    err = tx.Commit().Error
    if err != nil {
        return err
    }

    return nil
}

func GetAllTransactions() ([]types.Transaction, error) {
	var transactions []types.Transaction
	result := config.DB.Preload("User").Preload("Book").Find(&transactions)
	return transactions, result.Error
}

func UpdateTransactionStatus(transactionID uint, status string) error {
    var transaction types.Transaction
    result := config.DB.First(&transaction, transactionID)
    if result.Error != nil {
        return result.Error
    }

    tx := config.DB.Begin()

    transaction.Status = status
    if err := tx.Save(&transaction).Error; err != nil {
        tx.Rollback()
        return err
    }

    var book types.Book
    if err := tx.First(&book, transaction.BookID).Error; err != nil {
        tx.Rollback()
        return err
    }

    switch status {
    case "approved":
        if transaction.TransactionType == "checkout" {
            book.CheckedOutCopies++
        } else if transaction.TransactionType == "checkin" {
            book.CheckedOutCopies--
        }
    }

    if err := tx.Save(&book).Error; err != nil {
        tx.Rollback()
        return err
    }

    tx.Commit()
    return nil
}

func GetAllAdminRequests() ([]types.AdminRequest, error) {
	var adminRequests []types.AdminRequest
	result := config.DB.Preload("User").Find(&adminRequests)
	if result.Error != nil {
		return nil, result.Error
	}
	return adminRequests, nil
}

func UpdateAdminRequestStatus(requestID uint, status string) error {
	var adminRequest types.AdminRequest
	result := config.DB.First(&adminRequest, requestID)
	if result.Error != nil {
		return result.Error
	}
	adminRequest.Status = status
	return config.DB.Save(&adminRequest).Error
}