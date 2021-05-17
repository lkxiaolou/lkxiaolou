package toolx

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	DefaultTimeout = 8
)

func HttpGet(destUrl string, args map[string]string) (string, error) {
	var err error
	var resp *http.Response
	var body []byte

	url := destUrl

	if len(args) > 0 {
		first := true
		if strings.Contains(url, "?") {
			first = false
		}
		for k, v := range args {
			if first {
				url = fmt.Sprintf("%s?%s=%s", url, k, v)
				first = false
			} else {
				url = fmt.Sprintf("%s&%s=%s", url, k, v)
			}
		}
	}

	client := &http.Client{
		Timeout: DefaultTimeout * time.Second,
	}

	resp, err = client.Get(url)
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
