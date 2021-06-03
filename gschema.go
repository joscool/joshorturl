package main

import (
	"fmt"
	"log"

	"github.com/graphql-go/graphql"
)

// Schema
var ShortenURLSchema = graphql.Fields{
	"hello": &graphql.Field{
		Type: graphql.String,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			return "whatever the mind can conceive it can achieve", nil
		},
	},
}

// ShortenURL mutation definition
var mutationType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		"shortenURL": &graphql.Field{
			Type:        graphql.String,
			Description: "Create new shortened url",
			Args: graphql.FieldConfigArgument{
				"url": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				url, _ := params.Args["url"].(string)
				shortURL, err := businessLogic.GenerateShortURLCode(url)
				if err != nil {
					log.Printf("ERROR: %s", err)
					return "uable to generate short url. please try again", err
				}
				returnURL := fmt.Sprintf("%s%s", HOST, shortURL)
				return returnURL, nil
			},
		},
	},
})
