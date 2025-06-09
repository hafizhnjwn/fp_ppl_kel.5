package repositories

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Vehicle struct {
	ID            string `json:"id"`
	UserID        uint   `json:"user_id"`           // Foreign key field
	User          User   `gorm:"foreignKey:UserID"` // Association
	Queue         int    `json:"queue"`
	Name          string `json:"name"`
	Plate         string `json:"plate"`
	Process       string `json:"process"`
	Contact       string `json:"contact"`
	Date          string `json:"date"`
	Time     string `json:"time"`
}

type CreateVehicleRequest struct {
	UID     string `form:"username"`
	Name    string `form:"name" binding:"required"`
	Plate   string `form:"plate" binding:"required"`
	Contact string `form:"contact"`
	Process string `form:"process"`
}

type VehicleRepository struct {
	db *gorm.DB
}

func NewVehicleRepository(db *gorm.DB) *VehicleRepository {
	return &VehicleRepository{db: db}
}

func (r *VehicleRepository) Create(vehicle *CreateVehicleRequest) (string, error) {
	// Get count of user's vehicles to generate ID
	var count int64
	if err := r.db.Model(&Vehicle{}).Where("user_id = ?", vehicle.UID).Count(&count).Error; err != nil {
		return "", err
	}

	today := time.Now().Format("2006-01-02")

	// Get count of vehicles created today to generate the queue number
	var countqueue int64
	if err := r.db.Model(&Vehicle{}).Where("date = ?", today).Count(&countqueue).Error; err != nil {
		return "", err
	}

	// Fetch the user's username from the User model
	var user User
	if err := r.db.Where("id = ?", vehicle.UID).First(&user).Error; err != nil {
		return "", err
	}
	id := fmt.Sprintf("%s-%d", user.Username, count+1)

	// Create a new vehicle instance with ID format (username-vehiclecount)
	newVehicle := Vehicle{
		ID:      id,
		UserID:  user.ID, // Set the UserID foreign key
		Name:    vehicle.Name,
		Plate:   vehicle.Plate,
		Contact: vehicle.Contact,
		Process: "Waiting",
		// Date:    "2025-01-20",
		Date:          time.Now().Format("2006-01-02"),
		Time:     time.Now().Format("3:04 PM"),
		Queue:         int(countqueue + 1), // Set the queue number
	}
	return id, r.db.Create(&newVehicle).Error
}

func (r *VehicleRepository) FindByProcess(process string) ([]Vehicle, error) {
	var vehicles []Vehicle
	today := time.Now().Format("2006-01-02") //"2025-01-20"
	err := r.db.Where("process = ? AND date = ?", process, today).Preload("User").Find(&vehicles).Error
	return vehicles, err
}

func (r *VehicleRepository) FindByUsername(username string) ([]Vehicle, error) {
	var vehicles []Vehicle
	err := r.db.
		Joins("JOIN users ON users.id = vehicles.user_id").
		Where("users.username = ?", username).
		Preload("User").
		Find(&vehicles).Error
	return vehicles, err
}

func (r *VehicleRepository) FindByID(id string) (*Vehicle, error) {
	var vehicle Vehicle
	err := r.db.Where("id = ?", id).Preload("User").First(&vehicle).Error
	return &vehicle, err
}
func (r *VehicleRepository) Update(id string, vehicle *CreateVehicleRequest) error {
	var existingVehicle Vehicle
	if err := r.db.Where("id = ?", id).First(&existingVehicle).Error; err != nil {
		return err
	}

	// Update the fields of the existing vehicle
	existingVehicle.Name = vehicle.Name
	existingVehicle.Process = vehicle.Process
	existingVehicle.Plate = vehicle.Plate
	existingVehicle.Contact = vehicle.Contact

	return r.db.Save(&existingVehicle).Error
}

func (r *VehicleRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&Vehicle{}).Error
}