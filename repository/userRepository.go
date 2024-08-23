package repository

import (
	"go-gin-with-jwt-authentication-and-validation/models"

	"gorm.io/gorm"
)

type UserRepository interface {
    CreateUser(user *models.User) (*models.User, error)
    FindUserByUsername(username string) (*models.User, error)
    UpdateUser(user *models.User) (*models.User, error)
    DeactivateUser(userID int) error
    FindByUserID(userID uint) (*models.User, error)
	FindAllUsers() ([]models.User, error)
}

type userRepository struct {
    DB *gorm.DB
}

func NewUserRepository(DB *gorm.DB) UserRepository {
    return &userRepository{DB: DB}
}

func (r *userRepository) CreateUser(user *models.User) (*models.User, error) {
    if err := r.DB.Create(user).Error; err != nil {
        return nil, err
    }
    return user, nil
}

func (r *userRepository) FindUserByUsername(username string) (*models.User, error) {
    var user models.User
    if err := r.DB.Where("username = ?", username).First(&user).Error; err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *userRepository) UpdateUser(user *models.User) (*models.User, error) {
    if err := r.DB.Save(user).Error; err != nil {
        return nil, err
    }
    return user, nil
}

func (r *userRepository) DeactivateUser(userID int) error {
    if err := r.DB.Model(&models.User{}).Where("id = ?", userID).Update("is_active", false).Error; err != nil {
        return err
    }
    return nil
}

func (r *userRepository) FindByUserID(userID uint) (*models.User, error) {
    var user models.User
    if err := r.DB.First(&user, userID).Error; err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *userRepository) FindAllUsers() ([]models.User, error) {
    var users []models.User
    if err := r.DB.Find(&users).Error; err != nil {
        return nil, err
    }
    return users, nil
}