package mapper

import (
	ara "github.com/diegogub/aranGO"
	"testing"
	"log"
)


type DocTest struct {
	ara.Document // Must include arango Document in every struct you want to save id, key, rev after saving it
	Name     string
	Age      int
	Likes    []string
}

func TestArangoDb_basic(t *testing.T) {

	t.Log("GGOOOOOO 2")

	//change false to true if you want to see every http request
	//Connect(host, user, password string, log bool) (*Session, error) {
	// "http://localhost:8529"
	connection := "http://192.168.99.100:32775/"
	user := ""
	pwd := ""
	s,err := ara.Connect(connection,user,pwd,false)
	if err != nil{
		panic(err)
	}

	// CreateDB(name string,users []User) error
	s.CreateDB("test",nil)

	// create Collections test if exist
	if !s.DB("test").ColExist("docs1"){
		// CollectionOptions has much more options, here we just define name , sync
		docs1 := ara.NewCollectionOptions("docs1",true)
		err = s.DB("test").CreateCollection(docs1)
		if (err != nil) {
			log.Fatal("Failed to create collection 1:", err)
		}
	}

	if !s.DB("test").ColExist("docs2"){
		docs2 := ara.NewCollectionOptions("docs2",true)
		err = s.DB("test").CreateCollection(docs2)
		if (err != nil) {
			log.Fatal("Failed to create collection 2:", err)
		}

	}

	if !s.DB("test").ColExist("ed"){
		edges := ara.NewCollectionOptions("ed",true)
		edges.IsEdge() // set to Edge
		err = s.DB("test").CreateCollection(edges)
		if (err != nil) {
			log.Fatal("Failed to create collection 3:", err)
		}

	}

	t.Log(" Create and Relate documents")

	var d1,d2 DocTest
	d1.Name = "Diego"
	d1.Age = 22
	d1.Likes = []string { "arangodb", "golang", "linux" }

	d2.Name = "Facundo"
	d2.Age = 25
	d2.Likes = []string { "php", "linux", "python" }


	err =s.DB("test").Col("docs1").Save(&d1)
	err =s.DB("test").Col("docs2").Save(&d2)
	if err != nil {
		panic(err)
	}

	// could also check error in document
	/*
	if d1.Error {
	  panic(d1.Message)
	}
	*/

	// update document
	d1.Age = 23
	err =s.DB("test").Col("docs1").Replace(d1.Key,d1)
	if err != nil {
		panic(err)
	}

	t.Log(" Relate documents")

	s.DB("test").Col("ed").Relate(d1.Id,d2.Id,map[string]interface{}{ "is" : "friend" })
}



