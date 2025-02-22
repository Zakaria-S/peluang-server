package main

import (
	"log"
	"peluang-server/internal/component"
	"peluang-server/internal/config"
	"peluang-server/internal/modules/otp"
	"peluang-server/internal/modules/user"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	conf := config.NewConfig()

	db := component.GetDatabaseConnection(conf)

	component.Migrate(db)

	userRepository := user.NewRepository(db)
	userOtpRepository := otp.NewRepository(db)

	userService := user.NewService(userRepository, userOtpRepository)

	app := fiber.New()
	app.Use(cors.New())
	app.Use(logger.New())
	app.Use(recover.New())

	user.NewRoute(app, userService)

	log.Fatal(app.Listen(":" + conf.Srv.Port))
}
