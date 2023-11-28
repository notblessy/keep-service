package model

import "time"

type Bill struct {
	ID           string    `gorm:"primaryKey" json:"id" validate:"required"`
	OwnerID      string    `json:"owner_id" validate:"required"`
	Name         string    `json:"name" validate:"required"`
	Qty          int       `json:"qty" validate:"required"`
	Taken        int       `json:"taken" gorm:"-:all"`
	Price        int       `json:"price" validate:"required"`
	SplitPayment bool      `json:"split_payment"`
	ExpiredAt    time.Time `gorm:"->;<-:create" json:"expired_at"`
	CreatedAt    time.Time `gorm:"->;<-:create" json:"created_at"`
	Shares       []Share   `gorm:"->;-:migration" json:"shares"`
}

type AssignedBillResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name" `
	Me    bool   `json:"me" `
	Total int    `json:"total"`
	Bills []Bill `json:"bills,omitempty"`
}

func NewAssignedBill(bill Bill) AssignedBillResponse {
	var total int
	var assigneeID string
	var assigneeName string

	for _, share := range bill.Shares {
		assigneeID = share.MateID
		assigneeName = share.MateName

		total += share.Price
	}
	return AssignedBillResponse{
		ID:    assigneeID,
		Name:  assigneeName,
		Total: total,
	}
}
