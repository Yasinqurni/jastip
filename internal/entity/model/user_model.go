package model

import "time"

type UserModel struct {
	ID          string
	Name        string
	Address     string
	Email       string
	PhoneNumber string
	Password    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time

	//relation
	Subscribe *SubscribeModel `gorm:"-"`
}

func (UserModel) TableName() string {
	return "users"
}
