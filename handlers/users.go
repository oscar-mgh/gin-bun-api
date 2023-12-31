package handlers

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/oscar-mgh/gin_bun/db"
	"github.com/oscar-mgh/gin_bun/dto"
	"github.com/oscar-mgh/gin_bun/jwt"
	"github.com/oscar-mgh/gin_bun/middleware"
	"github.com/oscar-mgh/gin_bun/model"
	"github.com/oscar-mgh/gin_bun/validations"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"golang.org/x/crypto/bcrypt"
)

func RevalidateToken(c *gin.Context) {
	if middleware.ValidarJWT(c.GetHeader("Authorization")) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "invalid jwt",
		})
		return
	}
	token := strings.Split(c.GetHeader("Authorization"), " ")[1]
	claims, _ := jwt.ExtractClaims(token)

	id := claims["id"].(float64)
	email := claims["email"].(string)
	name := claims["name"].(string)

	newToken, _ := jwt.GenerateJWT(int64(id), email, name)
	c.JSON(http.StatusOK, gin.H{
		"user":  email,
		"token": newToken,
	})
}
func RegisterUser(c *gin.Context) {
	userDto := dto.UserDto{}
	if err := c.ShouldBindJSON(&userDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if len(userDto.Name) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "field name required",
		})
	}
	if len(userDto.Email) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "field email required",
		})
		return
	}
	if validations.ValidEmail.FindStringSubmatch(userDto.Email) == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "invalid email",
		})
		return
	}
	if !validations.ValidPassword(userDto.Password) {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "password must be between 8 and 22 characters long, have an upper case letter and a number",
		})
		return
	}

	db := bun.NewDB(db.Connect(), mysqldialect.New())
	exists, _ := db.NewSelect().Model((*model.UserModel)(nil)).Where("email = ?", userDto.Email).Exists(context.TODO())

	if exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "unexpected error",
		})
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(userDto.Password), bcrypt.DefaultCost)
	model := &model.UserModel{
		Name:      userDto.Name,
		Email:     userDto.Email,
		Phone:     userDto.Phone,
		Password:  string(hash),
		ProfileID: userDto.ProfileID,
	}

	_, err := db.NewInsert().Model(model).Exec(context.TODO())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "unexpected error",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "created",
		"message": "model created successfully",
	})
}

func LoginUser(c *gin.Context) {
	loginDto := dto.LoginDto{}

	if err := c.ShouldBindJSON(&loginDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if len(loginDto.Email) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "field email required",
		})
		return
	}

	if validations.ValidEmail.FindStringSubmatch(loginDto.Email) == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "invalid email",
		})
		return
	}
	if !validations.ValidPassword(loginDto.Password) {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "password must be between 8 and 22 characters long, have an upper case letter and a number",
		})
		return
	}

	user := model.UserModel{}

	db := bun.NewDB(db.Connect(), mysqldialect.New())
	err := db.NewSelect().Model(&user).Where("email = ?", loginDto.Email).Scan(context.TODO())

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "invalid credentials",
		})
		return

	} else {
		password := []byte(loginDto.Password)
		hash := []byte(user.Password)

		err = bcrypt.CompareHashAndPassword(hash, password)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "invalid credentials",
			})
			return

		} else {

			jwt, err := jwt.GenerateJWT(user.ID, user.Email, user.Name)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  "error",
					"message": "cannot generate jwt: " + err.Error(),
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"user":  user.Email,
				"token": jwt,
			})
			return
		}
	}
}
