package main

import (
	"fmt"
	"log"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func main() {
	driver, err := initializeDriver("bolt://localhost:7687", neo4j.NoAuth())
	if err != nil {
		log.Fatalf("Failed to initialize driver: %v", err)
	}
	defer driver.Close()
	/*
		databases, err := listDatabases(driver)

		if err != nil {
			log.Fatalf("Failed to list databases: %v", err)
		}

		for _, db := range databases {
			fmt.Printf("\nCounting nodes in database '%s':\n", db)
			nodeCount, err := countNodesInDatabase(driver, db)
			if err != nil {
				log.Printf("Failed to count nodes in database %s: %v", db, err)
				continue
			}
			fmt.Printf("Total nodes: %v\n", nodeCount)
		}

		createnodeinDatabase(driver, "neo4j")
	*/
	//createdatabase(driver, "introverts")
	//deleteallnodes(driver, "neo4j")
	//createnodeinDatabase(driver, "neo4j", "Abdullah")
	//addOrUpdateProperty(driver, "neo4j", "Abdullah", "email", "abdullah1779@gmail.com")
	calculateMatchScore(driver, "neo4j", "Abdul Sami", "Abdullah")
}

/*
	func createdatabase(driver neo4j.Driver, dbname string) {
		session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
		defer session.Close()

		result, err := session.Run(fmt.Sprintf("CREATE DATABASE %s", dbname), map[string]interface{}{})
		if err != nil {
			log.Fatalf("Failed to create database: %v", err)
		}

		if err := result.Err(); err != nil {
			log.Fatalf("Error creating database: %v", err)
		}
	}
*/

func calculateMatchScore(driver neo4j.Driver, dbname string, personName1 string, personName2 string) {
	session := driver.NewSession(neo4j.SessionConfig{DatabaseName: dbname, AccessMode: neo4j.AccessModeRead})
	defer session.Close()

	query := `
		MATCH (p1:Person {name: $name1}), (p2:Person {name: $name2})
		WITH p1, p2, keys(p1) AS keys1, keys(p2) AS keys2
		RETURN size([key IN keys1 WHERE key IN keys2]) AS matchCount
	`

	params := map[string]interface{}{
		"name1": personName1,
		"name2": personName2,
	}

	result, err := session.Run(query, params)
	if err != nil {
		log.Fatalf("Failed to calculate match score: %v", err)
	}

	if result.Next() {
		matchCount, ok := result.Record().Values[0].(int64)
		if !ok {
			log.Fatalf("Failed to cast match count to int64")
		}
		fmt.Printf("Match score between '%s' and '%s': %d properties matched\n", personName1, personName2, matchCount)
	} else {
		fmt.Printf("No matching nodes found for '%s' and '%s'\n", personName1, personName2)
	}

	if err := result.Err(); err != nil {
		log.Fatalf("Error during result iteration: %v", err)
	}
}

func createnodeinDatabase(driver neo4j.Driver, dbname string, personName string) {
	session := driver.NewSession(neo4j.SessionConfig{DatabaseName: dbname, AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	result, err := session.Run(
		"CREATE (n:Person {name: $name})",
		map[string]interface{}{
			"name": personName,
		},
	)
	if err != nil {
		log.Fatalf("Failed to create node: %v", err)
	}

	if err := result.Err(); err != nil {
		log.Fatalf("Error creating node: %v", err)
	}

	fmt.Printf("Person node with name '%s' created successfully\n", personName)
}

func addOrUpdateProperty(driver neo4j.Driver, dbname string, personName string, propertyName string, propertyValue interface{}) {
	session := driver.NewSession(neo4j.SessionConfig{DatabaseName: dbname, AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	query := `
		MATCH (n:Person {name: $name})
		SET n[$property] = $value
		RETURN n
	`

	params := map[string]interface{}{
		"name":     personName,
		"property": propertyName,
		"value":    propertyValue,
	}

	result, err := session.Run(query, params)
	if err != nil {
		log.Fatalf("Failed to update property: %v", err)
	}

	if result.Next() {
		fmt.Printf("Property '%s' updated/added successfully for person '%s'\n", propertyName, personName)
	} else {
		fmt.Printf("No person with name '%s' found. Property '%s' not updated.\n", personName, propertyName)
	}
}

func essentialNodeData(driver neo4j.Driver, dbname string, personName string, personAge int, personGender string, personOccupation string, personInstitute string) {
	session := driver.NewSession(neo4j.SessionConfig{DatabaseName: dbname, AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	query := `
		MERGE (n:Person {name: $name})
		ON CREATE SET n.age = $age, n.gender = $gender, n.occupation = $occupation, n.institute = $institute
		ON MATCH SET n.age = $age, n.gender = $gender, n.occupation = $occupation, n.institute = $institute
	`

	params := map[string]interface{}{
		"name":       personName,
		"age":        personAge,
		"gender":     personGender,
		"occupation": personOccupation,
		"institute":  personInstitute,
	}

	result, err := session.Run(query, params)
	if err != nil {
		log.Fatalf("Failed to upsert node: %v", err)
	}

	if err := result.Err(); err != nil {
		log.Fatalf("Error during node creation or update: %v", err)
	}

	fmt.Printf("Person node with name '%s' created or updated successfully\n", personName)
}

func deleteallnodes(driver neo4j.Driver, dbname string) {
	session := driver.NewSession(neo4j.SessionConfig{DatabaseName: dbname, AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	result, err := session.Run("MATCH (n) DETACH DELETE n", map[string]interface{}{})
	if err != nil {
		log.Fatalf("Failed to delete nodes: %v", err)
	}
	if err := result.Err(); err != nil {
		log.Fatalf("Error deleting nodes: %v", err)
	}

}
func initializeDriver(uri string, auth neo4j.AuthToken) (neo4j.Driver, error) {
	driver, err := neo4j.NewDriver(uri, auth)
	if err != nil {
		return nil, fmt.Errorf("failed to create driver: %w", err)
	}
	return driver, nil
}

func listDatabases(driver neo4j.Driver) ([]string, error) {
	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()

	result, err := session.Run("SHOW DATABASES", map[string]interface{}{})
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	var databases []string
	for result.Next() {
		dbName, _ := result.Record().Get("name")
		databases = append(databases, dbName.(string))
	}

	if err := result.Err(); err != nil {
		return nil, fmt.Errorf("error iterating result: %w", err)
	}

	return databases, nil
}

func countNodesInDatabase(driver neo4j.Driver, dbName string) (int64, error) {
	session := driver.NewSession(neo4j.SessionConfig{
		DatabaseName: dbName,
		AccessMode:   neo4j.AccessModeRead,
	})
	defer session.Close()

	result, err := session.Run("MATCH (n) RETURN COUNT(n) AS count", map[string]interface{}{})
	if err != nil {
		return 0, fmt.Errorf("query failed for database %s: %w", dbName, err)
	}

	if result.Next() {
		nodeCount, _ := result.Record().Get("count")
		return nodeCount.(int64), nil
	}

	if err := result.Err(); err != nil {
		return 0, fmt.Errorf("error iterating result for database %s: %w", dbName, err)
	}

	return 0, nil
}

func findNodeByName(driver neo4j.Driver, dbname string, personName string) {
	session := driver.NewSession(neo4j.SessionConfig{DatabaseName: dbname, AccessMode: neo4j.AccessModeRead})
	defer session.Close()

	query := `
		MATCH (n:Person {name: $name})
		RETURN n
	`

	params := map[string]interface{}{
		"name": personName,
	}

	result, err := session.Run(query, params)
	if err != nil {
		log.Fatalf("Failed to find node: %v", err)
	}

	if result.Next() {
		node := result.Record().Values[0]
		fmt.Printf("Node found: %v\n", node)
	} else {
		fmt.Printf("No node found with name '%s'\n", personName)
	}
}

func createRelationship(driver neo4j.Driver, dbname string, person1 string, person2 string, relationshipType string) {
	session := driver.NewSession(neo4j.SessionConfig{DatabaseName: dbname, AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	query := `
		MATCH (a:Person {name: $name1}), (b:Person {name: $name2})
		MERGE (a)-[r:%s]->(b)
		RETURN r
	`

	query = fmt.Sprintf(query, relationshipType)

	params := map[string]interface{}{
		"name1": person1,
		"name2": person2,
	}

	result, err := session.Run(query, params)
	if err != nil {
		log.Fatalf("Failed to create relationship: %v", err)
	}

	if result.Next() {
		fmt.Printf("Relationship '%s' created between '%s' and '%s'\n", relationshipType, person1, person2)
	} else {
		fmt.Printf("Failed to create relationship: No matching nodes found\n")
	}
}

func getRelationships(driver neo4j.Driver, dbname string, personName string) {
	session := driver.NewSession(neo4j.SessionConfig{DatabaseName: dbname, AccessMode: neo4j.AccessModeRead})
	defer session.Close()

	query := `
		MATCH (n:Person {name: $name})-[r]-(m)
		RETURN type(r) AS relationship, m.name AS relatedNode
	`

	params := map[string]interface{}{
		"name": personName,
	}

	result, err := session.Run(query, params)
	if err != nil {
		log.Fatalf("Failed to get relationships: %v", err)
	}

	fmt.Printf("Relationships for '%s':\n", personName)

	for result.Next() {
		record := result.Record()

		relationship, _ := record.Get("relationship")
		relatedNode, _ := record.Get("relatedNode")

		fmt.Printf("- %s -> %s\n", relationship, relatedNode)
	}

	if err := result.Err(); err != nil {
		log.Fatalf("Error iterating result: %v", err)
	}
}

func updateNodeProperties(driver neo4j.Driver, dbname string, personName string, properties map[string]interface{}) {
	session := driver.NewSession(neo4j.SessionConfig{DatabaseName: dbname, AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	setClauses := ""
	for key := range properties {
		setClauses += fmt.Sprintf("n.%s = $%s, ", key, key)
	}
	setClauses = setClauses[:len(setClauses)-2]

	query := fmt.Sprintf(`
		MATCH (n:Person {name: $name})
		SET %s
	`, setClauses)

	params := map[string]interface{}{
		"name": personName,
	}
	for key, value := range properties {
		params[key] = value
	}

	_, err := session.Run(query, params)
	if err != nil {
		log.Fatalf("Failed to update properties: %v", err)
	}

	fmt.Printf("Properties for '%s' updated successfully\n", personName)
}

func countRelationships(driver neo4j.Driver, dbname string, personName string) {
	session := driver.NewSession(neo4j.SessionConfig{DatabaseName: dbname, AccessMode: neo4j.AccessModeRead})
	defer session.Close()

	query := `
		MATCH (n:Person {name: $name})-[r]-()
		RETURN COUNT(r) AS relationshipCount
	`

	params := map[string]interface{}{
		"name": personName,
	}

	result, err := session.Run(query, params)
	if err != nil {
		log.Fatalf("Failed to count relationships: %v", err)
	}

	if result.Next() {
		count := result.Record().Values[0].(int64)
		fmt.Printf("Node '%s' has %d relationships\n", personName, count)
	}
}
