package api

import "github.com/gofiber/fiber/v2"

func HandleGetUsers(c *fiber.Ctx) error {
	return c.JSON("James")
}

func HandleGetUser(c *fiber.Ctx) error {
	return c.JSON("use by Id ")
}
