package models

import (
	"errors"
	"gorm.io/gorm"
)

type User struct {
	ID        int    `json:"id" gorm:"primaryKey"`
	Username  string `json:"username" gorm:"uniqueIndex;not null;size:64"`
	Password  string `json:"password" gorm:"not null;size:64"`
	Email     string `json:"email" gorm:"size:64"`
	Age       int8   `json:"age" gorm:"not null"`
	Gender    string `json:"gender" gorm:"not null;type:char(8)"`
	Phone     string `json:"phone" gorm:"not null;type:char(13)"`
	CreatedAt MyTime `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt MyTime `json:"updated_at" gorm:"autoUpdateTime"`
	IsPassed  bool   `json:"is_passed" gorm:"default:false"`
	Point     int    `json:"point" gorm:"default:0"`
}

//func (user User) MarshalJSON() ([]byte, error) {
//	var userMap = map[string]interface{}{
//		"id":       user.ID,
//		"username": user.Username,
//		"email":    user.Email,
//		"age":      user.Age,
//		"gender":   user.Gender,
//		"phone":    user.Phone,
//	}
//	return json.Marshal(userMap)
//}

func GetUsers() (users []User, err error) {
	err = DB.Transaction(func(tx *gorm.DB) error {
		return tx.Find(&users).Error
	})
	return
}

func GetUserByID(userID int) (user *User, err error) {
	err = DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Omit("Password").Take(&user, userID).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// insert user if not found
				user.ID = userID
				err = tx.Create(&user).Error
				if err != nil {
					return err
				}
			} else {
				return err
			}
		}
		return nil
	})
	return
}

func GetUserByUsername(username string) (user *User, err error) {
	//Omit("Password").
	err = DB.Where("username = ?", username).Take(&user).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, nil
	}
	return
}

func (user *User) Create() error {
	return DB.Create(&user).Error
}

func (user *User) Update() error {
	return DB.Transaction(func(tx *gorm.DB) error {
		result := tx.Updates(&user)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errors.New("user not found or no change")
			//	todo: return nil or error
		}
		return nil
	})
}

func DeleteUserByID(id int) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		result := tx.Delete(&User{}, id)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errors.New("user not found")
		}
		return nil
	})
}

func ValidateDriver(id int) (bool, error) {
	var user User
	err := DB.Take(&user, id).Error
	if err != nil {
		return false, err
	}
	return user.Age >= 18 && user.IsPassed, nil
}

func (user *User) AddPunishment(score int) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		var user User
		err := tx.Take(&user, user.ID).Error
		if err != nil {
			return err
		}
		if user.Point <= score {
			user.Point = 0
			user.IsPassed = false
		} else {
			user.Point -= score
		}
		return tx.Model(&user).Select("Point", "IsPassed").Updates(&user).Error
	})
}

func (user *User) ValidateDriver() bool {
	return user.Age >= 18 && user.IsPassed
}
