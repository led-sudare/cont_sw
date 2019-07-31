package main

import (
	"os"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	echoLog "github.com/labstack/gommon/log"
	"github.com/neko-neko/echo-logrus/v2/log"
	"github.com/sirupsen/logrus"
)

func initLogger() {
	log.Logger().SetOutput(os.Stdout)
	log.Logger().SetLevel(echoLog.DEBUG)
	log.Logger().SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})
}

func setLogger(e *echo.Echo) {
	e.Logger = log.Logger()
	e.Use(middleware.Logger())
	log.Info("Logger enabled.")
}
