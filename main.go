package main

import (
	"log"
	"my-go-project/routers"
	"my-go-project/utils"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func main() {
	driver, err := utils.InitializeDriver("bolt://localhost:7687", neo4j.NoAuth())
	if err != nil {
		log.Fatalf("Failed to initialize Neo4j driver: %v", err)
	}
	defer driver.Close()

	router := routers.SetupRouter(driver)
	log.Fatal(router.Run(":8070"))
}
