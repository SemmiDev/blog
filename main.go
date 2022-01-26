package main

import (
	"context"
	. "github.com/SemmiDev/blog/config"
	zerolog "github.com/SemmiDev/blog/internal/common/logger"
	userserver "github.com/SemmiDev/blog/internal/user/server"
	"github.com/SemmiDev/blog/internal/user/token"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fLog "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
	"time"
)

func main() {
	// set up the configurations.
	LoadConfig(".")

	// set up the logger.
	zerolog.Log = zerolog.NewConsole(false)

	// set up the database.
	dbPool, err := pgxpool.Connect(context.Background(), Env.DBSource)
	if err != nil {
		zerolog.Log.Error().Interface("db connection", err).Send()
	}
	defer dbPool.Close()

	// set up the token manager.
	tokenMaker, err := token.NewPasetoMaker(Env.TokenSymmetricKey)
	if err != nil {
		zerolog.Log.Error().Interface("token maker", err).Send()
	}

	// set up the auth server.
	authServer, err := userserver.NewAuthServer(dbPool, tokenMaker)
	if err != nil {
		zerolog.Log.Error().Interface("auth server", err).Send()
	}

	// set up the user server.
	userServer, err := userserver.NewUserServer(dbPool, tokenMaker)
	if err != nil {
		zerolog.Log.Error().Interface("user server", err).Send()
	}

	// set up the fiber app.
	app := fiber.New(
		fiber.Config{
			ReadTimeout:  time.Second * 5,
			WriteTimeout: time.Second * 5,
		},
	)

	// set up the middlewares.
	app.Use(fLog.New())
	app.Use(recover.New())
	app.Use(cors.New())

	// set up the auth routes.
	authGroup := app.Group("/auth")
	authServer.Mount(authGroup)

	// set up the user routes.
	userGroup := app.Group("/users")
	userServer.Mount(userGroup)

	// start the app on the server address port.
	log.Fatal(app.Listen(Env.ServerAddress))
}
