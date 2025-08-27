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
		AuthRequired: true,
	},
	{
		URI: "/users/{userID}",
		Method: "GET",
		Function: controllers.GetUserByID,
		AuthRequired: true,
	},
	{
		URI: "/users/nickname/{nickname}",
		Method: "GET",
		Function: controllers.GetUserByNickname,
		AuthRequired: true,
	},
	{
		URI: "/users/{userID}",
		Method: "PUT",
		Function: controllers.UpdateUserByID,
		AuthRequired: true,
	},
	{
		URI: "/users/{userID}",
		Method: "DELETE",
		Function: controllers.DeleteUserByID,
		AuthRequired: true,
	},
}