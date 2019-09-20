package main

import (
	"log"

	"github.com/gorilla/mux"
	"github.com/joelmdesouza/prest-admin/handlers"
	"github.com/prest/cmd"
	"github.com/prest/config"
	"github.com/prest/config/router"
	"github.com/prest/middlewares"
	"github.com/urfave/negroni"
)

func main() {
	config.Load()
	err := handlers.Load()
	if err != nil {
		log.Panic(err)
	}
	registerAllMiddlewares()

	r := router.Get()

	padminRouter := mux.NewRouter().PathPrefix("/padmin").Subrouter().StrictSlash(true)
	padminRouter.HandleFunc("/list/{table}", handlers.ListHandler).Methods("GET")
	padminRouter.HandleFunc("/create/{table}", handlers.CreateHandler).Methods("GET")
	padminRouter.HandleFunc("/create/{table}", handlers.CreateHandlerPost).Methods("POST")
	padminRouter.HandleFunc("/edit/{table}/{key}", handlers.EditHandler).Methods("GET")
	r.PathPrefix("/padmin").Handler(negroni.New(
		negroni.Wrap(padminRouter),
	))
	//r.HandleFunc("/list", handlers.ListHandler)
	//r.HandleFunc("/create", handlers.CreateHandler)
	//r.HandleFunc("/edit", handlers.EditHandler)
	//r.HandleFunc("/test", re).Methods("POST")
	//r.HandleFunc("/test/{id}", re).Methods("PUT")

	cmd.Execute()
}

func registerAllMiddlewares() {
	middlewares.MiddlewareStack = append(middlewares.MiddlewareStack, negroni.NewRecovery())
	middlewares.MiddlewareStack = append(middlewares.MiddlewareStack, negroni.NewLogger())
	//middlewares.MiddlewareStack = append(middlewares.MiddlewareStack, middlewares.HandlerSet())
}
