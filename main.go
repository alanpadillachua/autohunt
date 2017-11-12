package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {

	req, err := http.NewRequest("GET", "https://data.usajobs.gov/api/search?JobCategoryCode=2210", nil)
	if err != nil {
		// handle err
	}
	req.Header.Set("Host", "data.usajobs.gov")
	//req.Header.Set("User-Agent", os.ExpandEnv("axp141330@utdallas.edu"))
	//req.Header.Set("Authorization-Key", os.ExpandEnv("YbqpIDaiMNCy7BnTbVVaSevXNgzsfl7S1KwA9beNtro="))

	req.Header.Set("User-Agent", "axp141330@utdallas.edu")
	req.Header.Set("Authorization-Key", "YbqpIDaiMNCy7BnTbVVaSevXNgzsfl7S1KwA9beNtro=") // I know Probably should hide this but ehh

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Print("responce Error")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print("read error")
	}
	var f interface{}
	err = json.Unmarshal(body, &f)
	if err != nil {
		log.Print("json error")
	}
	m := f.(map[string]interface{})
	fmt.Println(m)

	defer resp.Body.Close()
}
