package api

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/amradel55/hotel-reservation/db"
	"github.com/amradel55/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthHandler struct {
	userStore db.UserStore
}

func NewAuthHandler(userStore db.UserStore) *AuthHandler {
	return &AuthHandler{
		userStore: userStore,
	}
}

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	User  *types.User `json:"user"`
	Token string      `json:"token"`
}

type genericResp struct {
	Status string `json:"status"`
	Msg    string `json:"msg"`
}

func (h *AuthHandler) HandleAuthenticate(c *fiber.Ctx) error {
	var params AuthParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	user, err := h.userStore.GetUserByEmail(c.Context(), params.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.Status(http.StatusBadRequest).JSON(genericResp{
				Status: "error",
				Msg:    "invalid Credentials",
			})
		}
		return err
	}
	if !types.IsValidPassword(user.EncryptedPassword, params.Password) {
		return c.Status(http.StatusBadRequest).JSON(genericResp{
			Status: "error",
			Msg:    "invalid Credentials",
		})
	}
	res := AuthResponse{
		User:  user,
		Token: createTokenFromUser(user),
	}
	return c.JSON(res)
}

func createTokenFromUser(user *types.User) string {
	now := time.Now()
	expires := now.Add(time.Hour * 4).Unix()
	claims := jwt.MapClaims{
		"id":      user.ID,
		"email":   user.Email,
		"expires": expires,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := "youCanDoIt" // you can apply the secret in different ways one is to fetch it through the os.Getenv or have a config file returns it from .env file
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Println("failed to sign token with secret", err)
	}
	return tokenStr
}
