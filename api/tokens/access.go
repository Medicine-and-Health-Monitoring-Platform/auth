package tokens

import (
	"Auth/config"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
)

func GenerateAccessToken(id, email, role string) (string, error) {
	token := *jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = id
	claims["email"] = email
	claims["role"] = role
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(1 * time.Hour).Unix()

	newToken, err := token.SignedString([]byte(config.Load().ACCESS_TOKEN))

	if err != nil {
		log.Println(err)
		return "", errors.Wrap(err, "failed to generate access token")
	}

	return newToken, nil
}
