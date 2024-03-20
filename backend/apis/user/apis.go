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

	user, err := GetUserByID(tmpUser.ID)
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
	return user.Update()
}

// ChangePassword @ChangePassword
// @Router /api/users/password [put]
// @Summary Change password
// @Description Change password by old password or email
// @Tags User
// @Accept json
// @Produce json
// @Param json body ChangePasswordRequest true "json"
// @Success 200 {object} TokenResponse
// @Failure 400 {object} common.HttpError
// @Failure 401 {object} common.HttpError
func ChangePassword(c *fiber.Ctx) error {
	tmpUser, err := GetGeneralUser(c)
	if err != nil {
		return err
	}

	var body ChangePasswordRequest
	err = common.ValidateBody(c, &body)
	if err != nil {
		return err
	}

	user, err := GetUserByID(tmpUser.ID)
	if err != nil {
		return err
	}

	if len(body.OldPassword) > 0 {
		ok := utils.CheckPassword(body.OldPassword, user.Password)
		if !ok {
			return common.Unauthorized("密码错误")
		}
	} else if len(body.Email) > 0 {
		if body.Email != user.Email {
			return common.Unauthorized("邮箱错误")
		}
	} else {
		return common.BadRequest("参数错误")
	}

	user.Password = utils.MakePassword(body.NewPassword)
	err = user.Update()
	if err != nil {
		return err
	}

	tmpUser.Password = user.Password

	access, err := tmpUser.CreateJWTToken()
	if err != nil {
		return err
	}
	return c.Status(200).JSON(TokenResponse{
		Access: access,
	})
}
