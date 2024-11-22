package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func loginHandler(c *gin.Context, driver neo4j.Driver) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()

	query := `
		MATCH (u:User {username: $username, password: $password})
		RETURN u.username AS username
	`
	params := map[string]interface{}{
		"username": input.Username,
		"password": input.Password,
	}

	result, err := session.Run(query, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.Next() {
		username, _ := result.Record().Get("username")
		c.JSON(http.StatusOK, gin.H{"message": "Login successful", "username": username})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
	}
}

func registerHandler(c *gin.Context, driver neo4j.Driver) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	query := `
		CREATE (u:User {username: $username, password: $password})
		RETURN u.username AS username
	`
	params := map[string]interface{}{
		"username": input.Username,
		"password": input.Password,
	}

	_, err := session.Run(query, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created", "username": input.Username})
}
