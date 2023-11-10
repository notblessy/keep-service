package model

import "time"

type Bill struct {
	ID        string    `gorm:"primaryKey" json:"id" validate:"required"`
	OwnerID   string    `json:"owner_id" validate:"required"`
	Name      string    `json:"name" validate:"required"`
	Qty       int       `json:"qty" validate:"required"`
	Price     int       `json:"price" validate:"required"`
	CreatedAt time.Time `gorm:"->;<-:create" json:"created_at"`
}
