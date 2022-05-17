package utils

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
)

type requestUtil struct {
}

var RequestUtil = requestUtil{}

type RequestContext struct {
	Url    string
	Cookie string
}

func (rq requestUtil) Get(env RequestContext, params map[string]string) []byte {

	defer func() {
		if err := recover(); err != nil {
			log.Errorln(err)
			panic(err)
		}
	}()

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
	log.Infof("the resp body is %s", string(all))
	return all
}

func (rq requestUtil) Post(env RequestContext, params []byte) ([]byte, error) {

	defer func() {
		if err := recover(); err != nil {
			log.Errorf("the request occur error, msg is %+v", err)
			panic(err)
		}
	}()

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
		return nil, err
	}
	//log.Infof("the resp body is %s", string(all))
	return all, nil
}

func sendRequest(request *http.Request) io.ReadCloser {
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}
	if response.StatusCode != 200 {
		defer response.Body.Close()
		all, err := ioutil.ReadAll(response.Body)
		if err != nil {
			panic(err)
		}
		panic("the status " + response.Status + ", " + string(all))
	}
	return response.Body
}
