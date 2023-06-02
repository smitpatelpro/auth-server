package handler

import (
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

// Validation
// type ErrorResponse struct {
// 	FailedField string
// 	Tag         string
// 	Value       string
// }

// var validate = validator.New()

// func ValidateStruct(user User) []*ErrorResponse {
// 	var errors []*ErrorResponse
// 	err := validate.Struct(user)
// 	if err != nil {
// 		for _, err := range err.(validator.ValidationErrors) {
// 			var element ErrorResponse
// 			element.FailedField = err.StructNamespace()
// 			element.Tag = err.Tag()
// 			element.Value = err.Param()
// 			errors = append(errors, &element)
// 		}
// 	}
// 	return errors
// }

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

func Signup(c *fiber.Ctx) error {
	data := new(SignupRequest)
	if err := c.BodyParser(data); err != nil {
		return err
	}

	// Check that user with same username is not already present
	db := database.DB
	user_model_check := new(model.User)
	if err := db.First(&user_model_check, "username = ?", data.Name).Error; err == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "User with given user_name already exists. please use different username."})
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
	return c.SendString("Hello World")
}
