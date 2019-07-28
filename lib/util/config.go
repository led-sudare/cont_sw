package util

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

func ReadConfig(configs interface{}) error {
	buf, err := ioutil.ReadFile("config.yml")
	if err == nil {
		err := yaml.Unmarshal(buf, configs)
		if err != nil {
			return err
		}
	}
	return nil
}
