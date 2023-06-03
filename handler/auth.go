package handler

import (
	"auth-server/config"
	"auth-server/database"
	"auth-server/model"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

// Util Methods
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	fmt.Println("err = ", err, "ref:", password, " - ", hash)
	return err == nil
}

// Request Schemas
type LoginRequest struct {
	Name string `json:"username" xml:"username" form:"username"`
	Pass string `json:"password" xml:"password" form:"password"`
}

type SignupRequest struct {
	Name     string `json:"username" xml:"username" form:"username"`
	Pass     string `json:"password" xml:"password" form:"password"`
	Email    string `json:"email" xml:"email" form:"email"`
	FullName string `json:"full_name" xml:"full_name" form:"full_name"`
}

// Handlers
func Login(c *fiber.Ctx) error {
	data := new(LoginRequest)
	if err := c.BodyParser(data); err != nil {
		return err
	}

	// Get User from DB
	db := database.DB
	user_model := new(model.User)
	if err := db.First(&user_model, "username = ?", data.Name).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Invalid username"})
	}

	if !CheckPasswordHash(data.Pass, user_model.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Invalid password"})
	}

	// Create the Claims
	claims := jwt.MapClaims{
		"username": user_model.Username,
		"is_admin": user_model.IsAdmin,
		"expiry":   time.Now().Add(time.Hour * 72).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(config.Config("SECRET")))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"token": t})
}

func Signup(c *fiber.Ctx) error {
	data := new(SignupRequest)
	if err := c.BodyParser(data); err != nil {
		// return err
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	// Check that user with same username is not already present
	db := database.DB
	user_model_check := new(model.User)
	if err := db.First(&user_model_check, "username = ?", data.Name).Error; err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "User with given user_name already exists. please use different username."})
	}

	pass, _ := hashPassword(data.Pass)
	user_model := model.User{
		Username: data.Name,
		Password: pass,
		Email:    data.Email,
		FullName: data.FullName,
		IsAdmin:  false,
	}
	if err := db.Create(&user_model).Error; err != nil {
		return c.JSON(fiber.Map{"message": "error during user creation", "data": new(struct{})})
	}

	return c.JSON(fiber.Map{"message": "user registered successfully", "data": user_model})
}

func HandleRoot(c *fiber.Ctx) error {
	fmt.Print("p0")
	user := c.Locals("user").(*jwt.Token)
	fmt.Print("p1")
	claims := user.Claims.(jwt.MapClaims)
	fmt.Print(claims)
	fmt.Print("p2")
	username, ok := claims["username"].(string)
	fmt.Print("p3", username, " - ", ok)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "malformed JWT"})
	}
	// Get User from DB
	db := database.DB
	user_model := new(model.User)
	if err := db.First(&user_model, "username = ?", username).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Unexpected error in retrival of profile"})
	}
	return c.JSON(fiber.Map{"data": user_model})
}
