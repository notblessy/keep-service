package model

import (
	"encoding/json"
	"time"
)

type SplitEntity struct {
	ID          string    `gorm:"primaryKey" json:"id"`
	OwnerID     string    `json:"owner_id"`
	OwnerDetail User      `json:"owner_detail" gorm:"foreignKey:OwnerID"`
	SplitMates  string    `json:"split_mates"`
	CreatedAt   time.Time `json:"created_at"`
}

type Split struct {
	ID          string      `json:"id"`
	OwnerID     string      `json:"owner_id"`
	OwnerDetail User        `json:"owner_detail" gorm:"foreignKey:OwnerID"`
	SplitMates  []SplitMate `json:"split_mates"`
	CreatedAt   time.Time   `json:"created_at"`
}

type SplitMate struct {
	ID         string      `json:"id"`
	Name       string      `json:"name"`
	GrandTotal int         `json:"grand_total"`
	SplitItems []SplitItem `json:"split_items"`
}

type SplitItem struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Qty   int    `json:"qty"`
	Price int    `json:"price"`
	Total int    `json:"total"`
}

func (s *Split) ToEntity() SplitEntity {
	json, err := json.Marshal(s.SplitMates)
	if err != nil {
		return SplitEntity{}
	}

	split := SplitEntity{
		ID:          s.ID,
		OwnerID:     s.OwnerID,
		OwnerDetail: s.OwnerDetail,
		SplitMates:  string(json),
	}

	return split
}
