package handlers

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/oscar-mgh/gin_bun/db"
	"github.com/oscar-mgh/gin_bun/model"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
)

func UploadImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	extension := strings.Split(file.Filename, ".")[1]
	time := strings.Split(time.Now().String(), " ")
	image := string(time[4][7:14]) + "." + extension
	archive := "public/uploads/images/" + image

	log.Println(file.Filename)

	c.SaveUploadedFile(file, archive)

	id := c.Param("id")
	numId, _ := strconv.ParseInt(id, 10, 64)
	db := bun.NewDB(db.Connect(), mysqldialect.New())
	img := &model.MovieImage{
		Name:    image,
		MovieID: int64(numId),
	}

	_, err = db.NewInsert().Model(img).Exec(context.TODO())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "unexpected error",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "created",
		"message": "image uploaded",
	})
}
func GetMovieImages(c *gin.Context) {
	id := c.Param("id")
	db := bun.NewDB(db.Connect(), mysqldialect.New())

	exists, err := db.NewSelect().Model((*model.MovieModel)(nil)).Where("id = ?", id).Exists(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "resource not found",
		})
		return
	}

	images := []model.MovieImage{}

	err = db.NewSelect().Model(&images).Where("movie_id = ?", id).Scan(context.TODO())
	if err != nil {
		log.Println(err)
	}

	c.JSON(http.StatusOK, images)
}
func DeleteImage(c *gin.Context) {
	id := c.Param("imageId")
	image := model.MovieImage{}
	db := bun.NewDB(db.Connect(), mysqldialect.New())
	err := db.NewSelect().Model(image).Scan(context.TODO())
	if err != nil {
		log.Println(err)
	}

	filePath := "public/uploads/images/" + image.Name
	err = os.Remove(filePath)
	if err != nil {
		log.Println(err)
	}

	_, err = db.NewDelete().Model((*model.MovieImage)(nil)).Where("id = ?", id).Exec(context.TODO())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "unexpected error",
		})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{
		"status":  "no content",
		"message": "model deleted successfully",
	})
}
