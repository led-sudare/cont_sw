package util

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

func ReadConfig(configs interface{}) error {
	val, err := ioutil.ReadFile("./config.yml")

	if err == nil {
		err := yaml.Unmarshal(val, &configs)
		if err != nil {
			return err
		}
	}
	return nil
}
