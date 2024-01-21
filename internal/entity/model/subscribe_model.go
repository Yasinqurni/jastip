package model

import "time"

type SubScribeStatusEnum string

const (
	STATUS_WAITING SubScribeStatusEnum = "waiting"
	STATUS_ACTIVE  SubScribeStatusEnum = "active"
	STATUS_EXPIRED SubScribeStatusEnum = "expired"
)

type SubscribeModel struct {
	ID         string
	UserID     string
	ClassID    string
	Status     SubScribeStatusEnum
	ExpireDate time.Time
	active_at  time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  time.Time

	//relation
	User  *UserModel
	Class *SubscribeClassModel
}
