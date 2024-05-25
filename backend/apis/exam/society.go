package exam

import (
	"github.com/gofiber/fiber/v2"
	"github.com/opentreehole/go-common"
	"gorm.io/gorm"
	. "src/models"
)

func GetSocietyPunishment(c *fiber.Ctx) (err error) {
	tmpUser, err := GetGeneralUser(c)
	if err != nil {
		return
	}
	punishmentID, err := c.ParamsInt("punishment_id")
	if err != nil {
		return
	}
	var punishments Punishments
	err = DB.Where("id = ? and user_id = ? and type = ?", punishmentID, tmpUser.ID, SOCIETY).Find(&punishments).Error
	if err != nil {
		return
	}
	return c.JSON(&punishments)
}

func AddSocietyPunishment(c *fiber.Ctx) (err error) {
	tmpUser, err := GetGeneralUser(c)
	if err != nil {
		return
	}
	var addPunishmentRequest AddPunishmentRequest
	err = common.ValidateBody(c, &addPunishmentRequest)
	if err != nil {
		return
	}
	punishment := Punishment{
		UserID: tmpUser.ID,
		Reason: addPunishmentRequest.Reason,
		Score:  addPunishmentRequest.Score,
		Type:   SOCIETY,
	}
	err = DB.Transaction(func(tx *gorm.DB) (err error) {
		var user User
		err = tx.Take(&user, tmpUser.ID).Error
		if err != nil {
			return
		}
		if !user.IsPassed {
			return common.BadRequest("You don't have the driver's license")
		}
		if user.Point < addPunishmentRequest.Score {
			user.Point = 0
			user.IsPassed = false
		} else {
			user.Point -= addPunishmentRequest.Score
		}
		err = tx.Model(&user).Select("Point", "IsPassed").Updates(&user).Error
		if err != nil {
			return
		}
		err = tx.Create(&punishment).Error
		if err != nil {
			return
		}
		return
	})
	if err != nil {
		return
	}
	return
}

func ListSocietyPunishments(c *fiber.Ctx) (err error) {
	tmpUser, err := GetGeneralUser(c)
	if err != nil {
		return
	}
	var punishments Punishments
	err = DB.Where("user_id = ? and type = ?", tmpUser.ID, SOCIETY).Find(&punishments).Error
	if err != nil {
		return
	}
	return c.JSON(punishments)
}
