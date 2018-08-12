package dyndns

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Urls contains a set of mirrors in which a
// raw IP string can be retreived. It is exported
// for the intent of modification.
var (
	Urls = []string{"http://myexternalip.com/raw", "https://api.ipify.org"}
)

func tryMirror(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(contents), nil
}

// GetExternalIP retrieves the external facing IP Address.
// If multiple mirrors are provided in Urls,
// it will try each one (in order), should
// preceding mirrors fail.
func GetExternalIP() (string, error) {
	for _, url := range Urls {
		resp, err := tryMirror(url)
		if err == nil {
			return resp, nil
		} else {
			fmt.Printf("Error in request: %s", err.Error())
		}
	}

	return "", errors.New("Could not retreive external IP")
}
