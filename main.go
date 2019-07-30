package main

import (
	"cont_sw/lib/util"

	"os"

	"github.com/neko-neko/echo-logrus/v2/log"
)

type Configs struct {
	Port    string `yaml:"port"`
	Timeout int    `yaml:"timeout"`

	Contents []Content `yaml:"contents"`
}

type Content struct {
	Name   string `yaml:"name"`
	Addr   string `yaml:"addr"`
	Enable bool   `yaml:"enable"`
}

func newConfig() Configs {
	return Configs{
		Port:     "8080",
		Timeout:  5,
		Contents: make([]Content, 0),
	}
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

	err = StartContentSwitch(conf)
	if err != nil {
		log.Info("fatal switch.")
		os.Exit(-1)
	}
	log.Info("finished switch.")
}
