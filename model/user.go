package model

import "time"

type User struct {
	ID        string     `gorm:"primaryKey" json:"id"`
	Email     string     `json:"email"`
	Name      string     `json:"name" validate:"required"`
	CreatedAt time.Time  `gorm:"->;<-:create" json:"created_at"`
	Mates     []UserMate `gorm:"->;-:migration;foreignKey:OwnerID" json:"mates"`
}

type UserMate struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	OwnerID   string    `json:"owner_id" validate:"required"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `gorm:"->;<-:create" json:"created_at"`
	User      User      `gorm:"->;-:migration;foreignKey:UserID" json:"user_detail"`
}

type Mate struct {
	ID     string  `gorm:"primaryKey" json:"id"`
	Name   string  `json:"name" validate:"required"`
	Shares []Share `gorm:"->;-:migration" json:"shares"`
}
