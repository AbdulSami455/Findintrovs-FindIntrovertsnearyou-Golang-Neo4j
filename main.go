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
		AccessMode: neo4j.AccessModeWrite,
	})
	defer session.Close()

	result, err := session.Run("MATCH (n) RETURN COUNT(n)", map[string]interface{}{})
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}

	if result.Next() {
		count, _ := result.Record().Get("COUNT(n)")
		fmt.Printf("Node count: %v\n", count)
	}

	if err = result.Err(); err != nil {
		log.Fatalf("Error iterating result: %v", err)
	}
}
