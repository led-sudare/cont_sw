package main

import (
	"cont_sw/lib/util"
	"net/http"

	"fmt"
	"os"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	echoLog "github.com/labstack/gommon/log"
	"github.com/neko-neko/echo-logrus/v2/log"
	"github.com/sirupsen/logrus"
)

type Configs struct {
	port string `json:"port"`
	// contents []Content `json:"contents`
}

type Content struct {
	name   string `json:"name"`
	addr   string `json:"addr"`
	enable bool   `json:"enable`
}

func newConfig() *Configs {
	return &Configs{
		port: "8002",
		// contents: make([]Content, 0),
	}
}

func createEcho() *echo.Echo {
	e := echo.New()
	e.HideBanner = true

	return e
}

func initLogger() {
	log.Logger().SetOutput(os.Stdout)
	log.Logger().SetLevel(echoLog.INFO)
	log.Logger().SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})

}

func setLogger(e *echo.Echo) {
	e.Logger = log.Logger()
	e.Use(middleware.Logger())
	log.Info("Logger enabled.")
}

func startEcho(e *echo.Echo, port string) {
	e.Logger.Fatal(e.Start(":" + port))
}

func setRouter(e *echo.Echo) {
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World!\n")
		// TODO WebUI
	})
}

func main() {
	initLogger()

	log.Info("start")
	conf := newConfig()
	util.ReadConfig(conf)
	log.Info("dump: ", fmt.Sprintf("%#v", conf))

	e := createEcho()
	setLogger(e)
	setRouter(e)
	startEcho(e, conf.port)

	log.Info("done.")
}
