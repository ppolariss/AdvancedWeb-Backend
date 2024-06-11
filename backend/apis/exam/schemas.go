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
	Score    int    `json:"score"`
	IsPassed bool   `json:"is_passed"`
	IsDriver bool   `json:"is_driver"`
	Info     string `json:"info"`
}

type AddExamRequest struct {
	Title       string `json:"title" validate:"omitempty,max=255"`
	Description string `json:"description" validate:"omitempty,max=255"`
	StartTime   MyTime `json:"start_time" validate:"required"`
	EndTime     MyTime `json:"end_time" validate:"required"`
	Score       int    `json:"score" validate:"required"`
}

type AddPunishmentRequest struct {
	PunishmentType int8   `json:"punishment_type"`
	Reason         string `json:"reason"`
	Score          int    `json:"score" validate:"required,min=0"`
}

type ModifyExamRequest struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type DriverPunishmentResponse struct {
	ID             int    `json:"id"`
	CreatedAt      MyTime `json:"created_at"`
	PunishmentType string `json:"punishment_type"`
	Reason         string `json:"reason"`
	Score          int    `json:"score"`
}

type DriverPunishmentResponses []DriverPunishmentResponse
