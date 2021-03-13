package main

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"os"
)

var (
	jwtToken              = []byte(os.Getenv("JWT_TOKEN"))
	accessTokenCookieName = "fuckyouwhoeverislurkingatthis"
)

type JWTClaims struct {
	UserID string `json:"userID"`
	jwt.StandardClaims
}
type LoginInfo struct {
	userID string
	token  string
}

func createJWTToken() (string, bool) {
	claims := jwt.StandardClaims{
		ExpiresAt: 15000,
		Issuer:    "UnblockHub",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token.Header["userID"] = utils.UUID()
	key, err := token.SignedString(jwtToken)
	if err != nil {
		logger.Printf("ERROR: Couldn't generate JWT token: %s", err)
		return "", false
	}
	return key, true

}

func getUserId(token string) string {
	user, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return jwtToken, nil
	})
	if err != nil {
		return ""
	}
	userID := user.Header["userID"].(string)
	return userID
}

func requireLogin(request *fiber.Ctx) string {
	token := request.Cookies(accessTokenCookieName, "")
	userId := getUserId(token)
	if userId == "" {
		token, err := createJWTToken()
		if err == false {
			panic(nil)
		}
		request.Cookie(&fiber.Cookie{
			Name:  accessTokenCookieName,
			Value: token,
		})
		userId = getUserId(token)
	}
	return userId
}
