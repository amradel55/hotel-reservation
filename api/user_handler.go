package api

import (
	"github.com/amradel55/hotel-reservation/db"
	"github.com/amradel55/hotel-reservation/types"
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

func (h *USerHandler) HandlePostUser(c *fiber.Ctx) error {
	var params types.CreateUserParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	if err := params.Validate(); len(err) > 0 {
		return c.JSON(err)
	}
	user, err := types.NewUserFromParams(params)
	if err != nil {
		return err
	}
	insertedUser, err := h.userStore.InsertUser(c.UserContext(), user)
	if err != nil {
		return err
	}
	return c.JSON(insertedUser)
}

func (h *USerHandler) HandleGetUsers(c *fiber.Ctx) error {
	users, err := h.userStore.GetUsers(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(users)
}

func (h *USerHandler) HandleGetUser(c *fiber.Ctx) error {
	var (
		id = c.Params("id")
	)

	user, err := h.userStore.GetUserByID(c.Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(user)
}
