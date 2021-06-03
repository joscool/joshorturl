package main

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/graphql-go/graphql"
)

func TestGqlShortenURL(t *testing.T) {

	// set mock service, mock service define in service test
	businessLogic = serviceTest

	expectedJson := `{"data":{"shortenURL":"https//www.joshua.heroku.com/ODEzYW"}}`

	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: ShortenURLSchema}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery),
		Mutation: mutationType,
	}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}

	query := `
	mutation {
		shortenURL(url:"https://www.golangprograms.com/example-function-that-takes-an-interface-type-as-value-and-pointer.html")
	   }	
	`
	params := graphql.Params{Schema: schema, RequestString: query}
	r := graphql.Do(params)
	if len(r.Errors) > 0 {
		log.Fatalf("failed to execute graphql operation, errors: %+v", r.Errors)
	}
	rJSON, _ := json.Marshal(r)

	//log.Printf("%s \n", rJSON) // {"data":{"hello":"world"}}
	if expectedJson != string(rJSON) {
		t.Errorf("want: %s got: %s", expectedJson, rJSON)
	}

}
