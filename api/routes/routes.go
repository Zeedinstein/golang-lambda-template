package routes

import (
	"github.com/GotBot-AI/users-api/api/controllers"
	"github.com/gofiber/fiber/v2"
)

func CreateRoutes(router fiber.Router) {
	router.Get("/users", controllers.GetAllUsers)
	router.Put("/users", controllers.CreateUser)
	router.Post("/users", controllers.UpdateUser)
	router.Get("/users/:userID", controllers.GetUserByID)
	router.Delete("/users/:userID", controllers.DeleteUser)
}
