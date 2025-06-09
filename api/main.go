package main

import (
	"html/template"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"washo.com/main/database"
	"washo.com/main/handlers"
	"washo.com/main/middleware"
	"washo.com/main/repositories"
	"washo.com/main/services"
)

var db *gorm.DB

func contains(substring, str string) bool {
	return strings.Contains(str, substring)
}

func init() {
	database.LoadEnvs()
	database.InitializeDatabaseLayer()

	// Check if tables exist first
	if !database.TablesExist() {
		log.Println("Tables do not exist. Running migrations...")
		if err := database.Migrate(); err != nil {
			log.Fatalf("Failed to migrate database: %v", err)
		}
		log.Println("Migrations completed successfully")
	} else {
		log.Println("Tables already exist. Skipping migrations")
	}
	db = database.GetDB()
}

func main() {
	// Create repository
	vehicleRepo := repositories.NewVehicleRepository(db)

	// Create service
	vehicleService := services.NewVehicleService(vehicleRepo)

	// Create handler
	vehicleHandler := handlers.NewVehicleHandler(vehicleService)

	// setup gin router
	router := gin.Default()
	router.Use(gin.Logger())

	// Load HTML templates
	router.SetFuncMap(template.FuncMap{
		"contains": contains, // Now you can use {{contains}} in templates
	})
	router.LoadHTMLGlob("templates/*")
	// Auth routes
	auth := router.Group("/")
	{
		auth.GET("", handlers.Home)
		auth.GET("/login", handlers.Login)
		auth.GET("/logout", handlers.Logout)
		auth.GET("/register", handlers.CreateUser)
		auth.POST("/login", handlers.Login)
		auth.POST("/register", handlers.CreateUser)
	}

	// Vehicle routes
	snip := router.Group("/vehicles")
	{
		// Guest routes
		snip.GET("", vehicleHandler.GetVehiclesByProcess)
		snip.GET("/:id", vehicleHandler.GetVehicleByID)

		// Authenticated routes
		snip.GET("/new", middleware.CheckAuth, vehicleHandler.CreateVehicle)
		snip.POST("/new", middleware.CheckAuth, vehicleHandler.CreateVehicle)
		snip.GET("/:id/edit", middleware.CheckAuth, vehicleHandler.UpdateVehicle)
		snip.POST("/:id/edit", middleware.CheckAuth, vehicleHandler.UpdateVehicle)
		snip.POST("/:id/delete", middleware.CheckAuth, vehicleHandler.DeleteVehicle)
		snip.GET("/:id/delete", middleware.CheckAuth, vehicleHandler.DeleteVehicle)
	}

	// start server
	log.Println("starting server on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
