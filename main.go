package main

import (
	"github.com/go-pg/pg/v10"
	"github.com/joho/godotenv"
	"log"
	"market/auth"
	"market/graph"
	"market/graph/generated"
	"market/repositories"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	DB := repositories.New(&pg.Options{
		Addr:     os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		Database: "market",
	})

	defer func(DB *pg.DB) {
		err := DB.Close()
		if err != nil {
			log.Println("Failed to close db connection")
		}
	}(DB)

	usersRepo := repositories.NewUsers(DB)
	productsRepo := repositories.NewProducts(DB)

	srv := handler.NewDefaultServer(
		generated.NewExecutableSchema(
			generated.Config{
				Resolvers: &graph.Resolver{
					Auth:         auth.New(usersRepo),
					ProductsRepo: productsRepo,
				},
			},
		),
	)

	mux := http.NewServeMux()

	mux.Handle("/", playground.Handler("GraphQL playground", "/query"))
	mux.Handle("/query", auth.CurrentUserMiddleware(usersRepo)(srv))

	log.Println("Listening", os.Getenv("HTTP_PORT"))
	log.Fatal(http.ListenAndServe(":"+os.Getenv("HTTP_PORT"), mux))
}
