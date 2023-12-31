package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/oscar-mgh/gin_bun/handlers"
)

func AssignRoutes(r *gin.Engine) {
	prefixPath := r.Group("/api/v1")
	{
		genres := prefixPath.Group("/genres")
		{
			genres.GET("", handlers.FindAllGenres)
			genres.POST("", handlers.CreateGenre)
			genres.GET("/:id", handlers.GetGenreById)
			genres.PUT("/:id", handlers.UpdateGenre)
			genres.DELETE("/:id", handlers.DeleteGenre)
		}
		movies := prefixPath.Group("/movies")
		{
			movies.GET("", handlers.FindAllMovies)
			movies.POST("", handlers.CreateMovie)
			movies.GET("/:id", handlers.GetMovieById)
			movies.PUT("/:id", handlers.UpdateMovie)
			movies.DELETE("/:id", handlers.DeleteMovie)
		}
		images := prefixPath.Group("/images")
		{
			images.GET("/:id", handlers.GetMovieImages)
			images.POST("/:id", handlers.UploadImage)
			images.DELETE("/:imageId", handlers.DeleteImage)
		}
		users := prefixPath.Group("/users")
		{
			users.POST("/register", handlers.RegisterUser)
			users.POST("/login", handlers.LoginUser)
			users.POST("/revalidate", handlers.RevalidateToken)
		}
	}
}
