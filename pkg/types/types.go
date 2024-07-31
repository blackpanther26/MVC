package types

import (
	"errors"
	"gorm.io/gorm"
	"time"
	"github.com/shopspring/decimal"
)

type User struct {
	ID           uint   `gorm:"primaryKey"`
	Username     string `gorm:"size:50;uniqueIndex" validate:"required"`
	PasswordHash string `gorm:"size:256" validate:"required"`
	IsAdmin      bool   `gorm:"default:false"`
}

type Book struct {
	ID          uint   `gorm:"primaryKey"`
	Title       string `gorm:"size:255" validate:"required"`
	Author      string `gorm:"size:255" validate:"required"`
	ISBN        string `gorm:"size:13;uniqueIndex" validate:"required"`
	TotalCopies int    `gorm:"default:1" validate:"required"`
	CheckedOutCopies  int    `gorm:"default:0"`
}

type Transaction struct {
	ID              uint            `gorm:"primaryKey"`
	UserID          uint            `gorm:"index"`
	BookID          uint            `gorm:"index"`
	TransactionType string          `gorm:"size:10" validate:"required"`
	TransactionDate time.Time       `gorm:"default:CURRENT_TIMESTAMP"`
	DueDate         *time.Time
	ReturnDate      *time.Time
	Fine            decimal.Decimal `gorm:"type:decimal(10,2);default:0"`
	Status          string          `gorm:"default:'pending';size:255" validate:"required"`
	User            User            `gorm:"foreignKey:UserID" json:"-"`
	Book            Book            `gorm:"foreignKey:BookID" json:"-"`
}

type AdminRequest struct {
	ID     uint   `gorm:"primaryKey"`
	UserID uint   `gorm:"index"`
	Status string `gorm:"size:255;default:'pending'" validate:"required"`
	User   User   `gorm:"foreignKey:UserID"`
}

var validStatuses = []string{"pending", "approved", "denied"}

func (ar *AdminRequest) BeforeSave(tx *gorm.DB) (err error) {
	if !isValidStatus(ar.Status) {
		return errors.New("invalid status")
	}
	return
}

func isValidStatus(status string) bool {
	for _, s := range validStatuses {
		if s == status {
			return true
		}
	}
	return false
}
