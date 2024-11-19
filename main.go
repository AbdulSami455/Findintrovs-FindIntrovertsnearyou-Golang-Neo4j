package main

import (
	"log"
	"my-go-project/routers"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func main() {
	driver, err := initializeDriver("bolt://localhost:7687", neo4j.NoAuth())
	if err != nil {
		log.Fatalf("Failed to initialize Neo4j driver: %v", err)
	}
	defer driver.Close()

	router := routers.SetupRouter(driver)
	log.Fatal(router.Run(":8080")) // Start the server on port 8080
}
