package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func CreateRelationshipHandler(c *gin.Context, driver neo4j.Driver) {
	var input struct {
		Person1      string `json:"person1"`
		Person2      string `json:"person2"`
		Relationship string `json:"relationship"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	query := fmt.Sprintf(`
		MATCH (a:User {name: $person1}), (b:User {name: $person2})
		MERGE (a)-[r:%s]->(b)
	`, input.Relationship)

	params := map[string]interface{}{
		"person1": input.Person1,
		"person2": input.Person2,
	}

	_, err := session.Run(query, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Relationship '%s' created between '%s' and '%s'", input.Relationship, input.Person1, input.Person2)})
}
