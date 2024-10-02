package api

import (
	"context"

	"github.com/amradel55/hotel-reservation/db"
	"github.com/gofiber/fiber/v2"
)

type USerHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *USerHandler {
	return &USerHandler{
		userStore: userStore,
	}
}

func (h *USerHandler) HandleGetUsers(c *fiber.Ctx) error {
	return c.JSON("use by Id ")

}

func (h *USerHandler) HandleGetUser(c *fiber.Ctx) error {
	var (
		id  = c.Params("id")
		ctx = context.Background()
	)

	user, err := h.userStore.GetUserByID(ctx, id)
	if err != nil {
		return err
	}
	return c.JSON(user)
}
