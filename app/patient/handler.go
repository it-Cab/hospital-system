package patient

import (
	"fmt"
	"net/http"
	"time"

	"example.com/myapp/app/database"
	"example.com/myapp/app/model"
	"github.com/gin-gonic/gin"
)

func GetPatientByID(c *gin.Context) {
	id := c.Param("id")
	staffHospital, _ := c.Get("hospital")

	var patient models.Patient
	fmt.Println("Searching for ID:", id)
	result := database.DB.Where("id = ? AND hospital = ?", id, staffHospital).First(&patient)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "ไม่พบข้อมูลคนไข้ที่ระบุ",
		})
		return
	}
	c.JSON(http.StatusOK, patient)
}

func GetPatients(c *gin.Context) {
	staffHospital, _ := c.Get("hospital")
	var input struct {
		NationalID  string `json:"national_id"`
		PassportID  string `json:"passport_id"`
		FirstName   string `json:"first_name"`
		MiddleName  string `json:"middle_name"`
		LastName    string `json:"last_name"`
		DateOfBirth string `json:"date_of_birth"`
		Email       string `json:"email"`
	}

	if err := c.ShouldBindJSON(&input); err != nil && c.Request.ContentLength > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ข้อมูล Input ไม่ถูกต้อง"})
		return
	}
	query := database.DB.Model(&models.Patient{}).Where("hospital = ?", staffHospital)

	if input.NationalID != "" {
		query = query.Where("national_id = ?", input.NationalID)
	}
	if input.PassportID != "" {
		query = query.Where("passport_id = ?", input.PassportID)
	}
	if input.FirstName != "" {
		query = query.Where("first_name_th LIKE ? OR first_name_en LIKE ?", "%"+input.FirstName+"%", "%"+input.FirstName+"%")
	}
	if input.MiddleName != "" {
		query = query.Where("middle_name_th = ? OR middle_name_en = ?", input.MiddleName, input.MiddleName)
	}
	if input.LastName != "" {
		query = query.Where("last_name_th LIKE ? OR last_name_en LIKE ?", "%"+input.LastName+"%", "%"+input.LastName+"%")
	}
	if input.DateOfBirth != "" {
		query = query.Where("date_of_birth = ?", input.DateOfBirth)
	}
	if input.Email != "" {
		query = query.Where("email = ?", input.Email)
	}

	var patients []models.Patient
	result := query.Find(&patients)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "ไม่พบข้อมูลคนไข้",
		})
		return
	}
	c.JSON(http.StatusOK, patients)
}

func CreatePatient(c *gin.Context) {
	var input struct {
		ID           string    `json:"id" binding:"required"`
		PatientHN    string    `json:"patient_hn" binding:"required"`
		Hospital     string    `json:"hospital" binding:"required"`
		FirstNameTH  string    `json:"first_name_th"`
		MiddleNameTH string    `json:"middle_name_th"`
		LastNameTH   string    `json:"last_name_th"`
		FirstNameEN  string    `json:"first_name_en"`
		MiddleNameEN string    `json:"middle_name_en"`
		LastNameEN   string    `json:"last_name_en"`
		DateOfBirth  time.Time `json:"date_of_birth"`
		NationalID   string    `json:"national_id"`
		PassportID   string    `json:"passport_id"`
		PhoneNumber  string    `json:"phone_number"`
		Email        string    `json:"email"`
		Gender       string    `json:"gender"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "กรุณากรอกข้อมูลให้ครบถ้วน"})
		return
	}
	staffHospital, _ := c.Get("hospital")
	if input.Hospital != staffHospital {
		c.JSON(http.StatusForbidden, gin.H{"error": "คุณไม่มีสิทธิ์เพิ่มข้อมูลให้โรงพยาบาลอื่น"})
		return
	}

	newPatient := models.Patient{
		ID:           input.ID,
		PatientHN:    input.PatientHN,
		Hospital:     input.Hospital,
		FirstNameTH:  input.FirstNameTH,
		MiddleNameTH: input.MiddleNameTH,
		LastNameTH:   input.LastNameTH,
		FirstNameEN:  input.FirstNameEN,
		MiddleNameEN: input.MiddleNameEN,
		LastNameEN:   input.LastNameEN,
		DateOfBirth:  input.DateOfBirth,
		NationalID:   input.NationalID,
		PassportID:   input.PassportID,
		PhoneNumber:  input.PhoneNumber,
		Email:        input.Email,
		Gender:       input.Gender,
	}

	if err := database.DB.Create(&newPatient).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถเพิ่มข้อมูลคนไข้ได้: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "เพิ่มข้อมูลคนไข้สำเร็จ",
		"patient_id": newPatient.ID,
		"patient_hn": newPatient.PatientHN,
	})
}
