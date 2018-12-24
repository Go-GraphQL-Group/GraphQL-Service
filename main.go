package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/handler"
	"github.com/Go-GraphQL-Group/GraphQL-Service/graphql"
	"github.com/Go-GraphQL-Group/GraphQL-Service/resolver"
	"github.com/Go-GraphQL-Group/GraphQL-Service/service"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

const defaultPort = "9090"

func NewServer() *negroni.Negroni {
	router := mux.NewRouter()
	initRoutes(router)

	n := negroni.Classic()
	n.Use(negroni.NewStatic(http.Dir("./static")))
	n.UseHandler(router)
	return n
}

func initRoutes(router *mux.Router) {
	router.HandleFunc("/api/login", service.LoginHandler).Methods("POST")
	router.Use(service.TokenMiddleware)
	router.HandleFunc("/api", service.ApiHandler).Methods("GET")
	// router.HandleFunc("/", handler.Playground("GraphQL playground", "/api/query"))
	router.HandleFunc("/api/query", handler.GraphQL(graphql.NewExecutableSchema(graphql.Config{Resolvers: &resolver.Resolver{}})))
	router.HandleFunc("/api/logout", service.LogoutHandler).Methods("POST", "GET")
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	server := NewServer()
	server.Run(":" + port)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
}
