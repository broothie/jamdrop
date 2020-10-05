// +build prod_build_info
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"jamdrop/config"
	"net/http"
	"time"
)

func main() {
	res, err := http.Get("https://jamdrop.app/build_info")
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var info config.BuildInfo
	if err := json.Unmarshal(body, &info); err != nil {
		panic(err)
	}

	location, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		panic(err)
	}

	t, err := time.ParseInLocation("2006-01-02T15:04:05", info.BuildTime, location)
	if err != nil {
		panic(err)
	}

	fmt.Println("last build time:", t)
}
