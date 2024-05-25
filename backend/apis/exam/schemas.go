package exam

import (
	. "src/models"
)

type StartExamRequest struct {
}

type StartExamResponse struct {
	ID int `json:"id"`
}

type GetScoreRequest struct {
	ID int `json:"id"`
}

type GetScoreResponse struct {
	Score int `json:"score"`
}

type EndExamRequest struct {
	ID int `json:"id"`
}

type EndExamResponse struct {
	Score int `json:"score"`
}

type AddExamRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	StartTime   MyTime `json:"start_time"`
	EndTime     MyTime `json:"end_time"`
	Score       int    `json:"score" validate:"required"`
}

type AddPunishmentRequest struct {
	PunishmentType int8   `json:"punishment_type"`
	Reason         string `json:"reason"`
	Score          int    `json:"score"`
}

type ModifyExamRequest struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
