package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	FetchKLines("ETHUSDT", KLineInterval1d, 4)
}

func FetchKLines(pairSymbol string, interval KLineInterval, limit uint, options ...KLineRequestOption) {
	url := KLinesURL(NewKLineReq(pairSymbol, interval, limit))
	// fmt.Println(url)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("content-type", "application/json")
	req.Header.Add("cache-control", "no-cache")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body))

	kls, err := UnmarshalKLinesJSON(body)
	if err != nil {
		panic(err)
	}
	for _, kl := range kls {
		fmt.Println(kl)
	}
}
