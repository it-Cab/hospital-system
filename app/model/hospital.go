package models

import (
	"time"

	"gorm.io/gorm"
)

type Hospital struct {
	Name      string         `gorm:"size:255;not null;unique" json:"name"`
	ID        string         `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Address   string         `json:"address"`

	Staffs   []Staff   `gorm:"foreignKey:HospitalID" json:"staffs,omitempty"`
	Patients []Patient `gorm:"foreignKey:HospitalID" json:"patients,omitempty"`
}
