package main

import (
	"log"
	"love-remittance-be-apps/core/config"
	"love-remittance-be-apps/lib/handler"
	"os"
	"runtime"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	// load
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("err loading: %v", err)
	}
	config.StreamLog()
	// get
	addr := os.Getenv("SESSION_URL")
	if addr == "" {
		log.Fatalf("missing ssion: %v", err)
	}
	runtime.GC()
	// fmt.Printf("%s%s", os.Getenv("LOG_FOLDER"), time.Now().Format("January-2006.log"))
	app := fiber.New(config.GetConfig())
	handler.Router(app)
	app.ListenTLS(":"+os.Getenv("PORT"), "./host.pem", "./host.key")
}
