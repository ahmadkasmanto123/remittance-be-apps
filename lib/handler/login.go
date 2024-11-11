package handler

import (
	"love-remittance-be-apps/lib/repository"
	"love-remittance-be-apps/lib/service"

	"github.com/gofiber/fiber/v2"
)

func loginHandler(api fiber.Router) {
	loginRepository := repository.NewLoginRepository()
	loginService := service.NewLoginService(loginRepository)

	api.Post("/login", handler(loginService.Login))
	api.Post("/createOtp", handler(loginService.LogInOtpCreate))
	api.Post("/validateOtp", handler(loginService.LogInOtpValidate))
}
