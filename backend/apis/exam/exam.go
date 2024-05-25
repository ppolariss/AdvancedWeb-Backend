package exam

import (
	"github.com/gofiber/fiber/v2"
	"github.com/opentreehole/go-common"
	"gorm.io/gorm"
	. "src/models"
	"time"
)

func ListExams(ctx *fiber.Ctx) (err error) {
	tmpUser, err := GetGeneralUser(ctx)
	if err != nil {
		return
	}
	var exams Exams
	err = DB.Where("user_id = ?", tmpUser.ID).Find(&exams).Error
	if err != nil {
		return
	}
	return ctx.JSON(exams)
}

func GetExam(c *fiber.Ctx) (err error) {
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
	return c.JSON(exam)
}

func AddExam(ctx *fiber.Ctx) (err error) {
	tmpUser, err := GetGeneralUser(ctx)
	if err != nil {
		return
	}
	var addExamRequest AddExamRequest
	if err = ctx.BodyParser(&addExamRequest); err != nil {
		return common.BadRequest("Invalid request body")
	}
	var exam = Exam{
		UserID:      tmpUser.ID,
		Title:       addExamRequest.Title,
		Description: addExamRequest.Description,
		StartTime:   addExamRequest.StartTime,
		EndTime:     addExamRequest.EndTime,
		Duration:    addExamRequest.EndTime.Time.Sub(addExamRequest.StartTime.Time),
		Score:       addExamRequest.Score,
	}
	return DB.Create(&exam).Error
}

func StartExam(ctx *fiber.Ctx) (err error) {
	tmpUser, err := GetGeneralUser(ctx)
	if err != nil {
		return
	}
	return DB.Create(&Exam{
		UserID: tmpUser.ID,
		StartTime: MyTime{
			Time: time.Now(),
		},
	}).Error
}

func EndExam(ctx *fiber.Ctx) (err error) {
	tmpUser, err := GetGeneralUser(ctx)
	if err != nil {
		return
	}
	var endExamRequest EndExamRequest
	if err = ctx.BodyParser(&endExamRequest); err != nil {
		return common.BadRequest("Invalid request body")
	}
	var exam Exam
	err = DB.Take(&exam, endExamRequest.ID).Error
	if err != nil {
		return
	}
	if exam.UserID != tmpUser.ID {
		return common.Forbidden("You are not the owner of this exam")
	}
	exam.EndTime = MyTime{Time: time.Now()}
	exam.Duration = exam.EndTime.Time.Sub(exam.StartTime.Time)
	err = DB.Transaction(func(tx *gorm.DB) (err error) {
		err = DB.Model(&exam).Select("EndTime", "Duration").UpdateColumns(&exam).Error
		if err != nil {
			return
		}
		if exam.Score != 100 {
			return
		}
		var user User
		err = DB.Take(&user, tmpUser.ID).Error
		if err != nil {
			return
		}
		if !user.IsPassed {
			return
		}
		user.IsPassed = true
		user.Point = 12
		return DB.Model(&user).Select("IsPassed", "Point").UpdateColumns(&user).Error
	})
	var endExamResponse = EndExamResponse{
		Score: exam.Score,
	}
	return ctx.JSON(endExamResponse)
}

func DeleteExam(ctx *fiber.Ctx) error {
	tmpUser, err := GetGeneralUser(ctx)
	if err != nil {
		return err
	}
	examID, err := ctx.ParamsInt("id")
	if err != nil {
		return err
	}
	var exam Exam
	err = DB.Take(&exam, examID).Error
	if err != nil {
		return err
	}
	if exam.UserID != tmpUser.ID {
		return common.Forbidden("You are not the owner of this exam")
	}
	return DB.Delete(&exam).Error
}

func ModifyExam(c *fiber.Ctx) (err error) {
	tmpUser, err := GetGeneralUser(c)
	if err != nil {
		return
	}
	var modifyExamRequest ModifyExamRequest
	if err = c.BodyParser(&modifyExamRequest); err != nil {
		return common.BadRequest("Invalid request body")
	}
	var exam Exam
	err = DB.Take(&exam, modifyExamRequest.ID).Error
	if err != nil {
		return
	}
	if exam.UserID != tmpUser.ID {
		return common.Forbidden("You are not the owner of this exam")
	}
	if modifyExamRequest.Title != "" {
		exam.Title = modifyExamRequest.Title
	}
	if modifyExamRequest.Description != "" {
		exam.Description = modifyExamRequest.Description
	}
	return DB.Model(&exam).Select("Title", "Description").Updates(&exam).Error
}