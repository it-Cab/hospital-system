package database

import (
	"log"
	"os"
	"time"

	"example.com/myapp/app/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := os.Getenv("DB_SOURCE")
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	DB.AutoMigrate(&models.Hospital{}, &models.Patient{}, &models.Staff{})

	seedHospital()
	seedPatient()
}

func seedPatient() {
	var count int64
	DB.Model(&models.Patient{}).Count(&count)

	if count == 0 {
		newPatient := models.Patient{
			ID:          "001",
			PatientHN:   "HN123456",
			HospitalID:  "1",
			FirstNameTH: "สมชาย",
			LastNameTH:  "รักดี",
			FirstNameEN: "Somchai",
			LastNameEN:  "Rakdee",
			DateOfBirth: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
			NationalID:  "1234567890123",
			PhoneNumber: "0812345678",
			Gender:      "M",
		}

		result := DB.Create(&newPatient)
		if result.Error != nil {
			log.Println("Failed to seed patient:", result.Error)
		} else {
			log.Println("Successfully seeded first patient data!")
		}
	}
}

func seedHospital() {
    var count int64
    DB.Model(&models.Hospital{}).Count(&count)

    if count == 0 {
        hospitals := []models.Hospital{
            {
                ID:      "1",
                Name:    "BKK Hospital",
                Address: "Bangkok, Thailand",
            },
            {
                ID:      "2",
                Name:    "Bangna Medical",
                Address: "Samut Prakan, Thailand",
            },
        }

        for _, h := range hospitals {
            DB.Create(&h)
        }
        log.Println("Successfully seeded hospitals!")
    }
}
