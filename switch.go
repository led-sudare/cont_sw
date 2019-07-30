package main

import (
	"cont_sw/lib/api"
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/neko-neko/echo-logrus/v2/log"
)

type contentSwitch struct {
	nodes []contentNode
}

type contentNode struct {
	prop   Content
	client *api.Node
}

func StartContentSwitch(conf Configs) error {

	nodes, err := createNodes(&conf)
	if err != nil {
		log.Errorf("create nodes: %s", err)
		return err
	}

	cs := contentSwitch{
		nodes: nodes,
	}
	cs.disableNodes()

	e := createEcho()
	startEcho(e, conf.Port)

	return nil
}

func createEcho() *echo.Echo {
	e := echo.New()
	e.HideBanner = true

	setLogger(e)
	setRouter(e)
	return e
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

func createNodes(conf *Configs) ([]contentNode, error) {
	nodes := make([]contentNode, 0)
	for _, content := range conf.Contents {
		node, err := api.NewNode(content.Addr, api.Timeout(conf.Timeout))
		if err != nil {
			return nil, err
		}

		nodes = append(nodes, contentNode{
			prop:   content,
			client: node,
		})
	}
	return nodes, nil
}

func changeStateNode(node contentNode, enable bool) error {
	// check alive & status
	stat, err := node.client.GetStatus(nil)
	if err != nil {
		log.Warn(fmt.Sprintf("node is dead.: %s\n%s", node.prop.Name, err))
		return err
	}

	for stat.Enable != enable {
		// try change status
		node.client.PostConfig(nil, enable)
		stat, err = node.client.GetStatus(nil)
		if err != nil {
			log.Warn(fmt.Sprintf("node is dead.: %s\n%s", node.prop.Name, err))
			break
		}
	}
	log.Infof("changed status (%s): %t", node.prop.Name, enable)
	return nil
}

func (cs *contentSwitch) enableNodes(name string) error {
	for _, node := range cs.nodes {
		if node.prop.Name == name {
			go changeStateNode(node, true)
		} else {
			go changeStateNode(node, false)
		}
	}
	return nil
}

func (cs *contentSwitch) disableNodes() {
	for _, node := range cs.nodes {
		changeStateNode(node, false)
	}
}
