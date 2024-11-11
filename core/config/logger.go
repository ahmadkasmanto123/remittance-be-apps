package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func Logger() func(*fiber.Ctx) error {
	file, err := os.OpenFile(fmt.Sprintf("%s%s ", os.Getenv("LOG_FOLDER"), time.Now().Format("January-2006.log")), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	// file, err := os.OpenFile(time.Now().Format("./January-2006.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	return logger.New(logger.Config{
		Format:     "===========================[ ${time} ]==========================\n::INFO:: ${pid} | IP : ${ip} | ${status} - ${method} ${host}${path}\n::PARAMS::\n${queryParams}\n::LATENCY::${latency}\n::REQUEST::\n${body}\n::RESPONSE::\n${resBody}\n::ERROR::\n${error}\n",
		TimeFormat: "02 January 2006 | 15:04:05",
		TimeZone:   "Asia/Jakarta",
		Output:     file,
	})

}

func StreamLog() {
	f, err := os.OpenFile(fmt.Sprintf("%s%s", os.Getenv("LOG_FOLDER"), time.Now().Format("January-2006.log")), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	// f, err := os.OpenFile(time.Now().Format("./January-2006.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return
	}
	log.SetOutput(f)
}
