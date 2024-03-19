package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/opentreehole/go-common"
	. "src/models"
	"src/utils"
)

// GetUserInfo @GetUserInfo
// @Router /api/users/data [get]
// @Summary Get user info
// @Description Get user info
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} models.User
// @Failure 403 {object} common.HttpError
// @param Authorization header string true "Bearer和token空格拼接"
func GetUserInfo(c *fiber.Ctx) error {
	tmpUser, err := GetGeneralUser(c)
	if err != nil {
		return err
	}

	getUser, err := GetUserByID(tmpUser.ID)
	if err != nil {
		return err
	}

	return c.JSON(&getUser)
}

// UpdateUser @UpdateUser
// @Router /api/users [put]
// @Summary Update user
// @Description Update user
// @Tags User
// @Accept json
// @Produce json
// @Param json body models.User true "json"
// @Success 200
// @Failure 400 {object} common.HttpError
// @Failure 403 {object} common.HttpError
// @Failure 404 {object} common.HttpError
// @param Authorization header string true "Bearer和token空格拼接"
func UpdateUser(c *fiber.Ctx) error {
	tmpUser, err := GetGeneralUser(c)
	if err != nil {
		return err
	}

	var updateBody UpdateUserRequest
	err = common.ValidateBody(c, &updateBody)
	if err != nil {
		return err
	}

	//var pwdChanged bool

	user, err := GetUserByID(tmpUser.ID)
	if len(updateBody.Password) != 0 {
		user.Password = utils.MakePassword(updateBody.Password)
		//pwdChanged = true
	}
	if updateBody.Age != 0 {
		user.Age = updateBody.Age
	}
	if len(updateBody.Email) != 0 {
		user.Email = updateBody.Email
	}
	if len(updateBody.Phone) != 0 {
		user.Phone = updateBody.Phone
	}
	if len(updateBody.Gender) != 0 {
		user.Gender = updateBody.Gender
	}
	err = user.Update()
	if err != nil {
		return err
	}
	//if pwdChanged {
	//	return DeleteUserJWTSecret(tmpUser.ID)
	//}
	return nil
}
