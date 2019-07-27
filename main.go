package main

import (
	"cont_sw/lib/util"

	"net/http"
	"fmt"

	// "github.com/labstack/echo"
	// "github.com/labstack/echo/middleware"

	log "github.com/cihub/seelog"
)

type Configs struct {
	rsServer string `json:"rs_server"`
	demos    string `json:"demos"`
	iguchi   string `json:"iguchi"`
}

func NewConfig() Configs {
	return Configs{
		rsServer: "localhost:5002",
		demos:    "localhost:5003",
		iguchi:   "localhost:5004",
	}
}

func main() {
	log.Info("start")
	conf := NewConfig()
	util.ReadConfig(&conf)
	log.Info("dump: ", fmt.Sprintf("%#v", conf))

	// e := echo.New()
	// e.HideBanner = true

	log.Info("done.")
	log.Flush()
}
