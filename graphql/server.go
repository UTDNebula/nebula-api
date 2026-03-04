package main

import (
	"graphql/configs"
	"graphql/graph"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/vektah/gqlparser/v2/ast"
)

func main() {
	port := configs.GetPortString()
	resolver := graph.Resolver{
		CourseCollection:   configs.GetCollection("courses"),
		SectionCollection:  configs.GetCollection("sections"),
		ProfCollection:     configs.GetCollection("professors"),
		BuildingCollection: configs.GetCollection("rooms"),
	}

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &resolver}))
	srv.Use(extension.FixedComplexityLimit(100)) // Avoid unlimited nesting (later)

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
