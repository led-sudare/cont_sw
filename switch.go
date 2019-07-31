package main

import (
	"cont_sw/lib/api"
	. "cont_sw/lib/model"

	"fmt"
	"net/http"
	"strings"

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

func createEcho(sw *contentSwitch) *echo.Echo {
	e := echo.New()
	e.HideBanner = true

	setLogger(e)
	setRouter(e, sw)
	return e
}

func startEcho(e *echo.Echo, port string) {
	e.Logger.Fatal(e.Start(":" + port))
}

func getName(c echo.Context) (string, error) {
	name := c.QueryParam("name")
	log.Debugf("get request params: name=%s\n", name)
	return name, nil
}

func getStatus(c echo.Context) (*Status, error) {
	s := &Status{}
	if err := c.Bind(s); err != nil {
		log.Error(err)
		return nil, fmt.Errorf("invalid body parameter.")
	}
	return s, nil
}

func setRouter(e *echo.Echo, sw *contentSwitch) {
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World!\n")
		// TODO WebUI
	})

	e.GET("/api/status", func(c echo.Context) error {
		name, err := getName(c)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		var stat Status
		stat.Enable, err = sw.getStatus(name)
		return c.JSON(http.StatusOK, stat)
	})

	e.POST("/api/enable", func(c echo.Context) error {
		name, err := getName(c)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		err = sw.enableNodes(name)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		return c.String(http.StatusOK, fmt.Sprintf("changed enable status: %s\n", name))
	})

	e.POST("/api/alldisable", func(c echo.Context) error {
		sw.disableNodes()
		return c.String(http.StatusOK, "changed all status: disable\n")
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
		if strings.EqualFold(node.prop.Name, name) {
			go changeStateNode(node, true)
		} else {
			go changeStateNode(node, false)
		}
	}
	return nil
}

func (cs *contentSwitch) disableNodes() {
	for _, node := range cs.nodes {
		go changeStateNode(node, false)
	}
}

func (cs *contentSwitch) getStatus(name string) (bool, error) {
	for _, node := range cs.nodes {
		if strings.EqualFold(node.prop.Name, name) {
			stat, err := node.client.GetStatus(nil)
			if err != nil {
				log.Warn(fmt.Sprintf("node is dead.: %s\n%s", node.prop.Name, err))
				return false, err
			}
			return stat.Enable, nil
		}
	}
	return false, fmt.Errorf("node not founded.: %s\n", name)
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

	e := createEcho(&cs)
	startEcho(e, conf.Port)

	return nil
}
