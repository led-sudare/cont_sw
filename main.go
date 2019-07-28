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
	// log "github.com/cihub/seelog"
)

type configs struct {
	port string `json:"port"`

	rsServer string `json:"rs_server"`
	demos    string `json:"demos"`
	iguchi   string `json:"iguchi"`
}

func newConfig() *configs {
	return &configs{
		port:     "8002",
		rsServer: "localhost:5002",
		demos:    "localhost:5003",
		iguchi:   "localhost:5004",
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
	util.ReadConfig(&conf)
	log.Info("dump: ", fmt.Sprintf("%#v", conf))

	e := createEcho()
	setLogger(e)
	setRouter(e)
	startEcho(e, conf.port)

	log.Info("done.")
}
