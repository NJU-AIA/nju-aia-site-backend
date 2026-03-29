package auth

import (
	"errors"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	db.AutoMigrate(&User{})
	return &Repository{db: db}
}

func (r *Repository) Create(user *User) error {
	return r.db.Create(user).Error
}

func (r *Repository) ExistsAdmin() (bool, error) {
	var count int64
	err := r.db.Model(&User{}).Where("role = ?", RoleAdmin).Count(&count).Error
	return count > 0, err
}

func (r *Repository) FindFirstAdmin() (*User, error) {
	var user User
	err := r.db.Where("role = ?", RoleAdmin).Order("id asc").First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) FindByUsername(username string) (*User, error) {
	var user User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) UpdatePasswordHash(userID uint, hash string) error {
	result := r.db.Model(&User{}).Where("id = ?", userID).Update("password_hash", hash)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}
