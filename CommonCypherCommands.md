# 1. Basic CRUD Operations
## CREATE: Create nodes and relationships.
CREATE (n:Person {name: "Bob", age: 25})

## MATCH: Retrieve nodes or relationships that match specified conditions.
MATCH (n:Person {name: "Bob"})
RETURN n

## DELETE: Remove nodes or relationships.
MATCH (n:Person {name: "Alice"})
DELETE n

## DETACH DELETE: Delete nodes along with all their relationships.
MATCH (n:Person {name: "Alice"})
DETACH DELETE n

## SET: Update properties or labels.
MATCH (n:Person {name: "Bob"})
SET n.age = 26
SET n :Employee

## REMOVE: Remove properties or labels from a node.
MATCH (n:Person {name: "Bob"})
REMOVE n.age
REMOVE n:Employee

# 2. Relationships

## CREATE Relationship: Define relationships between nodes.
MATCH (a:Person {name: "Alice"}), (b:Person {name: "Bob"})
CREATE (a)-[:KNOWS]->(b)

## MATCH Relationship: Find specific relationships.
MATCH (a:Person)-[r:KNOWS]->(b:Person)
RETURN a, b

DELETE Relationship: Remove a relationship.
MATCH (a:Person {name: "Alice"})-[r:KNOWS]->(b:Person {name: "Bob"})
DELETE r

3. Aggregation and Filtering
## COUNT: Count nodes or relationships.
MATCH (n:Person)
RETURN COUNT(n)

## WHERE: Filter based on specific conditions.
MATCH (n:Person)
WHERE n.age > 25 AND n.city = "New York"
RETURN n

ORDER BY: Sort query results.

cypher
Copy code
MATCH (n:Person)
RETURN n.name, n.age
ORDER BY n.age DESC
LIMIT: Limit the number of returned results.

cypher
Copy code
MATCH (n:Person)
RETURN n.name
LIMIT 5
4. Pattern Matching
Variable-Length Relationships: Find relationships with a specific length.

cypher
Copy code
MATCH (a:Person)-[:FRIENDS*1..3]-(b:Person)
WHERE a.name = "Alice"
RETURN b
Shortest Path: Find the shortest path between two nodes.

cypher
Copy code
MATCH p = shortestPath((a:Person {name: "Alice"})-[:FRIENDS*..]-(b:Person {name: "Bob"}))
RETURN p
5. With Clauses for Complex Queries
WITH: Allows chaining queries or passing variables to the next part of a query.

cypher
Copy code
MATCH (n:Person)
WITH n, n.age AS age
WHERE age > 25
RETURN n.name
UNWIND: Expand a list into individual rows.

cypher
Copy code
UNWIND [1, 2, 3] AS number
RETURN number
6. Indexes and Constraints
CREATE INDEX: Speed up searches on certain properties.

cypher
Copy code
CREATE INDEX FOR (n:Person) ON (n.name)
CREATE UNIQUE CONSTRAINT: Ensure unique property values.

cypher
Copy code
CREATE CONSTRAINT ON (n:Person) ASSERT n.name IS UNIQUE
DROP INDEX/CONSTRAINT: Remove indexes or constraints.

cypher
Copy code
DROP INDEX person_name_index
DROP CONSTRAINT unique_person_name
7. Advanced Graph Traversals
Collect: Gather values into a list.

cypher
Copy code
MATCH (n:Person)
RETURN n.city, collect(n.name)
Foreach: Loop through a list for operations.

cypher
Copy code
MATCH (n:Person {name: "Alice"})
FOREACH (x IN ["reading", "traveling", "cooking"] |
  CREATE (n)-[:LIKES]->(:Hobby {name: x}))