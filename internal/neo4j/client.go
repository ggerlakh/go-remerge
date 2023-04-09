package neo4j

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"log"
	"strings"
)

func ExecCypher(ctx context.Context, uri, username, password string, cypherQueries []string) error {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	driver, err := neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		return err
	}
	defer func(driver neo4j.DriverWithContext, ctx context.Context) {
		err := driver.Close(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}(driver, ctx)

	session := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer func(session neo4j.SessionWithContext, ctx context.Context) {
		err := session.Close(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}(session, ctx)

	_, err = session.ExecuteWrite(ctx, func(transaction neo4j.ManagedTransaction) (any, error) {
		//var result neo4j.ResultWithContext
		for _, cypher := range cypherQueries {
			if strings.Contains(cypher, `\u`) {
				cypher = strings.ReplaceAll(cypher, `\u`, `\\u`)
			}
			_, err = transaction.Run(ctx, cypher, nil)
		}
		if err != nil {
			return nil, err
		}
		//log.Println("neo4j result:", result.IsOpen(), result)
		return nil, nil
	})
	if err != nil {
		return err
	}

	return nil
}
