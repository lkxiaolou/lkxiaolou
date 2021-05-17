package toolx

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	DefaultTimeout = 8
)

var contentCache = make(map[string]string)

func HttpGetWithCache(destUrl string) (string, error) {
	if content, ok := contentCache[destUrl]; ok {
		return content, nil
	}
	content, err := HttpGet(destUrl)
	if err == nil {
		contentCache[destUrl] = content
	}
	return content, err
}

func HttpGet(destUrl string) (string, error) {
	var err error
	var resp *http.Response
	var body []byte

	client := &http.Client{
		Timeout: DefaultTimeout * time.Second,
	}

	resp, err = client.Get(destUrl)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = fmt.Errorf("The http resp code is not 200! Code:%d", resp.StatusCode)
		return "", err
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
