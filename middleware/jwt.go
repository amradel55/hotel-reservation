package middleware

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthentication(c *fiber.Ctx) error {
	fmt.Println("-- JWT auth --")
	token, ok := c.GetReqHeaders()["X-Api-Token"]
	if !ok {
		return fmt.Errorf("Unauthorized")
	}
	claims, err := ValidateJWTToken(token[0])
	if err != nil {
		return err
	}
	expiresFloat := claims["expires"].(float64)
	expires := int64(expiresFloat)
	if time.Now().Unix() > expires {
		return fmt.Errorf("token expired")
	}
	return c.Next()
}

func ValidateJWTToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("invalid signing method", token.Header["alg"])
			return nil, fmt.Errorf("Unauthorized")
		}
		secret := "youCanDoIt"
		return []byte(secret), nil
	})
	if err != nil {
		fmt.Println("Failed to parse JWT token:", err)

		return nil, fmt.Errorf("Unauthorized")
	}
	if !token.Valid {
		fmt.Println("invalid token")
		return nil, fmt.Errorf("Unauthorized")

	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("Unauthorized")

	}
	return claims, nil
}
