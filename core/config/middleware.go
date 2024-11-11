package config

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/storage/sqlite3/v2"
)

func Middleware(appPublic *fiber.App) {
	appPublic.Use(Logger())
	appPublic.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, HEAD, PUT, PATCH, POST, DELETE",
	}))
	storage := sqlite3.New()
	appPublic.Use(limiter.New(limiter.Config{
		Max:        5,
		Expiration: 5 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			query := fmt.Sprintf("[%s|%s|%s|%s|%s]", c.Method(), c.IP(), c.Path(), c.Context().UserAgent(), c.GetReqHeaders()["Device-Id"])
			log.Print(query)
			return query
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(map[string]interface{}{
				"message": "Too many request reached, please try again later",
			})
		},
		SkipFailedRequests:     false,
		SkipSuccessfulRequests: false,
		Storage:                storage,
		LimiterMiddleware:      limiter.SlidingWindow{},
	}))

	//uncomment this for the better apps

	appPublic.Use(recover.New(recover.Config{
		Next:             nil,
		EnableStackTrace: false,
		StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
			log.Printf("%v", e)

			// panic("I'm an errors")
		},
	}))

}
