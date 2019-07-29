package main

import (
	"cont_sw/lib/api"
	"cont_sw/lib/util"

	"net/http"
	"os"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

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
	Enable bool   `yaml:"enable"`
}

func newConfig() Configs {
	return Configs{
		Port:     "8080",
		Contents: make([]Content, 0),
	}
}

type ContentNode struct {
	prop   *Content
	client *api.Node
}

func createEcho() *echo.Echo {
	e := echo.New()
	e.HideBanner = true

	return e
}

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

func startEcho(e *echo.Echo, port string) {
	e.Logger.Fatal(e.Start(":" + port))
}

func setRouter(e *echo.Echo) {
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World!\n")
		// TODO WebUI
	})
}

func createNodes(conf *Configs) ([]ContentNode, error) {
	nodes := make([]ContentNode, 0)
	for _, content := range conf.Contents {
		node, err := api.NewNode(content.Addr)
		if err != nil {
			return nil, err
		}

		nodes = append(nodes, ContentNode{
			prop:   &content,
			client: node,
		})
	}
	return nodes, nil
}

func main() {
	initLogger()

	log.Info("start contents switch !!")
	conf := newConfig()
	err := util.ReadConfig(&conf)
	if err != nil {
		log.Errorf("read configs: %s", err)
		os.Exit(-1)
	}

	nodes, err := createNodes(&conf)
	if err != nil {
		log.Errorf("create nodes: %s", err)
		os.Exit(-1)
	}

	err = nodes[0].client.PostConfig(nil, false)
	if err != nil {
		log.Errorf("response: %s", err)
		os.Exit(-1)
	}
	stat, err := nodes[0].client.GetStatus(nil)
	if err != nil {
		log.Errorf("response: %s", err)
		os.Exit(-1)
	}
	log.Infof("demos status: %#v", stat)

	e := createEcho()
	setLogger(e)
	setRouter(e)
	startEcho(e, conf.Port)

	log.Info("done.")
}
