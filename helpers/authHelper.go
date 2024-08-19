package helpers

import (
	"errors"
	"github.com/gin-gonic/gin"
)

func MatchUserTypeToUid(context *gin.Context, userId string) (err error) {
	userType := context.GetString("user_type")
	uid := context.GetString("uid")
	err = nil

	if userType == "ROLE_USER" || uid != userId {
		err = errors.New("access denied")
		return err
	}
	err = CheckUserType(context, userType)
	return err
}

func CheckUserType(context *gin.Context, role string) (err error) {
	userType := context.GetString("user_type")
	err = nil
	if userType != role {
		err = errors.New("access denied")
		return err
	}
	return err
}
