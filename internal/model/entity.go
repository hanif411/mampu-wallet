package model

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Wallet   Wallet `gorm:"foreignKey:UserID"`
}

type Wallet struct {
	ID           uint          `gorm:"primaryKey"`
	UserID       uint          `gorm:"not null"`
	Balance      int64         `gorm:"default:0"`
	Pin          string        `gorm:"not null"`
	Signature    string        `gorm:"not null"`
	Transactions []Transaction `gorm:"foreignKey:WalletID"`
}

type Transaction struct {
	ID       uint   `gorm:"primaryKey"`
	WalletID uint   `gorm:"not null"`
	Amount   int64  `gorm:"not null"`
	Type     string `gorm:"type:varchar(10)"`
}
