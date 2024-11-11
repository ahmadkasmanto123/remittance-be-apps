package config

import (
	"encoding/json"
	"os"

	"github.com/gofiber/fiber/v2"
)

func GetConfig() fiber.Config {
	return fiber.Config{
		ErrorHandler: errorHandlerConfig,
		AppName:      os.Getenv("APP_NAME"),
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
	}
}
