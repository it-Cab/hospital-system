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

func createTestToken(hospitalId uint, expired bool) string {
	expiration := time.Now().Add(time.Hour * 1).Unix()
	if expired {
		expiration = time.Now().Add(-time.Hour * 1).Unix()
	}

	claims := jwt.MapClaims{
		"hospital_id": hospitalId,
		"exp":      expiration,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, _ := token.SignedString(jwtKey)
	return tokenString
}

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	if len(jwtKey) == 0 {
		jwtKey = []byte("test_secret_key")
	}

	t.Run("No Header Authorization - 401", func(t *testing.T) {
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

	t.Run("Invalid or Expire Token - 401", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)

		r.Use(AuthMiddleware())
		r.GET("/test", func(ctx *gin.Context) {
			ctx.Status(http.StatusOK)
		})

		expiredToken := createTestToken(01, true) //expire token

		c.Request, _ = http.NewRequest("GET", "/test", nil)
		c.Request.Header.Set("Authorization", "Bearer " + expiredToken)
		r.ServeHTTP(w, c.Request)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Token ไม่ถูกต้องหรือหมดอายุ")
	})

	t.Run("Valid Token - 200", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)

		var hospitalInContext uint

		r.Use(AuthMiddleware())
		r.GET("/test", func(ctx *gin.Context) {
			h, _ := ctx.Get("hospital_id")
			hospitalInContext = h.(uint)
			ctx.Status(http.StatusOK)
		})

		validToken := createTestToken(01, false)

		c.Request, _ = http.NewRequest("GET", "/test", nil)
		c.Request.Header.Set("Authorization", "Bearer "+ validToken)
		r.ServeHTTP(w, c.Request)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, uint(01), hospitalInContext)
	})
}