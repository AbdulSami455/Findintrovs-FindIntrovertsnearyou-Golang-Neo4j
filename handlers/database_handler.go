package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func ListDatabasesHandler(c *gin.Context, driver neo4j.Driver) {
	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()

	result, err := session.Run("SHOW DATABASES", nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var databases []string
	for result.Next() {
		dbName, _ := result.Record().Get("name")
		databases = append(databases, dbName.(string))
	}

	c.JSON(http.StatusOK, gin.H{"databases": databases})
}

func CountNodesHandler(c *gin.Context, driver neo4j.Driver) {
	dbName := c.Param("dbname")

	session := driver.NewSession(neo4j.SessionConfig{DatabaseName: dbName, AccessMode: neo4j.AccessModeRead})
	defer session.Close()

	result, err := session.Run("MATCH (n) RETURN COUNT(n) AS count", nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.Next() {
		nodeCount := result.Record().Values[0].(int64)
		c.JSON(http.StatusOK, gin.H{"node_count": nodeCount})
	} else {
		c.JSON(http.StatusOK, gin.H{"node_count": 0})
	}
}
