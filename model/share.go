package model

import "time"

type Share struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	OwnerID   string    `json:"owner_id" validate:"required"`
	BillID    string    `json:"bill_id" validate:"required"`
	MateID    string    `json:"mate_id" validate:"required"`
	MateName  string    `json:"mate_name" validate:"required"`
	Qty       int       `json:"qty" validate:"required"`
	Price     int       `json:"price" validate:"required"`
	CreatedAt time.Time `gorm:"->;<-:create" json:"created_at"`
	Bill      Bill      `gorm:"->;-:migration;foreignKey:BillID" json:"bill_detail" validate:"omitempty"`
}

type BulkShare struct {
	BillID  string      `json:"bill_id" validate:"required"`
	OwnerID string      `json:"owner_id" validate:"required"`
	Mates   []ShareMate `json:"mates" validate:"required"`
}

type ShareMate struct {
	MateID   string `json:"mate_id"`
	MateName string `json:"mate_name"`
}

type TotalShareQty struct {
	BillID string `json:"bill_id"`
	Qty    int    `json:"qty"`
}
