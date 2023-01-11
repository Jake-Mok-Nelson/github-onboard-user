package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

// Convert the http response body into a string for debugging
func readResponseBody(resp *http.Response) string {
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return string(bodyBytes)
}
