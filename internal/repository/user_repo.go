package repository

import (
	"go-finance-wallet/internal/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateWithWallet(user *model.User, hashedPin string, initialSignature string) error
	GetByUsername(username string) (*model.User, error)
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{db}
}

func (r *userRepo) CreateWithWallet(user *model.User, hashedPin string, initialSignature string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}

		wallet := model.Wallet{
			UserID:    user.ID,
			Balance:   0,
			Pin:       hashedPin,
			Signature: "pending",
		}

		if err := tx.Create(&wallet).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *userRepo) GetByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	return &user, err
}
