package model

import "time"

type RoleUser string

const (
	ROLE_ADMIN RoleUser = "admin"
	ROLE_USER  RoleUser = "user"
)

type UserModel struct {
	ID          string
	Name        string
	Address     string
	Email       string
	PhoneNumber string
	Password    string
	Role        RoleUser
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time

	//relation
	Subscribe *SubscribeModel `gorm:"-"`
}

func (UserModel) TableName() string {
	return "users"
}
