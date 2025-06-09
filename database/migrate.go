package database

import (
	"log"

	"washo.com/main/repositories"
)

func TablesExist() bool {
	db := GetDB()

	// Check if tables exist
	hasUser := db.Migrator().HasTable(&repositories.User{})
	hasVehicle := db.Migrator().HasTable(&repositories.Vehicle{})

	return hasUser && hasVehicle
}

func Migrate() error {
	db := GetDB()

	// Run migrations
	err := db.AutoMigrate(
		&repositories.User{},
		&repositories.Vehicle{}, // Note: Changed from Vehicles to Vehicle to match model name
	)

	if err != nil {
		log.Printf("Failed to migrate database: %v", err)
		return err
	}

	log.Println("Database migration completed successfully")
	return nil
}
