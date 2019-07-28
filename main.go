package main

import (
	"io/ioutil"
	"net/http"

	"fmt"
	"os"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gopkg.in/yaml.v2"

	echoLog "github.com/labstack/gommon/log"
	"github.com/neko-neko/echo-logrus/v2/log"
	"github.com/sirupsen/logrus"
)

type Configs struct {
	Port     string    `yaml:"port"`
	Contents []Content `yaml:"contents`
}

type Content struct {
	Name   string `yaml:"name"`
	Addr   string `yaml:"addr"`
	Enable string `yaml:"enable`
}

func newConfig() Configs {
	return Configs{
		Port:     "8080",
		Contents: make([]Content, 0),
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

	var conf Configs
	buf, err := ioutil.ReadFile("config.yml")
	log.Info("read: ", string(buf))

	err = yaml.Unmarshal(buf, &conf)
	if err != nil {
		panic(err)
	}
	// util.ReadConfig(&conf)
	log.Info("dump: ", fmt.Sprintf("%#v", conf))

	e := createEcho()
	setLogger(e)
	setRouter(e)
	startEcho(e, conf.Port)

	log.Info("done.")
}
