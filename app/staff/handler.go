package staff

import (
	"net/http"
	"os"
	"time"

	"example.com/myapp/app/database"
	"example.com/myapp/app/model"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

func StaffCreate(c *gin.Context) {
	var input struct {
		Username   string `json:"username" binding:"required"`
		Password   string `json:"password" binding:"required"`
		HospitalID string `json:"hospital_id" binding:"required"`
		FullName   string `json:"full_name"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "กรุณากรอกข้อมูลให้ครบถ้วน"})
		return
	}

	newStaff := models.Staff{
		Username:   input.Username,
		Password:   input.Password,
		HospitalID: input.HospitalID,
		FullName:   input.FullName,
	}

	if err := database.DB.Create(&newStaff).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถสร้างบัญชีได้ หรือ Username นี้มีอยู่แล้ว"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "สร้างบัญชี Staff สำเร็จ", "username": newStaff.Username})
}

func StaffLogin(c *gin.Context) {
	var credentials struct {
		Username   string `json:"username" binding:"required"`
		Password   string `json:"password" binding:"required"`
		HospitalID string `json:"hospital_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ข้อมูลไม่ถูกต้อง"})
		return
	}

	var staff models.Staff
	result := database.DB.Preload("Hospital").
		Where("username = ? AND password = ? AND hospital_id = ?",
			credentials.Username, credentials.Password, credentials.HospitalID).First(&staff)

	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username, Password หรือ HospitalID ไม่ถูกต้อง"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username":    staff.Username,
		"hospital_id": staff.HospitalID,
		"exp":         time.Now().Add(time.Hour * 24).Unix(), // Expire in 24 hr.
	})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถสร้าง Token ได้"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   tokenString,
	})
}
