package repository

import (
	"go-gin-with-jwt-authentication-and-validation/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *models.User) (*models.User, error)
	FindUserByUsername(username string) (*models.User)
	FindByUserID(userID any) (*models.User, error)
	UpdateUser(user *models.User) error
	DeactivateUser(userID any) error
	FindAllUsers() ([]models.User, error)	
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (ur *userRepository) CreateUser(user *models.User) (*models.User, error) {
	if err := ur.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *userRepository) FindUserByUsername(username string) (*models.User) {
	var user models.User
	ur.db.Where("username = ? AND is_active = ?", username, true).First(&user)
	return &user
}

func (ur *userRepository) FindByUserID(userID any) (*models.User, error) {
	var user models.User
	if err := ur.db.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *userRepository) UpdateUser(user *models.User) error {
	return ur.db.Save(user).Error
}

func (ur *userRepository) DeactivateUser(userID any) error {
	return ur.db.Model(&models.User{}).Where("id = ?", userID).Update("is_active", false).Error
}

func (r *userRepository) FindAllUsers() ([]models.User, error) {
    var users []models.User
    if err := r.db.Find(&users).Error; err != nil {
        return nil, err
    }
    return users, nil
}
