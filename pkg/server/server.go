package server

import (
	"net/http"
	"shop/pkg/generated"
	"shop/pkg/resolver"
	"shop/pkg/storage"
	"shop/pkg/user"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/99designs/gqlgen/handler"
)

func GetApp() (http.Handler, error) {

	mux := http.NewServeMux()

	st := storage.NewStorageInit()
	err := storage.AddInitialData("./data/testdata.json", st)
	if err != nil {
		return nil, err
	}

	userHandler := user.NewUserHandler()

	cfg := generated.Config{
		Resolvers: &resolver.Resolver{
			StorageData:   st,
			UserIdFromCtx: userHandler.Sm.GetUserIdFromCtx,
		},
		Directives: generated.DirectiveRoot{
			Authorized: userHandler.Sm.CheckSession,
		},
	}

	gqlHandler := handler.GraphQL(generated.NewExecutableSchema(cfg))
	mux.Handle("/", playground.Handler("GraphQL playground", "/query"))

	mux.Handle("/query", userHandler.Sm.AuthMiddleWare(gqlHandler))
	mux.HandleFunc("/register", userHandler.Register)

	return mux, nil
}
