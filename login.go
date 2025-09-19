package main

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type User struct {
	Realname string `json:"realname"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

var ctx = context.Background()
var rdb *redis.Client

func main() {
	// menyambungkan ke redis
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})

	app := fiber.New()

	// endpoint login
	app.Post("/login", loginHandler)

	fmt.Println("Server running on http://localhost:3000")
	log.Fatal(app.Listen(":3000"))
}

func loginHandler(c *fiber.Ctx) error {
	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var req LoginRequest

	if err := c.BodyParser(&req); err != nil {
		// request tidak valid
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "invalid request body",
		})
	}

	key := "login:" + req.Username
	val, err := rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  false,
			"message": "user not found",
		})
	} else if err != nil {
		log.Printf("redis get error for key %s: %v\n", key, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "internal server error",
		})
	}

	var user User
	if err := json.Unmarshal([]byte(val), &user); err != nil {
		log.Printf("error parsing user JSON for key %s: %v\n", key, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "error parsing user data",
		})
	}

	h := sha1.New()
	h.Write([]byte(req.Password))
	hashedInput := hex.EncodeToString(h.Sum(nil))

	if hashedInput != user.Password {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  false,
			"message": "invalid password",
		})
	}

	return c.JSON(fiber.Map{
		"status":   true,
		"message":  "login success",
		"realname": user.Realname,
		"email":    user.Email,
	})
}
