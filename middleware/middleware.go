package middleware

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"github.com/oscar-mgh/gin_bun/db"
	"github.com/oscar-mgh/gin_bun/model"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
)

type HttpResponse struct {
	Message     string
	Status      int
	Description string
}

func ValidarJWT(header string) int {
	err := godotenv.Load()
	if err != nil {
		return 0
	}

	jwtSecret := []byte(os.Getenv("JWT_SECRET"))

	if len(header) == 0 {
		return 0
	}
	splitBearer := strings.Split(header, " ")
	if len(splitBearer) != 2 {
		return 0
	}
	splitToken := strings.Split(splitBearer[1], ".")
	if len(splitToken) != 3 {
		return 0
	}
	tk := strings.TrimSpace(splitBearer[1])
	token, err := jwt.Parse(tk, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: ")
		}

		return jwtSecret, nil
	})
	if err != nil {
		return 0
	}
	if err != nil {
		return 0
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		user := model.UserModel{}

		db := bun.NewDB(db.Connect(), mysqldialect.New())
		err := db.NewSelect().Model(&user).Where("email = ?", claims["email"]).Scan(context.TODO())

		if err != nil {
			return 0
		}
		return 1

	} else {
		return 0
	}
}
