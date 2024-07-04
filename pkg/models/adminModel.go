package models

import (
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
	result := config.DB.Delete(&types.Book{}, id)
	return result.Error
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