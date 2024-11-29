package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func AddEssentailData(c *gin.Context, driver neo4j.Driver) {
	var input struct {
		Username   string `json:"username"`
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
		MERGE (n:User {username: $username})
		ON CREATE SET n.age = $age, n.gender = $gender, n.occupation = $occupation, n.institute = $institute
		ON MATCH SET n.age = $age, n.gender = $gender, n.occupation = $occupation, n.institute = $institute
	`
	params := map[string]interface{}{
		"username":   input.Username,
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

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Node for username '%s' created or updated", input.Username)})
}

func AddIntrovertPreferencesHandler(c *gin.Context, driver neo4j.Driver) {
	var input struct {
		Username         string   `json:"username"`
		MoviesLikes      []string `json:"movies_likes"`
		MoviesDislikes   []string `json:"movies_dislikes"`
		GamesLikes       []string `json:"games_likes"`
		GamesDislikes    []string `json:"games_dislikes"`
		BooksLikes       []string `json:"books_likes"`
		BooksDislikes    []string `json:"books_dislikes"`
		MusicLikes       []string `json:"music_likes"`
		MusicDislikes    []string `json:"music_dislikes"`
		ArtHobbies       []string `json:"art_hobbies"`
		OutdoorsLikes    []string `json:"outdoors_likes"`
		OutdoorsDislikes []string `json:"outdoors_dislikes"`
		FitnessHobbies   []string `json:"fitness_hobbies"`
		SocialHobbies    []string `json:"social_hobbies"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	query := `
		MATCH (n:User {username: $username})
		SET n.movies_likes = $moviesLikes,
			n.movies_dislikes = $moviesDislikes,
			n.games_likes = $gamesLikes,
			n.games_dislikes = $gamesDislikes,
			n.books_likes = $booksLikes,
			n.books_dislikes = $booksDislikes,
			n.music_likes = $musicLikes,
			n.music_dislikes = $musicDislikes,
			n.art_hobbies = $artHobbies,
			n.outdoors_likes = $outdoorsLikes,
			n.outdoors_dislikes = $outdoorsDislikes,
			n.fitness_hobbies = $fitnessHobbies,
			n.social_hobbies = $socialHobbies
		RETURN n
	`

	params := map[string]interface{}{
		"username":         input.Username,
		"moviesLikes":      input.MoviesLikes,
		"moviesDislikes":   input.MoviesDislikes,
		"gamesLikes":       input.GamesLikes,
		"gamesDislikes":    input.GamesDislikes,
		"booksLikes":       input.BooksLikes,
		"booksDislikes":    input.BooksDislikes,
		"musicLikes":       input.MusicLikes,
		"musicDislikes":    input.MusicDislikes,
		"artHobbies":       input.ArtHobbies,
		"outdoorsLikes":    input.OutdoorsLikes,
		"outdoorsDislikes": input.OutdoorsDislikes,
		"fitnessHobbies":   input.FitnessHobbies,
		"socialHobbies":    input.SocialHobbies,
	}

	result, err := session.Run(query, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.Next() {
		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Introvert preferences added/updated for username '%s'", input.Username)})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
	}
}
