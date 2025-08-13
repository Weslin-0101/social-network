package routes

import (
	controllers "backend/src/controllers/login"
)

var loginRoutes = Route {
	URI: 			"/login",
	Method: 		"POST",
	Function: 		controllers.Login,
	AuthRequired: 	false,
}