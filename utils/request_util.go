package main

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
)

type environment struct {
	Url    string
	Cookie string
}

func Get(env environment, params map[string]string) []byte {
	request, err := http.NewRequest("GET", env.Url, nil)
	if err != nil {
		panic(err)
	}
	request.Header.Add("Cookie", env.Cookie)
	values := request.URL.Query()
	for key, value := range params {
		values.Add(key, value)
	}
	request.URL.RawQuery = values.Encode()
	readCloser := sendRequest(request)

	defer func(readCloser io.ReadCloser) {
		err := readCloser.Close()
		if err != nil {
			panic(err)
		}
	}(readCloser)

	all, err := ioutil.ReadAll(readCloser)
	if err != nil {
		panic(err)
	}
	log.Infof("the body is %s", string(all))
	return all
}

func Post(env environment, params []byte) []byte {
	request, err := http.NewRequest("POST", env.Url, bytes.NewBuffer(params))
	if err != nil {
		panic(err)
	}

	request.Header.Add("Cookie", env.Cookie)
	request.Header.Add("X-Requested-With", "XMLHttpRequest")
	request.Header.Add("Content-Type", "application/json")

	resp := sendRequest(request)

	defer func(resp io.ReadCloser) {
		err := resp.Close()
		if err != nil {
			panic(err)
		}
	}(resp)

	all, err := ioutil.ReadAll(resp)
	if err != nil {
		panic(err)
	}
	log.Infof("the body is %s", string(all))
	return nil
}

/**
 * 使用panic，是否继续由上层处理
 */
func sendRequest(request *http.Request) io.ReadCloser {
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}
	if response.StatusCode != 200 {
		all, err := ioutil.ReadAll(response.Body)
		if err != nil {
			panic(err)
		}
		panic("status code not 200, " + string(all))
	}
	return response.Body
}
