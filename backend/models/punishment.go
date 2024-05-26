package models

var ExamTypeMap = map[int8]string{
	EXAM:    "考试",
	SOCIETY: "社会",
}

const (
	EXAM = iota
	SOCIETY
)

var PunishmentTypeMap = map[int8]string{
	OverSpeed:      "超速",
	IllegalParking: "违停",
	NoBelts:        "未系安全带",
	FlameOut:       "熄火",
	RedLight:       "闯红灯",
	NoLicensePlate: "无牌",
	Unknown:        "未知",
}

const (
	Unknown = iota
	OverSpeed
	IllegalParking
	NoBelts
	FlameOut
	RedLight
	NoLicensePlate
)

type ExamPunishment struct {
	ID             int
	CreatedAt      MyTime
	User           *User  `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	UserID         int    `json:"user_id" gorm:"index:idx_punishment_user,priority:1"`
	Exam           *Exam  `gorm:"foreignKey:ExamID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ExamID         int    `json:"exam_id" gorm:"index:idx_punishment_exam,priority:1"`
	PunishmentType int8   `json:"punishment_type" gorm:"index:idx_punishment_type,priority:3"`
	Reason         string `json:"reason"`
	Score          int    `json:"score"`
}

type ExamPunishments []ExamPunishment

type DriverPunishment struct {
	ID             int
	CreatedAt      MyTime
	User           *User  `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	UserID         int    `json:"user_id" gorm:"index:idx_punishment_user,priority:1"`
	PunishmentType int8   `json:"punishment_type" gorm:"index:idx_punishment_type,priority:3,default:0"`
	Reason         string `json:"reason"`
	Score          int    `json:"score"`
}

type DriverPunishments []DriverPunishment
