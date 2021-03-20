package auth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"log"
	"os"
)

var (
	logger                = log.New(os.Stdout, "[GATEWAY][AUTH]", 0)
	jwtSecret             = []byte(os.Getenv("JWT_SECRET"))
	AccessTokenCookieName = "fuckyouwhoeverislurkingatthis"
)

type AuthorizedUser struct {
	ID string
}

func createJWTToken() (string, bool) {
	claims := jwt.StandardClaims{
		ExpiresAt: 15000,
		Issuer:    os.Getenv("JWT_ISSUER"),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token.Header["userID"] = utils.UUID()
	key, err := token.SignedString(jwtSecret)
	if err != nil {
		logger.Printf("ERROR: Couldn't generate JWT token: %s", err)
		return "", false
	}
	return key, true

}

func GetUserId(token string) string {
	user, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return ""
	}
	userID := user.Header["userID"].(string)
	return userID
}

func RequireLogin(request *fiber.Ctx) string {
	token := request.Get(AccessTokenCookieName, "")
	userId := GetUserId(token)
	if userId == "" {
		token, err := createJWTToken()
		if err == false {
			panic(nil)
		}
		request.Set(AccessTokenCookieName, token)
		userId = GetUserId(token)
	}
	return userId
}
