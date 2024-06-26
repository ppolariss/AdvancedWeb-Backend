package exam

import (
	"fmt"
	. "src/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/opentreehole/go-common"
	"gorm.io/gorm"
)

// ListExams @ListExams
// @Router /api/exams [get]
// @Summary list my exams
// @Description list my exams
// @Tags Exam
// @Accept json
// @Produce json
// @Success 200 {object} Exams
// @Failure 400 {object} common.HttpError
// @Failure 401 {object} common.HttpError
// @param Authorization header string true "Bearer和token空格拼接"
func ListExams(ctx *fiber.Ctx) (err error) {
	tmpUser, err := GetGeneralUser(ctx)
	if err != nil {
		return
	}
	var exams Exams
	err = DB.Where("user_id = ? and exam_type = 'exam'", tmpUser.ID).Order("id desc").Find(&exams).Error
	if err != nil {
		return
	}
	return ctx.JSON(exams)
}

// ListDriverExams @ListDriverExams
// @Router /api/drivers [get]
// @Summary list my driver exams
// @Description list my driver exams
// @Tags Driver
// @Accept json
// @Produce json
// @Success 200 {object} Exams
// @Failure 400 {object} common.HttpError
// @Failure 401 {object} common.HttpError
// @param Authorization header string true "Bearer和token空格拼接"
func ListDriverExams(ctx *fiber.Ctx) (err error) {
	tmpUser, err := GetGeneralUser(ctx)
	if err != nil {
		return
	}
	var exams Exams
	err = DB.Where("user_id = ? and exam_type = 'driver'", tmpUser.ID).Order("id desc").Find(&exams).Error
	if err != nil {
		return
	}
	return ctx.JSON(exams)
}

// GetExam @GetExam
// @Router /api/exams/{id} [get]
// @Summary Get exam by ID
// @Description Get exam by ID
// @Tags Exam
// @Accept json
// @Produce json
// @Param id path int true "Exam ID"
// @Success 200 {object} Exam
// @Failure 400 {object} common.HttpError
// @Failure 401 {object} common.HttpError
// @Failure 403 {object} common.HttpError
// @param Authorization header string true "Bearer和token空格拼接"
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

// AddExam @AddExam
// @Router /api/exams/add [post]
// @Summary Add exam once
// @Description Add exam once
// @Tags Exam
// @Accept json
// @Produce json
// @Param json body AddExamRequest true "json"
// @Success 200 {object} models.Exam
// @Failure 400 {object} common.HttpError
// @Failure 401 {object} common.HttpError
// @Param Authorization header string true "Bearer和token空格拼接"
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
		StartTime:   &addExamRequest.StartTime,
		EndTime:     &addExamRequest.EndTime,
		Duration:    addExamRequest.EndTime.Time.Sub(addExamRequest.StartTime.Time),
		Score:       addExamRequest.Score,
	}
	return DB.Create(&exam).Error
}

// StartExam @StartExam
// @Router /api/exams/start [post]
// @Summary Start exam
// @Description Start exam
// @Tags Exam
// @Accept json
// @Produce json
// @Success 200 {object} StartExamResponse
// @Failure 400 {object} common.HttpError
// @Failure 401 {object} common.HttpError
// @Param Authorization header string true "Bearer和token空格拼接"
func StartExam(ctx *fiber.Ctx) (err error) {
	tmpUser, err := GetGeneralUser(ctx)
	if err != nil {
		return
	}
	var startExamRequest StartExamRequest
	if err = ctx.BodyParser(&startExamRequest); err != nil {
		return common.BadRequest("Invalid request body")
	}
	examType := startExamRequest.ExamType
	if examType == "" {
		examType = "exam"
	}
	exam := Exam{
		UserID: tmpUser.ID,
		StartTime: &MyTime{
			Time: time.Now(),
		},
		ExamType: examType,
	}
	err = DB.Omit("EndTime").Create(&exam).Error
	if err != nil {
		return
	}
	return ctx.JSON(StartExamResponse{
		ID: exam.ID,
	})
}

// EndExam @EndExam
// @Router /api/exams/end [post]
// @Summary End exam
// @Description End exam
// @Tags Exam
// @Accept json
// @Produce json
// @Param json body EndExamRequest true "json"
// @Success 200 {object} EndExamResponse
// @Failure 400 {object} common.HttpError
// @Failure 401 {object} common.HttpError
// @Param Authorization header string true "Bearer和token空格拼接"
func EndExam(ctx *fiber.Ctx) (err error) {
	tmpUser, err := GetGeneralUser(ctx)
	if err != nil {
		return
	}
	isDriver := false
	info := "Oops, something went wrong"
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
	if exam.IsFinished() {
		return common.Forbidden("This exam has already ended")
	}
	exam.Normal = endExamRequest.Normal
	exam.EndTime = &MyTime{Time: time.Now()}
	exam.Duration = exam.EndTime.Time.Sub(exam.StartTime.Time)
	var endExamResponse EndExamResponse
	err = DB.Transaction(func(tx *gorm.DB) (err error) {
		err = tx.Model(&exam).Select("EndTime", "Duration", "Normal").UpdateColumns(&exam).Error
		if err != nil {
			return
		}
		var user User
		err = tx.Take(&user, tmpUser.ID).Error
		if err != nil {
			return
		}
		if exam.ExamType != "exam" {
			endExamResponse = EndExamResponse{
				Score:    user.Point,
				Info:     "Cool! You've had a great road trip",
				IsPassed: user.IsPassed,
				IsDriver: user.IsPassed,
			}
			return
		}
		if exam.Score != 100 {
			info = "Sorry, you didn't pass the exam"
			return
		}
		if !exam.Normal {
			info = "Sorry, you quit the exam in advance"
			return
		}

		if user.IsPassed {
			info = "Congratulations! You have passed the exam"
			isDriver = true
			return
		}
		if user.Age < 18 {
			info = "Congratulations! You have passed the exam. However, you are not eligible for a driver's license because you are " + fmt.Sprintf("%d", user.Age) + " years old"
			return
		}
		info = "Congratulations! You have passed the exam and obtained the driver's license"
		isDriver = true
		user.IsPassed = true
		user.Point = 12
		return tx.Model(&user).Select("IsPassed", "Point").UpdateColumns(&user).Error
	})
	if err != nil {
		return
	}
	if exam.ExamType == "exam" {
		endExamResponse = EndExamResponse{
			Score:    exam.Score,
			IsPassed: exam.Score == 100,
			IsDriver: isDriver,
			Info:     info,
		}
	}
	return ctx.JSON(endExamResponse)
}

// DeleteExam @DeleteExam
// @Router /api/exams/{id} [delete]
// @Summary Delete exam by ID
// @Description Delete exam by ID
// @Tags Exam
// @Accept json
// @Produce json
// @Param id path int true "Exam ID"
// @Success 200
// @Failure 400 {object} common.HttpError
// @Failure 401 {object} common.HttpError
// @Failure 403 {object} common.HttpError
// @param Authorization header string true "Bearer和token空格拼接"
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

// ModifyExam @ModifyExam
// @Router /api/exams/ [put]
// @Summary Modify exam by ID
// @Description Modify exam by ID
// @Tags Exam
// @Accept json
// @Produce json
// @Param id path int true "Exam ID"
// @Param json body ModifyExamRequest true "json"
// @Success 200
// @Failure 400 {object} common.HttpError
// @Failure 401 {object} common.HttpError
// @Failure 403 {object} common.HttpError
// @param Authorization header string true "Bearer和token空格拼接"
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
