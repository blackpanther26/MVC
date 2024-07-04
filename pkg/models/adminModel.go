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

    var pendingCount int64
    tx.Model(&types.Transaction{}).Where("book_id = ? AND status = ?", id, "pending").Count(&pendingCount)
    if pendingCount > 0 {
        tx.Rollback()
        return errors.New("cannot delete book because it has pending transactions")
    }

    var book types.Book
    tx.First(&book, id)
    if book.TotalCopies > 0 {
        var borrowedCount int64
        tx.Model(&types.Transaction{}).Where("book_id = ? AND status = ?", id, "checked_out").Count(&borrowedCount)
        if borrowedCount > 0 {
            tx.Rollback()
            return errors.New("cannot delete book because some copies are checked out")
        }
    }

    result := tx.Delete(&types.Book{}, id)
    if result.Error != nil {
        tx.Rollback()
        return result.Error
    }

    return tx.Commit().Error
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
	transaction.Status = status
	return config.DB.Save(&transaction).Error
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