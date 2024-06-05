package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/opentreehole/go-common"
	. "src/models"
	"src/utils"
	"time"
)

// Login godoc
// @Router /api/login [post]
// @Summary Login
// @Description Login
// @Tags Auth
// @Accept json
// @Produce json
// @Param json body LoginRequest true "json"
// @Success 200 {object} TokenResponse
// @Failure 400 {object} common.HttpError
// @Failure 401 {object} common.HttpError
// @Failure 403 {object} common.HttpError
// @Failure 500 {object} common.HttpError
func Login(c *fiber.Ctx) error {
	var body LoginRequest
	err := common.ValidateBody(c, &body)
	if err != nil {
		return err
	}

	var tmpUser LoginUser
	user, err := GetUserByUsername(body.Username)
	if err != nil {
		return common.InternalServerError()
	}
	if user == nil {
		return common.Unauthorized("账号未注册")
	}

	tmpUser.ID = user.ID
	tmpUser.Username = user.Username
	tmpUser.Password = user.Password

	ok := utils.CheckPassword(body.Password, tmpUser.Password)
	if !ok {
		//return fiber.NewError(fiber.StatusUnauthorized, "密码错误")
		return common.Unauthorized("密码错误")
	}

	access, err := tmpUser.CreateJWTToken()
	if err != nil {
		return err
	}

	return c.Status(200).JSON(TokenResponse{
		Access: access,
		UserID: user.ID,
		//Message: "登录成功",
	})
}

func Logout(c *fiber.Ctx) error {
	return nil
}

// Register godoc
// @Router /api/register [post]
// @Summary Register
// @Description Register
// @Tags Auth
// @Accept json
// @Produce json
// @Param json body RegisterRequest true "json"
// @Success 200
// @Failure 400 {object} common.HttpError
// @Failure 500 {object} common.HttpError
func Register(c *fiber.Ctx) error {
	var body RegisterRequest
	err := common.ValidateBody(c, &body)
	if err != nil {
		return err
	}

	user, err := GetUserByUsername(body.Username)
	if err != nil {
		return common.InternalServerError()
	}
	if user != nil {
		return common.BadRequest("用户名已存在")
	}
	user = &User{
		Username:  body.Username,
		Password:  utils.MakePassword(body.Password),
		Gender:    body.Gender,
		Email:     body.Email,
		Phone:     body.Phone,
		Age:       body.Age,
		CreatedAt: MyTime{Time: time.Now()},
	}
	return user.Create()
}

func ChangePassword(c *fiber.Ctx) error {
	return nil
}

func ForgetPassword(c *fiber.Ctx) error {
	return nil
}
