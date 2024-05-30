package models

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/opentreehole/go-common"
	"gorm.io/gorm"
	"src/config"
	"src/utils"
	"strconv"
	"strings"
	"time"
)

type LoginUser struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// GetGeneralUser
// return user from fiber.Ctx or jwt
func GetGeneralUser(c *fiber.Ctx) (user *LoginUser, err error) {
	if c.Locals("user") != nil {
		user = c.Locals("user").(*LoginUser)
		return
	}

	// get id and user_type from jwt
	token := common.GetJWTToken(c)
	if token == "" {
		return nil, common.Unauthorized("Unauthorized")
	}
	err = common.ParseJWTToken(token, &user)
	if err != nil {
		return nil, common.Unauthorized("Unauthorized")
	}
	// err = jwt.ParseWithClaims(token, &user, func(token *jwt.Token) (interface{}, error) {
	// 	return []byte(config.GetJwtSecret()), nil
	// })
	var userJwtSecret UserJwtSecret
	err = DB.Take(&userJwtSecret, user.ID).Error
	if err != nil {
		return nil, common.Unauthorized("Unauthorized")
	}
	_, err = CheckJWTToken(token, userJwtSecret.Secret)
	if err != nil {
		return nil, common.Unauthorized("Unauthorized")
	}

	// load user from database in transaction
	err = user.CheckUserID()
	if err != nil {
		return
	}

	// save user in c.Locals
	c.Locals("user", user)
	return
}

func (user *LoginUser) CheckUserID() error {
	return DB.Transaction(func(tx *gorm.DB) error {
		var newUser = User{ID: user.ID}
		err := tx.Take(&newUser).Error
		if err != nil {
			return err
		}
		// TODO
		//if user.Username != newUser.Username {
		//	return common.Unauthorized("Unauthorized")
		//}
		user.Username = newUser.Username

		return nil
	})
}

func GetUserCacheKey(userID int) string {
	return "moss_user:" + strconv.Itoa(userID)
}

const UserCacheExpire = 48 * time.Hour

func GetUserID(c *fiber.Ctx) (int, error) {

	id, err := strconv.Atoi(c.Get("X-Consumer-Username"))
	if err != nil {
		return 0, common.Unauthorized("Unauthorized")
	}

	return id, nil
}

// LoadUserByIDFromCache return value `err` is directly from DB.Take()
func LoadUserByIDFromCache(userID int, userPtr *User) error {
	cacheKey := GetUserCacheKey(userID)
	if config.GetCache(cacheKey, userPtr) != nil {
		err := DB.Take(userPtr, userID).Error
		if err != nil {
			return err
		}
		// err has been printed in SetCache
		_ = config.SetCache(cacheKey, *userPtr, UserCacheExpire)
	}
	return nil
}

func DeleteUserCacheByID(userID int) {
	cacheKey := GetUserCacheKey(userID)
	err := config.DeleteCache(cacheKey)
	if err != nil {
		utils.Logger.Error("err in DeleteUserCacheByID: ")
		//	TODO
	}
}

func LoadUserByID(userID int) (*User, error) {
	var user User
	err := LoadUserByIDFromCache(userID, &user)
	if err != nil { // something wrong in DB.Take() in LoadUserByIDFromCache()
		DeleteUserCacheByID(userID)
		return nil, err
	}
	updated := false

	if updated {
		DB.Model(&user).Select("ModelID", "PluginConfig").Updates(&user)
		err = config.SetCache(GetUserCacheKey(userID), user, UserCacheExpire)
	}
	return &user, err
}

func LoadUser(c *fiber.Ctx) (*User, error) {
	userID, err := GetUserID(c)
	if err != nil {
		return nil, err
	}
	return LoadUserByID(userID)
}

func GetUserIDFromWs(c *websocket.Conn) (int, error) {
	// get cookie named access or query jwt
	token := c.Query("jwt")
	if token == "" {
		token = c.Cookies("access")
		if token == "" {
			return 0, common.Unauthorized()
		}
	}
	// get data
	data, err := parseJWT(token, false)
	if err != nil {
		return 0, err
	}
	id, ok := data["id"] // get id
	if !ok {
		return 0, common.Unauthorized()
	}
	return int(id.(float64)), nil
}

func LoadUserFromWs(c *websocket.Conn) (*User, error) {
	userID, err := GetUserIDFromWs(c)
	if err != nil {
		return nil, err
	}
	return LoadUserByID(userID)
}

// parseJWT extracts and parse token
func parseJWT(token string, bearer bool) (Map, error) {
	if len(token) < 7 {
		return nil, errors.New("bearer token required")
	}

	if bearer {
		token = token[7:]
	}

	payloads := strings.SplitN(token[7:], ".", 3) // extract "Bearer "
	if len(payloads) < 3 {
		return nil, errors.New("jwt token required")
	}

	// jwt encoding ignores padding, so RawStdEncoding should be used instead of StdEncoding
	payloadBytes, err := base64.RawURLEncoding.DecodeString(payloads[1]) // the middle one is payload
	if err != nil {
		return nil, err
	}

	var value Map
	err = json.Unmarshal(payloadBytes, &value)
	return value, err
}

func GetUserByRefreshToken(c *fiber.Ctx) (*User, error) {
	// get id
	userID, err := GetUserID(c)
	if err != nil {
		return nil, err
	}

	tokenString := c.Get("Authorization")
	if tokenString == "" { // token can be in either header or cookie
		tokenString = c.Cookies("refresh")
	}

	payload, err := parseJWT(tokenString, true)
	if err != nil {
		return nil, err
	}

	if tokenType, ok := payload["type"]; !ok || tokenType != "refresh" {
		return nil, common.Unauthorized("refresh token invalid")
	}

	var user User
	err = LoadUserByIDFromCache(userID, &user)
	return &user, err
}
