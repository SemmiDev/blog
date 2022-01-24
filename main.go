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
	// init the configurations
	LoadConfig(".")

	// init the logger
	zerolog.Log = zerolog.NewConsole(false)

	// init the database pooling
	dbPool, err := pgxpool.Connect(context.Background(), Config.DBSource)
	if err != nil {
		zerolog.Log.Error().Interface("db connection", err).Send()
	}
	defer dbPool.Close()

	// init the token manager
	tokenMaker, err := token.NewPasetoMaker(Config.TokenSymmetricKey)
	if err != nil {
		zerolog.Log.Error().Interface("token maker", err).Send()
	}

	// init the auth server
	authServer, err := userserver.NewAuthServer(dbPool, tokenMaker)
	if err != nil {
		zerolog.Log.Error().Interface("auth server", err).Send()
	}

	// init the user server
	userServer, err := userserver.NewUserServer(dbPool, tokenMaker)
	if err != nil {
		zerolog.Log.Error().Interface("user server", err).Send()
	}

	// init the fiber app
	app := fiber.New(
		fiber.Config{
			ReadTimeout:  time.Second * 5,
			WriteTimeout: time.Second * 5,
		},
	)

	// init the middlewares
	app.Use(fLog.New())
	app.Use(recover.New())
	app.Use(cors.New())

	// init the auth routes
	authGroup := app.Group("/auth")
	authServer.Mount(authGroup)

	// init the user routes
	userGroup := app.Group("/users")
	userServer.Mount(userGroup)

	// start the app
	log.Fatal(app.Listen(Config.ServerAddress))
}
