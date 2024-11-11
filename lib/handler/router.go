package handler

import (
	"love-remittance-be-apps/core/config"

	"github.com/gofiber/fiber/v2"
)

func Router(app *fiber.App) {
	config.Middleware(app)

	//public
	calculatorHandler(app.Group("/calc"))
	loginHandler(app.Group("/apps"))
	registerHandler(app.Group("/reg"))

	//private
	domestic := app.Group("/domestic")
	domestic.Use(config.JwtConfig())

	domesticHandler(domestic.Group("/calc"))
	domTrxHandler(domestic.Group("/trx"))

	update := app.Group("/update")
	update.Use(config.JwtConfig())

	updateHandler(update.Group("/account"))

	getData := app.Group("/get")
	getData.Use(config.JwtConfig())

	//get data
	getDataHandler(getData.Group("/data"))

	international := app.Group("/international")
	international.Use(config.JwtConfig())

	internationalHandler(international.Group("/calc"))
	intlTrxHandler(international.Group("/trx"))
}
