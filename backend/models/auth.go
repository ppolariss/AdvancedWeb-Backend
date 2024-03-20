package models

import (
	"github.com/gofiber/fiber/v2"
	"github.com/opentreehole/go-common"
	"gorm.io/gorm"
)

type LoginUser struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// GetGeneralUser
// return user from fiber.Ctx or jwt
func GetGeneralUser(c *fiber.Ctx) (user *LoginUser, err error) {
	if c.Locals("user") != nil {
		user = c.Locals("user").(*LoginUser)
		return
	}

	// get id and user_type from jwt
	token := common.GetJWTToken(c)
	if token == "" {
		return nil, common.Unauthorized("Unauthorized")
	}
	err = common.ParseJWTToken(token, &user)
	if err != nil {
		return nil, common.Unauthorized("Unauthorized")
	}

	// load user from database in transaction
	err = user.CheckUserID()
	if err != nil {
		return
	}

	// save user in c.Locals
	c.Locals("user", user)
	return
}

func (user *LoginUser) CheckUserID() error {
	return DB.Transaction(func(tx *gorm.DB) error {
		var newUser = User{ID: user.ID}
		err := tx.Take(&newUser).Error
		if err != nil {
			return err
		}
		user.Username = newUser.Username

		return nil
	})
}
