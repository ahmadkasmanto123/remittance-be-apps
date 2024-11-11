package config

import (
	"errors"
	"fmt"
	"love-remittance-be-apps/lib/model"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	jwtware "github.com/gofiber/contrib/jwt"
)

var sckey = os.Getenv("SECRET_KEY")
var secretKey = []byte(sckey)

func CreateToken(username string, phone string, prefix string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"customer_name": username,
			"phone":         phone,
			"phone_prefix":  prefix,
			"exp":           time.Now().Add(time.Minute * 30).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}

func JwtConfig() func(*fiber.Ctx) error {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(sckey)},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			var rc = "401"
			if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
				rc = "401NOTV"
			} else if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
				rc = "401SignFail"
			} else {
				rc = "401Else"
				fmt.Print(c.Get("Authorization"))
			}
			return c.Status(fiber.StatusUnauthorized).JSON(map[string]interface{}{
				"rc": rc,
				"error": model.ErrorData{
					Description: err.Error(),
				},
				"message": "Unauthorized Customer",
			})
		},
	})
}
