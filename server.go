package main

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/99designs/gqlgen/handler"
)

func GetApp() http.Handler {

	mux := http.NewServeMux()

	storage := NewStorageInit()

	err := storage.AddInitialData("testdata.json")
	if err != nil {
		log.Println("error with initial data json:", err)
	}

	userHandler := NewUserHandler()

	cfg := Config{
		Resolvers: &Resolver{
			StorageData:   storage,
			UserIdFromCtx: userHandler.Sm.GetUserIdFromCtx,
		},
		Directives: DirectiveRoot{
			Authorized: userHandler.Sm.CheckSession,
		},
	}

	gqlHandler := handler.GraphQL(NewExecutableSchema(cfg))
	mux.Handle("/", playground.Handler("GraphQL playground", "/query"))

	mux.Handle("/query", userHandler.Sm.AuthMiddleWare(gqlHandler))
	mux.HandleFunc("/register", userHandler.Register)
	return mux
}
