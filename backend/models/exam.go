package models

import (
	"gorm.io/gorm"
	"time"
)

type Exam struct {
	//IsPublic    bool   `json:"is_public"`
	ID          int            `json:"id" gorm:"primaryKey"`
	Title       string         `json:"title" gorm:"default:'exam'"`
	Description string         `json:"description" gorm:"default:'description'"`
	StartTime   *MyTime        `json:"start_time"`
	EndTime     *MyTime        `json:"end_time"`
	Duration    time.Duration  `json:"duration"`
	User        *User          `validate:"omitempty" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	UserID      int            `json:"user_id" gorm:"index:idx_exam_user,priority:1"`
	Score       int            `json:"score" gorm:"default:100"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" validate:"omitempty"`
	//Punishments *Punishments   `json:"punishments"`
}

type Exams []Exam

func (exam Exam) IsFinished() bool {
	//return exam.EndTime.Time.Before(time.Now())
	return exam.EndTime != nil
}
