package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"github.com/oscar-mgh/gin_bun/db"
	"github.com/oscar-mgh/gin_bun/dto"
	"github.com/oscar-mgh/gin_bun/model"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
)

func FindAllMovies(c *gin.Context) {
	db := bun.NewDB(db.Connect(), mysqldialect.New())
	movies := []*model.MovieModel{}

	err := db.NewSelect().Model(&movies).Relation("Genre").Scan(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, movies)
}
func GetMovieById(c *gin.Context) {
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

	movie := model.MovieModel{}
	err = db.NewSelect().Model(&movie).Where("id = ?", id).Limit(1).Scan(context.TODO())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "resource not found",
		})
		return
	}

	c.JSON(http.StatusOK, movie)
}
func CreateMovie(c *gin.Context) {
	movieDto := dto.MovieDto{}
	if err := c.ShouldBindJSON(&movieDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if len(movieDto.Name) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "name required",
		})
		return
	}

	db := bun.NewDB(db.Connect(), mysqldialect.New())
	movie := &model.MovieModel{
		Name:        movieDto.Name,
		Description: movieDto.Description,
		Year:        movieDto.Year,
		GenreID:     int64(movieDto.GenreID),
		Slug:        slug.Make(movieDto.Name),
	}
	_, err := db.NewInsert().Model(movie).Exec(context.TODO())
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
func UpdateMovie(c *gin.Context) {
	movieDto := dto.MovieDto{}
	if err := c.ShouldBindJSON(&movieDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if len(movieDto.Name) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "name required",
		})
		return
	}

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

	movie := &model.MovieModel{
		Name:        movieDto.Name,
		Description: movieDto.Description,
		Year:        movieDto.Year,
		GenreID:     int64(movieDto.GenreID),
		Slug:        slug.Make(movieDto.Name),
	}
	_, err = db.NewUpdate().Model(movie).Where("id = ?", id).Exec(context.TODO())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "unexpected error",
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"status":  "updated",
		"message": "model updated successfully",
	})
}
func DeleteMovie(c *gin.Context) {
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

	_, err = db.NewDelete().Model((*model.MovieModel)(nil)).Where("id = ?", id).Exec(context.TODO())
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
