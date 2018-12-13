package main

import (
	"log"
	"os"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/99designs/gqlgen/handler"
	"github.com/Go-GraphQL-Group/GraphQL-Service"
	"github.com/Go-GraphQL-Group/GraphQL-Service/server/service"
)

const defaultPort = "8080"

func NewServer() *negroni.Negroni {
	router := mux.NewRouter()
	initRoutes(router)

	n := negroni.Classic()

	n.UseHandler(router)
	return n
}

func initRoutes(router *mux.Router) {
	router.HandleFunc("/login", service.LoginHandler).Methods("POST")
	router.HandleFunc("/", handler.Playground("GraphQL playground", "/query"))
	router.HandleFunc("/query", handler.GraphQL(GraphQL_Service.NewExecutableSchema(GraphQL_Service.Config{Resolvers: &GraphQL_Service.Resolver{}})))
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	server := NewServer()
	server.Run(":"+port)
	
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
}
