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
			handlers.AddEssentailData(c, driver)
		})
		api.POST("/nodes/data", func(c *gin.Context) {
			handlers.AddIntrovertPreferencesHandler(c, driver)
		})
		api.POST("/relationships", func(c *gin.Context) {
			handlers.CreateSimpleRelationshipHandler(c, driver)
		})
		api.POST("/login", func(c *gin.Context) {
			handlers.LoginHandler(c, driver)
		})
		api.POST("/register", func(c *gin.Context) {
			handlers.RegisterHandler(c, driver)
		})
		api.GET("/authtest", func(c *gin.Context) {
			handlers.Authtest(c)
		})
		api.POST("/change-password", func(c *gin.Context) {
			handlers.ChangePasswordHandler(c, driver)
		})
		api.POST("/AddEssentialData", func(c *gin.Context) {
			handlers.AddEssentailData(c, driver)
		})
		api.POST("/AddIntrovertPreferences", func(c *gin.Context) {
			handlers.AddIntrovertPreferencesHandler(c, driver)
		})
		api.POST("/match-and-assign-with-attributes", func(c *gin.Context) {
			handlers.MatchAndAssignRelationshipWithAttributes(c, driver)
		})

		return router
	}
}
