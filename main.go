package main

import (
	"fmt"
	"log"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func main() {
	uri := "bolt://localhost:7687"

	driver, err := neo4j.NewDriver(uri, neo4j.NoAuth())
	if err != nil {
		log.Fatalf("Failed to create driver: %v", err)
	}
	defer driver.Close()

	session := driver.NewSession(neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeRead,
	})
	defer session.Close()

	result, err := session.Run("SHOW DATABASES", map[string]interface{}{})
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}

	var databases []string
	for result.Next() {
		dbName, _ := result.Record().Get("name")
		databases = append(databases, dbName.(string))
	}

	for _, db := range databases {
		fmt.Printf("\nCounting nodes in database '%s':\n", db)

		session := driver.NewSession(neo4j.SessionConfig{
			DatabaseName: db,
			AccessMode:   neo4j.AccessModeRead,
		})
		defer session.Close()

		result, err := session.Run("MATCH (n) RETURN COUNT(n) AS count", map[string]interface{}{})
		if err != nil {
			log.Printf("Query failed for database %s: %v", db, err)
			continue
		}

		if result.Next() {
			nodeCount, _ := result.Record().Get("count")
			fmt.Printf("Total nodes: %v\n", nodeCount)
		}

		if err = result.Err(); err != nil {
			log.Printf("Error iterating result for database %s: %v", db, err)
		}
	}
}
