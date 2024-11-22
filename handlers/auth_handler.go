package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"golang.org/x/crypto/bcrypt"
)

func Authtest(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "auth test"})
}

func RegisterHandler(c *gin.Context, driver neo4j.Driver) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
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
		"password": string(hashedPassword),
	}

	_, err = session.Run(query, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created", "username": input.Username})
}

func LoginHandler(c *gin.Context, driver neo4j.Driver) {
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
		MATCH (u:User {username: $username})
		RETURN u.password AS password
	`
	params := map[string]interface{}{
		"username": input.Username,
	}

	result, err := session.Run(query, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.Next() {
		storedHashedPassword, _ := result.Record().Get("password")

		err := bcrypt.CompareHashAndPassword([]byte(storedHashedPassword.(string)), []byte(input.Password))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Login successful", "username": input.Username})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
	}
}
