package infrastructure

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Router struct {
	Gin *gin.Engine
}

func NewRouter() Router {
	httpRouter := gin.Default()
	httpRouter.GET("/health-check", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "Api is runnings."})
	})
	httpRouter.Use()
	return Router{
		Gin: httpRouter,
	}
}
