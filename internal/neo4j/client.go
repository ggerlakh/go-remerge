package neo4j

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"strings"
)

func HelloWorld(ctx context.Context, uri, username, password string) (string, error) {
	driver, err := neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		return "", err
	}
	defer driver.Close(ctx)

	session := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	greeting, err := session.ExecuteWrite(ctx, func(transaction neo4j.ManagedTransaction) (any, error) {
		result, err := transaction.Run(ctx,
			"CREATE (a:Greeting) SET a.message = $message RETURN a.message + ', from node ' + id(a)",
			map[string]any{"message": "hello, world"})
		if err != nil {
			return nil, err
		}

		if result.Next(ctx) {
			return result.Record().Values[0], nil
		}

		return nil, result.Err()
	})
	if err != nil {
		return "", err
	}

	return greeting.(string), nil
}

func ExecCypher(ctx context.Context, uri, username, password string, cypherQueries []string) error {
	driver, err := neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		return err
	}
	defer driver.Close(ctx)

	session := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	_, err = session.ExecuteWrite(ctx, func(transaction neo4j.ManagedTransaction) (any, error) {
		var result neo4j.ResultWithContext
		for _, cypher := range cypherQueries {
			if strings.Contains(cypher, `\u`) {
				cypher = strings.ReplaceAll(cypher, `\u`, `\\u`)
			}
			result, err = transaction.Run(ctx, cypher, nil)
		}
		if err != nil {
			fmt.Println("line 53")
			return nil, err
		}

		return nil, result.Err()
	})
	if err != nil {
		return err
	}

	return nil
}
