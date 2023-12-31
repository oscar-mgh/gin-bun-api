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

// func Query(c *gin.Context) {
// 	slug := c.Query("slug")
// 	if len(slug) == 0 {
// 		c.JSON(http.StatusBadRequest, "query param slug required")
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{
// 		"slug": slug,
// 	})
// }

func FindAllGenres(c *gin.Context) {
	db := bun.NewDB(db.Connect(), mysqldialect.New())
	genres := []*model.GenreModel{}

	err := db.NewSelect().Model(&genres).Scan(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, genres)
}
func GetGenreById(c *gin.Context) {
	id := c.Param("id")
	db := bun.NewDB(db.Connect(), mysqldialect.New())

	genre := model.GenreModel{}
	err := db.NewSelect().Model(&genre).Where("id = ?", id).Limit(1).Scan(context.TODO())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "resource not found",
		})
		return
	}

	c.JSON(http.StatusOK, genre)
}
func CreateGenre(c *gin.Context) {
	genreDto := dto.GenreDto{}
	if err := c.Copy().ShouldBindJSON(&genreDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if len(genreDto.Name) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "name required",
		})
		return
	}

	db := bun.NewDB(db.Connect(), mysqldialect.New())
	genre := &model.GenreModel{
		Name: genreDto.Name,
		Slug: slug.Make(genreDto.Name),
	}
	_, err := db.NewInsert().Model(genre).Exec(context.TODO())
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
func UpdateGenre(c *gin.Context) {
	genreDto := dto.GenreDto{}
	if err := c.Copy().ShouldBindJSON(&genreDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if len(genreDto.Name) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "name required",
		})
		return
	}

	id := c.Param("id")
	db := bun.NewDB(db.Connect(), mysqldialect.New())
	exists, err := db.NewSelect().Model((*model.GenreModel)(nil)).Where("id = ?", id).Exists(context.TODO())
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

	genre := &model.GenreModel{
		Name: genreDto.Name,
		Slug: slug.Make(genreDto.Name),
	}
	_, err = db.NewUpdate().Model(genre).Where("id = ?", id).Exec(context.TODO())
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
func DeleteGenre(c *gin.Context) {
	id := c.Param("id")
	db := bun.NewDB(db.Connect(), mysqldialect.New())
	exists, err := db.NewSelect().Model((*model.GenreModel)(nil)).Where("id = ?", id).Exists(context.TODO())
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

	_, err = db.NewDelete().Model((*model.GenreModel)(nil)).Where("id = ?", id).Exec(context.TODO())
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
