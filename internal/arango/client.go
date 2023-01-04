package arango

import (
	"context"
	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	"log"
)

func LoadGraph(ctx context.Context, endpoints []string, username, password, dbName, graphName string, nodes, edges []map[string]any) {
	var db driver.Database
	var edgeDefinition driver.EdgeDefinition
	var graphOptions driver.CreateGraphOptions
	nodeCollName := graphName + "_nodes"
	edgeCollName := graphName + "_edges"

	// Create an HTTP connection to the database
	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: endpoints,
	})
	if err != nil {
		log.Fatalf("Failed to create HTTP connection: %v", err)
	}
	// Create a client
	c, err := driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.BasicAuthentication(username, password),
	})

	dbExists, err := c.DatabaseExists(ctx, dbName)

	if !dbExists {
		// Create database
		db, err = c.CreateDatabase(ctx, dbName, nil)
		if err != nil {
			log.Fatalf("Failed to create database: %v", err)
		}
	} else {
		db, err = c.Database(ctx, dbName)

		if err != nil {
			log.Fatalf("Failed to open existing database: %v", err)
		}
	}

	nodeCollExists, err := db.CollectionExists(nil, nodeCollName)
	if !nodeCollExists {
		_, err = db.CreateCollection(ctx, nodeCollName, nil)

		if err != nil {
			log.Fatalf("Failed to create collection: %v", err)
		}
	}

	edgeCollExists, err := db.CollectionExists(nil, edgeCollName)
	if !edgeCollExists {
		edgeOptions := driver.CreateCollectionOptions{
			Type: driver.CollectionTypeEdge,
		}
		_, err = db.CreateCollection(ctx, edgeCollName, &edgeOptions)

		if err != nil {
			log.Fatalf("Failed to create collection: %v", err)
		}
	}
	// define the edgeCollection to store the edges
	edgeDefinition.Collection = edgeCollName
	// define a set of collections where an edge is going out...
	edgeDefinition.From = []string{nodeCollName}

	// repeat this for the collections where an edge is going into
	edgeDefinition.To = []string{nodeCollName}

	// A graph can contain additional vertex collections, defined in the set of orphan collections
	graphOptions.EdgeDefinitions = []driver.EdgeDefinition{edgeDefinition}

	// now it's possible to create a graph
	graph, err := db.CreateGraphV2(ctx, graphName, &graphOptions)
	if err != nil {
		log.Fatalf("Failed to create graph: %v", err)
	}

	// add vertex
	vertexCollection1, err := graph.VertexCollection(nil, nodeCollName)
	if err != nil {
		log.Fatalf("Failed to get vertex collection: %v", err)
	}

	//nodesCtx := driver.WithWaitForSync(ctx)
	// add nodes
	_, _, err = vertexCollection1.CreateDocuments(ctx, nodes)
	if err != nil {
		log.Fatalf("Failed to create vertex documents: %v", err)
	}

	//edgesCtx := driver.WithWaitForSync(ctx)
	// add edges
	edgeCollection, _, err := graph.EdgeCollection(ctx, edgeCollName)
	if err != nil {
		log.Fatalf("Failed to select edge collection: %v", err)
	}

	//_, _, err = edgeCollection.CreateDocuments(ctx, edges)
	for _, edge := range edges {
		_, err = edgeCollection.CreateDocument(ctx, edge)
		if err != nil {
			log.Fatalf("Failed to create edge document: %v", err)
		}
	}
}
