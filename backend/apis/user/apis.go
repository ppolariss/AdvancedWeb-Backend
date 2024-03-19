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

	var user User
	err = common.ValidateBody(c, &user)
	if err != nil {
		return err
	}
	if tmpUser.ID != user.ID {
		return common.Forbidden("You can only update your own info")
	}
	if len(user.Password) != 0 {
		user.Password = utils.MakePassword(user.Password)
	}
	return user.Update()
}
