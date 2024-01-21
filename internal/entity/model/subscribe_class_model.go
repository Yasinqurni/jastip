package model

import "time"

type SubscribeClassModel struct {
	ID          string
	Name        string
	Time        string
	Price       float64
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time
}
