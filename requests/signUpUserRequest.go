package request

import (
	"go-gin-with-jwt-authentication-and-validation/repository"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

var validate *validator.Validate

func init() {
    validate = validator.New()
}

type SignupRequest struct {
    Username string `json:"username" validate:"required,min=3,max=32,unique_username"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=6"`
}


func (r *SignupRequest) Validate(repo repository.UserRepository) error {
    validate := validator.New()
    validate.RegisterValidation("unique_username", UniqueUsernameValidation(repo))
    return validate.Struct(r)
}

type UpdateUserRequest struct {
    Username       string `json:"username" validate:"required,min=3,max=32,unique_username"`
    CurrentPassword string `json:"current_password" validate:"required_with=Password"`
    Password       string `json:"password" validate:"omitempty,min=6"`
}

func CurrentPasswordValidation(repo repository.UserRepository, userID *uint) validator.Func {
    return func(fl validator.FieldLevel) bool {
        if userID == nil {
            return false
        }
        currentPassword := fl.Field().String()
        user, _ := repo.FindByUserID(*userID) 
        if user == nil {
            return false
        }
        err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(currentPassword))
        return err == nil
    }
}


func (r *UpdateUserRequest) Validate(repo repository.UserRepository, userID *uint) error {
    validate := validator.New()
    validate.RegisterValidation("unique_username", UniqueUsernameValidation(repo, userID))
    validate.RegisterValidation("current_password_valid", CurrentPasswordValidation(repo, userID))
    return validate.Struct(r)
}

func UniqueUsernameValidation(repo repository.UserRepository, userID *uint) validator.Func {
    return func(fl validator.FieldLevel) bool {
        username := fl.Field().String()
        user, _ := repo.FindUserByUsername(username)
        return user == nil || (userID != nil && user.ID == *userID)
    }
}
