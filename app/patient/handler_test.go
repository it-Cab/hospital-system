package patient

import (
	"os"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"example.com/myapp/app/database"
	"example.com/myapp/app/middleware"
	"example.com/myapp/app/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"github.com/golang-jwt/jwt/v5"
)

// SetupTestDB
func SetupTestDB() {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&models.Patient{})
	database.DB = db
}

func generateTestToken(hospital string) string {
	claims := jwt.MapClaims{
		"hospital": hospital,
		"username": "testuser",
		"exp":      time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET"))) 
	return t
}

func TestPatientSearchById(t *testing.T) {
    SetupTestDB()
    gin.SetMode(gin.TestMode)

    // mock Patient DB
    database.DB.Create(&models.Patient{
        ID: "001", PatientHN: "HN001", Hospital: "BKK Hospital",
    })

    t.Run("Search Fail Case Not Login", func(t *testing.T) {
        r := gin.Default()
        r.Use(middleware.AuthMiddleware())
        r.GET("/patient/search/:id", GetPatientByID)

        req, _ := http.NewRequest("GET", "/patient/search/001", nil)
        w := httptest.NewRecorder()
        r.ServeHTTP(w, req)

        assert.Equal(t, http.StatusUnauthorized, w.Code)
    })

    t.Run("Search Success", func(t *testing.T) {
        r := gin.Default()
        r.Use(middleware.AuthMiddleware())
        r.GET("/patient/search/:id", GetPatientByID)

        token := generateTestToken("BKK Hospital")
        req, _ := http.NewRequest("GET", "/patient/search/001", nil)
        req.Header.Set("Authorization", "Bearer "+token)
        
        w := httptest.NewRecorder()
        r.ServeHTTP(w, req)

        assert.Equal(t, http.StatusOK, w.Code)
    })

    t.Run("Search Fail Case Invalid Hospital", func(t *testing.T) {
        r := gin.Default()
        r.Use(middleware.AuthMiddleware())
        r.GET("/patient/search/:id", GetPatientByID)

        token := generateTestToken("BNK Hospital")
        req, _ := http.NewRequest("GET", "/patient/search/001", nil)
        req.Header.Set("Authorization", "Bearer "+token)

        w := httptest.NewRecorder()
        r.ServeHTTP(w, req)

        assert.Equal(t, http.StatusNotFound, w.Code)
    })
}

func TestPatientSearch(t *testing.T) {
    SetupTestDB()
    gin.SetMode(gin.TestMode)

    // mock Patient DB
    database.DB.Create(&models.Patient{
        ID: "001", PatientHN: "HN001", Hospital: "BKK Hospital",
    })

    t.Run("Search Fail Case Not Login", func(t *testing.T) {
        r := gin.Default()
        r.Use(middleware.AuthMiddleware())
        r.GET("/patient/search", GetPatients)

        req, _ := http.NewRequest("GET", "/patient/search", nil)
        w := httptest.NewRecorder()
        r.ServeHTTP(w, req)

        assert.Equal(t, http.StatusUnauthorized, w.Code)
    })

    t.Run("Search Success", func(t *testing.T) {
        r := gin.Default()
        r.Use(middleware.AuthMiddleware())
        r.GET("/patient/search", GetPatients)

        token := generateTestToken("BKK Hospital")
        req, _ := http.NewRequest("GET", "/patient/search", nil)
        req.Header.Set("Authorization", "Bearer "+token)
        
        w := httptest.NewRecorder()
        r.ServeHTTP(w, req)

        assert.Equal(t, http.StatusOK, w.Code)
    })

    t.Run("Search Fail Case Invalid Hospital", func(t *testing.T) {
        r := gin.Default()
        r.Use(middleware.AuthMiddleware())
        r.GET("/patient/search", GetPatients)

        token := generateTestToken("BNK Hospital")
        req, _ := http.NewRequest("GET", "/patient/search", nil)
        req.Header.Set("Authorization", "Bearer "+token)

        w := httptest.NewRecorder()
        r.ServeHTTP(w, req)

        assert.Equal(t, http.StatusOK, w.Code)
    })
}

func TestPatientCreate(t *testing.T) {
	SetupTestDB()
	gin.SetMode(gin.TestMode)

	t.Run("Create Patient Success", func(t *testing.T) {
		r := gin.Default()
		r.Use(middleware.AuthMiddleware())
		r.POST("/patient/add", CreatePatient)

		token := generateTestToken("BKK Hospital")
		patientData := map[string]interface{}{
			"id":             "005",
			"patient_hn":     "HN005",
			"hospital":       "BKK Hospital",
			"first_name_th":  "สมหญิง",
			"last_name_th":   "จริงใจ",
			"date_of_birth":  "1995-01-01T00:00:00Z",
			"national_id":    "1234567890123",
			"gender":         "F",
		}
		body, _ := json.Marshal(patientData)
		
		req, _ := http.NewRequest("POST", "/patient/add", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), "เพิ่มข้อมูลคนไข้สำเร็จ")
	})

	t.Run("Create Patient Fail Case Incomplete Input", func(t *testing.T) {
		r := gin.Default()
		r.Use(middleware.AuthMiddleware())
		r.POST("/patient/add", CreatePatient)

		token := generateTestToken("BKK Hospital")
		patientData := map[string]interface{}{
			"id": "006",
		}
		body, _ := json.Marshal(patientData)
		
		req, _ := http.NewRequest("POST", "/patient/add", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Create Patient Fail Case Not Login", func(t *testing.T) {
		r := gin.Default()
		r.Use(middleware.AuthMiddleware())
		r.POST("/patient/add", CreatePatient)

		patientData := map[string]interface{}{
			"id":             "005",
			"patient_hn":     "HN005",
			"hospital":       "BKK Hospital",
			"first_name_th":  "สมหญิง",
			"last_name_th":   "จริงใจ",
			"date_of_birth":  "1995-01-01T00:00:00Z",
			"national_id":    "1234567890123",
			"gender":         "F",
		}
		body, _ := json.Marshal(patientData)
		
		req, _ := http.NewRequest("POST", "/patient/add", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer ")
		
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}