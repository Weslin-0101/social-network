package routes

import (
	controllers "backend/src/controllers/user"
)

var userRoutes = []Route {
	{
		URI: "/users",
		Method: "POST",
		Function: controllers.CreateUser,
		AuthRequired: false,
	},
	{
		URI: "/users",
		Method: "GET",
		Function: controllers.GetAllUsers,
		AuthRequired: false,
	},
	{
		URI: "/users/{userID}",
		Method: "GET",
		Function: controllers.GetUserByID,
		AuthRequired: false,
	},
	{
		URI: "/users/nickname/{nickname}",
		Method: "GET",
		Function: controllers.GetUserByNickname,
		AuthRequired: false,
	},
	{
		URI: "/users/{userID}",
		Method: "PUT",
		Function: controllers.UpdateUserByID,
		AuthRequired: false,
	},
	{
		URI: "/users/{userID}",
		Method: "DELETE",
		Function: controllers.DeleteUserByID,
		AuthRequired: false,
	},
}