package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func CreateSimpleRelationshipHandler(c *gin.Context, driver neo4j.Driver) {
	var input struct {
		Person1      string `json:"person1"`
		Person2      string `json:"person2"`
		Relationship string `json:"relationship"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	if input.Person1 == "" || input.Person2 == "" || input.Relationship == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "All fields (person1, person2, relationship) are required"})
		return
	}

	input.Relationship = strings.ToUpper(input.Relationship)
	if !isValidRelationshipType(input.Relationship) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid relationship type"})
		return
	}

	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	query := `
		MATCH (a:User {username: $person1}), (b:User {username: $person2})
		CREATE (a)-[r:` + input.Relationship + `]->(b)
		RETURN r
	`

	params := map[string]interface{}{
		"person1": input.Person1,
		"person2": input.Person2,
	}

	result, err := session.Run(query, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create relationship: " + err.Error()})
		return
	}

	if !result.Next() {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Relationship creation failed. Please ensure both users exist."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Relationship created successfully", "relationship": input.Relationship})
}

func isValidRelationshipType(relationship string) bool {
	validRelationships := []string{"FRIENDS", "LIKES", "COLLEAGUES", "FAMILY", "FOLLOWS"}
	for _, valid := range validRelationships {
		if relationship == valid {
			return true
		}
	}
	return false
}
func MatchAndAssignRelationshipWithAttributes(c *gin.Context, driver neo4j.Driver) {
	var input struct {
		Person1 string `json:"person1"`
		Person2 string `json:"person2"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	query := `
		MATCH (a:User {username: $person1}), (b:User {username: $person2})
		WITH a, b,
			[
				size([x IN a.movies_likes WHERE x IN b.movies_likes]),
				size([x IN a.movies_dislikes WHERE x IN b.movies_dislikes]),
				size([x IN a.games_likes WHERE x IN b.games_likes]),
				size([x IN a.games_dislikes WHERE x IN b.games_dislikes]),
				size([x IN a.books_likes WHERE x IN b.books_likes]),
				size([x IN a.books_dislikes WHERE x IN b.books_dislikes]),
				size([x IN a.music_likes WHERE x IN b.music_likes]),
				size([x IN a.music_dislikes WHERE x IN b.music_dislikes]),
				size([x IN a.art_hobbies WHERE x IN b.art_hobbies]),
				size([x IN a.outdoors_likes WHERE x IN b.outdoors_likes]),
				size([x IN a.outdoors_dislikes WHERE x IN b.outdoors_dislikes]),
				size([x IN a.fitness_hobbies WHERE x IN b.fitness_hobbies]),
				size([x IN a.social_hobbies WHERE x IN b.social_hobbies])
			] AS match_counts
		WITH a, b, 
			match_counts,
			reduce(total = 0, count IN match_counts | total + count) AS total_matches,
			CASE
				WHEN reduce(total = 0, count IN match_counts | total + count) <= 0 THEN 1
				WHEN reduce(total = 0, count IN match_counts | total + count) >= 100 THEN 10
				ELSE reduce(total = 0, count IN match_counts | total + count) / 10
			END AS score
		MERGE (a)-[r:SIMILARITY_SCORE]->(b)
		SET r.score = score,
			r.match_counts = match_counts,
			r.total_matches = total_matches
		RETURN r.score AS score, r.match_counts AS match_counts
	`

	params := map[string]interface{}{
		"person1": input.Person1,
		"person2": input.Person2,
	}

	result, err := session.Run(query, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process relationship: " + err.Error()})
		return
	}

	if result.Next() {
		score, _ := result.Record().Get("score")
		matchCounts, _ := result.Record().Get("match_counts")
		c.JSON(http.StatusOK, gin.H{
			"message":      "Relationship assigned successfully",
			"score":        score,
			"match_counts": matchCounts,
		})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "Users not found or no preferences matched"})
	}
}

func DeleteRelationshipHandler(c *gin.Context, driver neo4j.Driver) {
	var input struct {
		Person1      string `json:"person1"`
		Person2      string `json:"person2"`
		Relationship string `json:"relationship"`
	}

	// Parse input JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	// Validate input
	if input.Person1 == "" || input.Person2 == "" || input.Relationship == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "All fields (person1, person2, relationship) are required"})
		return
	}

	// Ensure relationship type is valid
	input.Relationship = strings.ToUpper(input.Relationship)
	if !isValidRelationshipType(input.Relationship) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid relationship type"})
		return
	}

	// Start a Neo4j session
	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	// Query to delete the specific relationship
	query := `
		MATCH (a:User {username: $person1})-[r:` + input.Relationship + `]->(b:User {username: $person2})
		DELETE r
		RETURN COUNT(r) AS deletedCount
	`

	params := map[string]interface{}{
		"person1": input.Person1,
		"person2": input.Person2,
	}

	result, err := session.Run(query, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete relationship: " + err.Error()})
		return
	}

	// Check if any relationships were deleted
	if result.Next() {
		deletedCount, _ := result.Record().Get("deletedCount")
		if deletedCount.(int64) > 0 {
			c.JSON(http.StatusOK, gin.H{"message": "Relationship deleted successfully"})
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "No matching relationship found"})
		}
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete relationship"})
	}
}
