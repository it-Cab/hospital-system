package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

// ฟังก์ชันช่วยสร้าง Token สำหรับใช้ใน Test Case ต่างๆ
func createTestToken(hospital string, expired bool) string {
	expiration := time.Now().Add(time.Hour * 1).Unix()
	if expired {
		expiration = time.Now().Add(-time.Hour * 1).Unix() // ย้อนเวลาให้ Expired
	}

	claims := jwt.MapClaims{
		"hospital": hospital,
		"exp":      expiration,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// หมายเหตุ: jwtKey ต้องถูกประกาศใน package middleware
	tokenString, _ := token.SignedString(jwtKey)
	return tokenString
}

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// กำหนด Key สำหรับ Test (กรณีใน Env ไม่มีค่า)
	if len(jwtKey) == 0 {
		jwtKey = []byte("test_secret_key")
	}

	t.Run("ไม่มี Header Authorization - ควรได้ 401", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)

		r.Use(AuthMiddleware())
		r.GET("/test", func(ctx *gin.Context) {
			ctx.Status(http.StatusOK)
		})

		c.Request, _ = http.NewRequest("GET", "/test", nil)
		r.ServeHTTP(w, c.Request)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "กรุณา Login ก่อนใช้งาน")
	})

	t.Run("Token ไม่ถูกต้องหรือหมดอายุ - ควรได้ 401", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)

		r.Use(AuthMiddleware())
		r.GET("/test", func(ctx *gin.Context) {
			ctx.Status(http.StatusOK)
		})

		// สร้าง Token ที่หมดอายุแล้ว
		expiredToken := createTestToken("BKK Hospital", true)

		c.Request, _ = http.NewRequest("GET", "/test", nil)
		c.Request.Header.Set("Authorization", "Bearer "+expiredToken)
		r.ServeHTTP(w, c.Request)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Token ไม่ถูกต้องหรือหมดอายุ")
	})

	t.Run("Token ถูกต้อง - ควรได้ 200 และมีการ Set Context", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)

		var hospitalInContext string

		r.Use(AuthMiddleware())
		r.GET("/test", func(ctx *gin.Context) {
			// ตรวจสอบว่าค่า hospital ถูกส่งต่อมาถึง Handler จริงไหม
			h, _ := ctx.Get("hospital")
			hospitalInContext = h.(string)
			ctx.Status(http.StatusOK)
		})

		validToken := createTestToken("BKK Hospital", false)

		c.Request, _ = http.NewRequest("GET", "/test", nil)
		c.Request.Header.Set("Authorization", "Bearer "+validToken)
		r.ServeHTTP(w, c.Request)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "BKK Hospital", hospitalInContext)
	})
}