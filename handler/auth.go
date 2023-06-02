package handler

import (
	"auth-server/database"
	"auth-server/model"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type User struct {
	Name string `json:"user" xml:"user" form:"user"`
	Pass string `json:"pass" xml:"pass" form:"pass"`
}

func Login(c *fiber.Ctx) error {
	u := new(User)

	if err := c.BodyParser(u); err != nil {
		return err
	}

	db := database.DB
	user_model := new(model.User)

	if err := db.First(&user_model, "username = ?", u.Name).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Invalid credentials"})
	}

	// Create the Claims
	claims := jwt.MapClaims{
		"user":  user_model.Username,
		"name":  user_model.FullName,
		"admin": user_model.IsAdmin,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"token": t})
}

func HandleRoot(c *fiber.Ctx) error {
	return c.SendString("Hello World")
}
