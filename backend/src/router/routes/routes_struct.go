package routes

import (
	"backend/src/middlewares"
	"net/http"

	"github.com/gorilla/mux"
)

type Route struct {
	URI 			string
	Method 			string
	Function 		func(http.ResponseWriter, *http.Request)
	AuthRequired 	bool 
}

func Config(r *mux.Router) *mux.Router {
	routes := userRoutes
	routes = append(routes, loginRoutes)

	for _, route := range routes {

		if route.AuthRequired {
			r.HandleFunc(
				route.URI,
				middlewares.Logger(
					middlewares.Authenticate(route.Function),
				)).Methods(route.Method)
		} else {
			r.HandleFunc(route.URI, middlewares.Logger(route.Function)).Methods(route.Method)
		}
	}

	return r
}