package exam

import (
	"github.com/gofiber/fiber/v2"
	"github.com/opentreehole/go-common"
	"gorm.io/gorm"
	. "src/models"
)

// GetDriverPunishment @GetDriverPunishment
// @Router /api/drivers/punishments/:id [get]
// @Summary Get driver punishment by ID
// @Description Get driver punishment by ID
// @Tags Driver
// @Accept json
// @Produce json
// @Param id path int true "DriverPunishment ID"
// @Success 200 {object} DriverPunishmentResponses
// @Failure 400 {object} common.HttpError
// @Failure 401 {object} common.HttpError
// @Failure 403 {object} common.HttpError
// @Failure 404 {object} common.HttpError
// @param Authorization header string true "Bearer和token空格拼接"
func GetDriverPunishment(c *fiber.Ctx) (err error) {
	tmpUser, err := GetGeneralUser(c)
	if err != nil {
		return
	}
	punishmentID, err := c.ParamsInt("id")
	if err != nil {
		return
	}
	var driverPunishments DriverPunishments
	err = DB.Where("id = ? and user_id = ?", punishmentID, tmpUser.ID).Find(&driverPunishments).Error
	if err != nil {
		return
	}
	return c.JSON(ToModel(driverPunishments))
}

// AddDriverPunishment @AddDriverPunishment
// @Router /api/drivers/punishments/ [post]
// @Summary Add driver punishment
// @Description Add driver punishment
// @Tags Driver
// @Accept json
// @Produce json
// @Param json body AddPunishmentRequest true "json"
// @Success 200 {object} models.DriverPunishments
// @Failure 400 {object} common.HttpError
// @Failure 401 {object} common.HttpError
// @Failure 403 {object} common.HttpError
// @Failure 404 {object} common.HttpError
// @param Authorization header string true "Bearer和token空格拼接"
func AddDriverPunishment(c *fiber.Ctx) (err error) {
	tmpUser, err := GetGeneralUser(c)
	if err != nil {
		return
	}
	var addPunishmentRequest AddPunishmentRequest
	err = common.ValidateBody(c, &addPunishmentRequest)
	if err != nil {
		return
	}
	punishment := DriverPunishment{
		UserID:         tmpUser.ID,
		Reason:         addPunishmentRequest.Reason,
		Score:          addPunishmentRequest.Score,
		PunishmentType: addPunishmentRequest.PunishmentType,
	}
	err = DB.Transaction(func(tx *gorm.DB) (err error) {
		var user User
		err = tx.Clauses(LockingClause).Take(&user, tmpUser.ID).Error
		if err != nil {
			return
		}
		if !user.ValidateDriver() {
			return common.Forbidden("You don't have the driver's license")
		}
		if user.Point <= addPunishmentRequest.Score {
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

// ListDriverPunishments @ListDriverPunishments
// @Router /api/drivers/punishments/ [get]
// @Summary List driver punishments
// @Description List driver punishments
// @Tags Driver
// @Accept json
// @Produce json
// @Success 200 {object} DriverPunishmentResponses
// @Failure 400 {object} common.HttpError
// @Failure 401 {object} common.HttpError
// @Failure 403 {object} common.HttpError
// @Failure 404 {object} common.HttpError
// @param Authorization header string true "Bearer和token空格拼接"
func ListDriverPunishments(c *fiber.Ctx) (err error) {
	tmpUser, err := GetGeneralUser(c)
	if err != nil {
		return
	}
	var punishments DriverPunishments
	err = DB.Where("user_id = ?", tmpUser.ID).Find(&punishments).Error
	if err != nil {
		return
	}

	return c.JSON(ToModel(punishments))
}

func ToModel(punishments DriverPunishments) (driverPunishmentResponses DriverPunishmentResponses) {
	for _, punishment := range punishments {
		punishmentType := PunishmentTypeMap[punishment.PunishmentType]
		if punishmentType == "" {
			punishmentType = "未知"
		}
		driverPunishmentResponses = append(driverPunishmentResponses, DriverPunishmentResponse{
			ID:             punishment.ID,
			Reason:         punishment.Reason,
			Score:          punishment.Score,
			PunishmentType: punishmentType,
		})
	}
	return
}
