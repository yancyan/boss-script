package main

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestPath(t *testing.T) {
	dir, err := ioutil.ReadDir("file")
	if err != nil {
		panic(err)
	}
	for _, info := range dir {
		fmt.Println(info.Name())
		fmt.Println(info)
	}

}
