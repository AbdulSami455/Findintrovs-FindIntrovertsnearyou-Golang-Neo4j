package routers

import (
	"my-go-project/handlers"

	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func SetupAuthRouter(router *gin.Engine, driver neo4j.Driver) {
	auth := router.Group("/auth")
	{
		auth.POST("/login", func(c *gin.Context) {
			handlers.LoginHandler(c, driver)
		})
		auth.POST("/register", func(c *gin.Context) {
			handlers.RegisterHandler(c, driver)

		})
	}
}
