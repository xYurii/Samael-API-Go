package router

import (
	"os"

	"github.com/gin-gonic/gin"
)

func Initiliaze() {
	r := gin.Default()

	// initiliaze the routes
	initializeRoutes(r)

	port := os.Getenv("PORT")
	r.Run(":" + port)
}
