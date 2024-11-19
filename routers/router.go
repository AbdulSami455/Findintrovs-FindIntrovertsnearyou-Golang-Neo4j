package routers

import (
	"my-go-project/handlers"

	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func SetupRouter(driver neo4j.Driver) *gin.Engine {
	router := gin.Default()

	api := router.Group("/api")
	{
		api.GET("/databases", func(c *gin.Context) {
			handlers.ListDatabasesHandler(c, driver)
		})
		api.GET("/databases/:dbname/count", func(c *gin.Context) {
			handlers.CountNodesHandler(c, driver)
		})
		api.POST("/nodes", func(c *gin.Context) {
			handlers.CreateOrUpdateNodeHandler(c, driver)
		})
		api.POST("/relationships", func(c *gin.Context) {
			handlers.CreateRelationshipHandler(c, driver)
		})
	}

	return router
}
