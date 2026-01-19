package models

import "time"

type Patient struct {
	ID        string `gorm:"primaryKey" json:"id"`
	PatientHN string `json:"patient_hn"`
	Hospital  string `json:"hospital"`

	FirstNameTH  string `gorm:"size:100" json:"first_name_th"`
	MiddleNameTH string `gorm:"size:100" json:"middle_name_th"`
	LastNameTH   string `gorm:"size:100" json:"last_name_th"`

	FirstNameEN  string `gorm:"size:100" json:"first_name_en"`
	MiddleNameEN string `gorm:"size:100" json:"middle_name_en"`
	LastNameEN   string `gorm:"size:100" json:"last_name_en"`

	DateOfBirth time.Time `json:"date_of_birth"`
	NationalID  string    `gorm:"size:13;index" json:"national_id"`
	PassportID  string    `gorm:"size:20;index" json:"passport_id"`

	PhoneNumber string `gorm:"size:20" json:"phone_number"`
	Email       string `gorm:"size:100" json:"email"`
	Gender      string `gorm:"size:1" json:"gender"`
}
