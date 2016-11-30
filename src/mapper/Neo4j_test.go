package mapper


import (
	_ "gopkg.in/cq.v1"
	"testing"
	//"database/sql"
	//"log"
	_ "github.com/jmcvetta/neoism"
	"database/sql"
	"log"
	"fmt"
	"github.com/jmcvetta/neoism"
)


func TestNeoDb_neoism(t *testing.T) {
	db, err := neoism.Connect("http://neo4j:admin@192.168.99.100:32781/db/data")
	if err != nil{
		panic(err)
	}

	kirk := "Captain Kirk"
	mccoy := "Dr McCoy"
	//
	// Create a node
	//
	n0, _ := db.CreateNode(neoism.Props{"name": kirk})
	defer n0.Delete()  // Deferred clean up
	n0.AddLabel("Person") // Add a label
	//
	// Create a node with a Cypher query
	//
	res0 := []struct {
		N neoism.Node // Column "n" gets automagically unmarshalled into field N
	}{}
	cq0 := neoism.CypherQuery{
		Statement: "CREATE (n:Person {name: {name}}) RETURN n",
		// Use parameters instead of constructing a query string
		Parameters: neoism.Props{"name": mccoy},
		Result:     &res0,
	}
	db.Cypher(&cq0)
	n1 := res0[0].N // Only one row of data returned
	n1.Db = db // Must manually set Db with objects returned from Cypher query
	//
	// Create a relationship
	//
	n1.Relate("reports to", n0.Id(), neoism.Props{}) // Empty Props{} is okay
	//
	// Issue a query
	//
	res1 := []struct {
		A   string `json:"a.name"` // `json` tag matches column name in query
		Rel string `json:"type(r)"`
		B   string `json:"b.name"`
	}{}
	cq1 := neoism.CypherQuery{
		// Use backticks for long statements - Cypher is whitespace indifferent
		Statement: `
			MATCH (a:Person)-[r]->(b)
			WHERE a.name = {name}
			RETURN a.name, type(r), b.name
		`,
		Parameters: neoism.Props{"name": mccoy},
		Result:     &res1,
	}
	db.Cypher(&cq1)
	r := res1[0]
	fmt.Println(r.A, r.Rel, r.B)
	//
	// Clean up using a transaction
	//
	qs := []*neoism.CypherQuery{
		&neoism.CypherQuery{
			Statement: `
				MATCH (n:Person)-[r]->()
				WHERE n.name = {name}
				DELETE r
			`,
			Parameters: neoism.Props{"name": mccoy},
		},
		&neoism.CypherQuery{
			Statement: `
				MATCH (n:Person)
				WHERE n.name = {name}
				DELETE n
			`,
			Parameters: neoism.Props{"name": mccoy},
		},
	}
	tx, _ := db.Begin(qs)
	tx.Commit()
	t.Log("WRITING::::::::::000:")
}

func TestNeoDb_basic0(t *testing.T) {

	db, err := sql.Open("neo4j-cypher", "http://neo4j:admin@192.168.99.100:32781")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	t.Log("WRITING::::::::::::::::::")

	/**
	 Write to DB
	 */
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt2, err := tx.Prepare("create (:User {screenName:{0}})")
	if err != nil {
		log.Fatal(err)
	}

	stmt2.Exec("wefreema")
	//stmt2.Exec("JnBrymn")
	//stmt2.Exec("technige")

	err2 := tx.Commit()
	if err2 != nil {
		log.Fatal(err2)
	}


}
func TestNeoDb_write0(t *testing.T) {
	db, err := sql.Open("neo4j-cypher", "http://neo4j:admin@192.168.99.100:32781")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	//stmt, err := tx.Prepare("create (:User {screenName:{0}})")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//stmt.Exec("wefreema")
	//stmt.Exec("JnBrymn")
	//stmt.Exec("technige")
	//
	err2 := tx.Commit()
	if err2 != nil {
		log.Fatal(err2)
	}
}

