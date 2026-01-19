package main

import (
	"example.com/myapp/app/database"
	"example.com/myapp/app/middleware"
	"example.com/myapp/app/patient"
	"example.com/myapp/app/staff"
	"github.com/gin-gonic/gin"
)

func main() {
	database.InitDB()
	r := gin.Default()

	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/patient/search/:id", patient.GetPatientByID)
		protected.GET("/patient/search", patient.GetPatients)
		protected.POST("/patient/add", patient.CreatePatient)
	}

	r.POST("/staff/add", staff.StaffCreate)
	r.POST("/staff/login", staff.StaffLogin)

	r.Run(":8080")
}