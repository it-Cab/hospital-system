package models

import "time"

type Staff struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Username   string    `gorm:"unique;not null" json:"username"`
	Password   string    `gorm:"not null" json:"-"`
	Hospital   string    `gorm:"not null" json:"hospital"`
	FullName   string    `json:"full_name"`
	Role       string    `json:"role"`
	CreatedAt  time.Time `json:"created_at"`
}