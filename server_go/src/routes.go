package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func Router() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	return router
}

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"TodoIndex",
		"GET",
		"/todos",
		TodoIndex,
	},
	Route{
		"TodoShow",
		"GET",
		"/todos/{todoId}/",
		TodoShow,
	},
	Route{
		"PostTodo",
		"POST",
		"/todo",
		PostTodo,
	},
	Route{
		"postTransaction",
		"POST",
		"/transaction",
		postTransaction,
	},
	Route{
		"getAllTransaction",
		"GET",
		"/transactions",
		getAllTransaction,
	},
	Route{
		"getBalances",
		"GET",
		"/balances",
		getBalances,
	},
}
