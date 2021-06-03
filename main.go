package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/patrickmn/go-cache"

	"github.com/friendsofgo/graphiql"
	"github.com/graphql-go/graphql"
	gqlhandler "github.com/graphql-go/graphql-go-handler"
)

const HOST = "https//www.joshua.heroku.com/"

var db *sql.DB

var server = "localhost"
var port = 1433
var user = "sa"
var password = "AdaJosh@2019"
var database = "report_db"

var businessLogic BusinessLogic
var dataStore Datastore

func main() {

	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: ShortenURLSchema}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery),
		Mutation: mutationType,
	}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}

	// Build connection string
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		server, user, password, port, database)

	// Create connection pool
	db, err = sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Error creating connection pool: ", err.Error())
	}
	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("Connected!\n")

	// Instantiate datastore
	dataStore = Datastore{store: db}

	// Instatiate cache
	// Create a cache with a default expiration time of 5 minutes, and which
	// purges expired items every 10 minutes
	c := cache.New(5*time.Minute, 10*time.Minute)

	// instantiate business logic
	businessLogic = BusinessLogic{db: &dataStore, c: *c}

	h := gqlhandler.New(&gqlhandler.Config{
		Schema: &schema,
		Pretty: true,
	})

	// Serve a GraphQL endpoint at `/graphql`
	http.Handle("/graphql", h)

	graphiqlHandler, err := graphiql.NewGraphiqlHandler("/graphql")
	if err != nil {
		panic(err)
	}
	http.Handle("/ui", graphiqlHandler)
	http.HandleFunc("/", handleGetShortURL)
	http.ListenAndServe(":8080", nil)
}

// Handles searching of shortkeycode and redirects to long url if a match is found
func handleGetShortURL(w http.ResponseWriter, req *http.Request) {
	shortURLCode := strings.TrimPrefix(req.URL.Path, "/")

	if shortURLCode == "" {
		fmt.Fprintf(w, "url code is missing!")
		return
	}

	log.Printf("INFO: received shortURL code: %s", shortURLCode)

	// Retrieve mapped url code
	longURL, err := businessLogic.search(shortURLCode)
	if err != nil {
		fmt.Fprintf(w, "unable to retrieve code")
		return
	}

	http.Redirect(w, req, longURL, http.StatusSeeOther)

}
