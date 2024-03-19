package models

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/thanhpk/randstr"
	"gorm.io/gorm"
	"time"
)

func main() {
	// 生成 JWT
	token := jwt.New(jwt.SigningMethodHS256)

	// 设置载荷（Payload）
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = "john.doe"
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	// 生成签名
	secret := []byte("your-secret-key")
	signedToken, err := token.SignedString(secret)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("JWT:", signedToken)

	// 解析 JWT
	parsedToken, err := jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		username := claims["username"].(string)
		expiration := claims["exp"].(float64)

		fmt.Println("Username:", username)
		fmt.Println("Expiration:", time.Unix(int64(expiration), 0))
	} else {
		fmt.Println("Invalid token:", err)
	}
}

type UserJwtSecret struct {
	ID     int    `json:"id" gorm:"primaryKey"`
	Secret string `json:"secret" gorm:"size:256"`
}

type UserClaims struct {
	jwt.RegisteredClaims
	ID       int    `json:"id"`
	UserType string `json:"user_type"`
	Type     string `json:"type"`
	//UserID               int       `json:"user_id"`
	//UID                  int       `json:"uid"`
	//Nickname             string    `json:"nickname"`
	//JoinedTime           time.Time `json:"joined_time"`
	//HasAnsweredQuestions bool      `json:"has_answered_questions"`
}

const (
	JWTTypeAccess = "access"
	//JWTTypeRefresh = "refresh"
)

//func DeleteUserJWTSecret(userID int) error {
//	return DB.Delete(&UserJwtSecret{}, userID).Error
//}

//func CheckJWTToken(token string, secret string) (claims *UserClaims, err error) {
//	claims = new(UserClaims)
//	_, err = jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
//		return []byte(secret), nil
//	})
//	return
//}

func (user *LoginUser) CreateJWTToken() (accessToken string, err error) {
	// get jwt key and secret
	var key, secret string
	// no gateway, store jwt secret in database
	var userJwtSecret UserJwtSecret
	err = DB.Take(&userJwtSecret, user.ID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			userJwtSecret = UserJwtSecret{
				ID:     user.ID,
				Secret: randstr.Base62(32),
			}
			err = DB.Create(&userJwtSecret).Error
			if err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	}

	key = fmt.Sprintf("user_%d", user.ID)
	secret = userJwtSecret.Secret

	// create JWT token
	claim := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    key,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(60 * time.Minute)), // // 60 minutes
		},
		ID: user.ID,
	}

	// access payload
	accessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claim).SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	// refresh payload
	claim.Type = JWTTypeAccess
	//claim.ExpiresAt = jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)) // 30 days
	//refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claim).SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return
}
