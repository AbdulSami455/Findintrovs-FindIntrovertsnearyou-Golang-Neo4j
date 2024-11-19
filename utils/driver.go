package utils

import (
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func InitializeDriver(uri string, auth neo4j.AuthToken) (neo4j.Driver, error) {
	driver, err := neo4j.NewDriver(uri, auth)
	if err != nil {
		return nil, fmt.Errorf("failed to create driver: %w", err)
	}
	return driver, nil
}
