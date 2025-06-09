package repositories

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint      `form:"id" gorm:"primary_key"`
	Username  string    `form:"username" gorm:"unique"`
	Password  string    `form:"password"`
	Admin     bool      `form:"admin"` // New attribute
	Vehicles  []Vehicle `gorm:"foreignKey:UserID"`          // Association
	CreatedAt time.Time
	UpdatedAt time.Time
}

type AuthInput struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) Update(user *User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&User{}, id).Error
}
