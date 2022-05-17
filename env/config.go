package main

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type config struct {
	Domain string `json:"domain"`
}

func (c *config) InitConfig() *config {
	f, err := ioutil.ReadFile("config/application.yaml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(f, c)
	if err != nil {
		panic(err)
	}
	return c
}
