package models

const (
	EXAM = iota
	SOCIETY
)

const (
	Unknown = iota
	OverSpeed
	IllegalParking
	NoBelts
	FlameOut
	RedLight
	NoLicensePlate
)

type Punishment struct {
	ID             int
	CreatedAt      MyTime
	User           *User  `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	UserID         int    `json:"user_id" gorm:"index:idx_punishment_user,priority:1"`
	Exam           *Exam  `gorm:"foreignKey:ExamID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ExamID         int    `json:"exam_id" gorm:"index:idx_punishment_exam,priority:1"`
	Type           int8   `json:"type" gorm:"index:idx_type,priority:2"`
	PunishmentType int8   `json:"punishment_type" gorm:"index:idx_punishment_type,priority:3"`
	Reason         string `json:"reason"`
	Score          int    `json:"score"`
}

type Punishments []Punishment
