package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/oscar-mgh/gin_bun/routes"
)

func main() {
	godotenv.Load()

	r := gin.Default()

	r.Static("/public", "./public")
	r.LoadHTMLFiles("public/index.html")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	routes.AssignRoutes(r)

	r.Run(":" + os.Getenv("PORT"))
}
