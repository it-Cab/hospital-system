package staff

import (
	// "os"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	// "time"

	"example.com/myapp/app/database"
	"example.com/myapp/app/model"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	// "github.com/golang-jwt/jwt/v5"
)

// SetupTestDB
func SetupTestDB() {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&models.Staff{})
	database.DB = db
}

func TestStaffCreate(t *testing.T) {
	SetupTestDB()
	gin.SetMode(gin.TestMode)

	t.Run("Create Staff Success", func(t *testing.T) {
		r := gin.Default()
		r.POST("/staff/add", StaffCreate)

		staffData := map[string]interface{}{
			"username":    "admin01",
			"password":    "password123",
			"hospital_id": "01",
			"full_name":   "John Doe",
		}
		body, _ := json.Marshal(staffData)

		req, _ := http.NewRequest("POST", "/staff/add", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), "สร้างบัญชี Staff สำเร็จ")
	})

	t.Run("Create Staff Fail Case Incomplete Data", func(t *testing.T) {
		r := gin.Default()
		r.POST("/staff/add", StaffCreate)

		staffData := map[string]interface{}{
			"username": "admin01",
			"password": "password123",
		}
		body, _ := json.Marshal(staffData)

		req, _ := http.NewRequest("POST", "/staff/add", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "กรุณากรอกข้อมูลให้ครบถ้วน")
	})

	t.Run("Create Staff Fail Case Duplicate", func(t *testing.T) {
		// mock Staff DB
		database.DB.Create(&models.Staff{
			Username: "admin01", Password: "password123", HospitalID: "01",
		})
		r := gin.Default()
		r.POST("/staff/add", StaffCreate)

		staffData := map[string]interface{}{
			"username":    "admin01",
			"password":    "password123",
			"hospital_id": "01",
			"full_name":   "John Doe",
		}
		body, _ := json.Marshal(staffData)

		req, _ := http.NewRequest("POST", "/staff/add", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "ไม่สามารถสร้างบัญชีได้ หรือ Username นี้มีอยู่แล้ว")
	})
}

func TestStaffLogin(t *testing.T) {
	SetupTestDB()
	gin.SetMode(gin.TestMode)

	database.DB.Create(&models.Staff{
		Username: "admin01", Password: "password123", HospitalID: "01",
	})
	t.Run("Login Success", func(t *testing.T) {
		r := gin.Default()
		r.POST("/staff/login", StaffLogin)

		staffData := map[string]interface{}{
			"username":    "admin01",
			"password":    "password123",
			"hospital_id": "01",
		}
		body, _ := json.Marshal(staffData)

		req, _ := http.NewRequest("POST", "/staff/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Login successful")
	})

	t.Run("Login Fail", func(t *testing.T) {
		r := gin.Default()
		r.POST("/staff/login", StaffLogin)

		staffData := map[string]interface{}{
			"username":    "admin01",
			"password":    "password",
			"hospital_id": "01",
		}
		body, _ := json.Marshal(staffData)

		req, _ := http.NewRequest("POST", "/staff/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Username, Password หรือ HospitalID ไม่ถูกต้อง")
	})

	t.Run("Login Fail Case Incomplete Data", func(t *testing.T) {
		r := gin.Default()
		r.POST("/staff/login", StaffLogin)

		staffData := map[string]interface{}{
			"username": "admin01",
			"password": "password",
		}
		body, _ := json.Marshal(staffData)

		req, _ := http.NewRequest("POST", "/staff/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

}
