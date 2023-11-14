package model

import "encoding/json"

type SplitEntity struct {
	ID         string `gorm:"primaryKey" json:"id"`
	OwnerID    string `json:"owner_id"`
	SplitMates string `json:"split_mates"`
}

type Split struct {
	ID         string      `json:"id"`
	OwnerID    string      `json:"owner_id"`
	SplitMates []SplitMate `json:"split_mates"`
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
	json, err := json.Marshal(s)
	if err != nil {
		return SplitEntity{}
	}

	split := SplitEntity{
		ID:         s.ID,
		OwnerID:    s.OwnerID,
		SplitMates: string(json),
	}

	return split
}
