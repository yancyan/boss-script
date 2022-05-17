package main

import (
	"fmt"
	"stariboss-script/env"
	"testing"
)

func TestConfig_InitConfig(t *testing.T) {
	c := env.Config
	fmt.Println(c)
}
