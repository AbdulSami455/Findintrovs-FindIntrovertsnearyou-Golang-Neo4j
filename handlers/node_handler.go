package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func CreateOrUpdateNodeHandler(c *gin.Context, driver neo4j.Driver) {
	var input struct {
		Name       string `json:"name"`
		Age        int    `json:"age"`
		Gender     string `json:"gender"`
		Occupation string `json:"occupation"`
		Institute  string `json:"institute"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	query := `
		MERGE (n:User {name: $name})
		ON CREATE SET n.age = $age, n.gender = $gender, n.occupation = $occupation, n.institute = $institute
		ON MATCH SET n.age = $age, n.gender = $gender, n.occupation = $occupation, n.institute = $institute
	`
	params := map[string]interface{}{
		"name":       input.Name,
		"age":        input.Age,
		"gender":     input.Gender,
		"occupation": input.Occupation,
		"institute":  input.Institute,
	}

	_, err := session.Run(query, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Node '%s' created or updated", input.Name)})
}
