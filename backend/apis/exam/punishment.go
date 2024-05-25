package exam

import (
	"github.com/gofiber/fiber/v2"
	"github.com/opentreehole/go-common"
	"gorm.io/gorm"
	. "src/models"
	"time"
)

func AddPunishment(c *fiber.Ctx) (err error) {
	tmpUser, err := GetGeneralUser(c)
	if err != nil {
		return
	}
	var addPunishmentRequest AddPunishmentRequest
	if err = c.BodyParser(&addPunishmentRequest); err != nil {
		return common.BadRequest("Invalid request body")
	}
	examID, err := c.ParamsInt("id")
	var exam Exam
	err = DB.Take(&exam, examID).Error
	if err != nil {
		return
	}
	if exam.UserID != tmpUser.ID {
		return common.Forbidden("You are not the owner of this exam")
	}
	var punishment = Punishment{
		Reason:         addPunishmentRequest.Reason,
		Score:          addPunishmentRequest.Score,
		PunishmentType: addPunishmentRequest.PunishmentType,
		CreatedAt:      MyTime{Time: time.Now()},
		UserID:         tmpUser.ID,
		ExamID:         examID,
		Type:           EXAM,
	}
	return DB.Transaction(func(tx *gorm.DB) (err error) {
		err = tx.Model(&exam).UpdateColumn("Score", exam.Score-addPunishmentRequest.Score).Error
		if err != nil {
			return
		}
		return tx.Create(&punishment).Error
	})
	//if
	//exam.Punishments = append(exam.Punishments, Punishment{
	//	Reason: addPunishmentRequest.Reason,
	//	Score:  addPunishmentRequest.Score,
	//})
	// DB.Save(&exam).Error
}

func ListPunishments(c *fiber.Ctx) (err error) {
	tmpUser, err := GetGeneralUser(c)
	if err != nil {
		return
	}
	examID, err := c.ParamsInt("id")
	if err != nil {
		return
	}
	var exam Exam
	err = DB.Take(&exam, examID).Error
	if err != nil {
		return
	}
	if exam.UserID != tmpUser.ID {
		return common.Forbidden("You are not the owner of this exam")
	}
	var punishments Punishments
	err = DB.Where("exam_id = ? and type = ?", exam.ID, EXAM).Find(&punishments).Error
	return c.JSON(exam.Punishments)
}

func GetPunishment(c *fiber.Ctx) (err error) {
	tmpUser, err := GetGeneralUser(c)
	if err != nil {
		return
	}
	punishmentID, err := c.ParamsInt("id")
	if err != nil {
		return
	}
	var punishment Punishment
	err = DB.Take(&punishment, punishmentID).Error
	if err != nil {
		return
	}
	if punishment.UserID != tmpUser.ID {
		return common.Forbidden("You are not the owner of this punishment")
	}
	return c.JSON(punishment)
}
