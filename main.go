package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func isEmpty(str string) bool {
	if len(str) == 0 {
		return true
	}
	return false
}

func scrape() {
	doc, err := goquery.NewDocument("https://www.governmentjobs.com/jobs?keyword=cyber+security+&location=")
	if err != nil {
		log.Fatal(err)
	}
	pos := make([]string, 0)
	app := make([]string, 0)
	loc := make([]string, 0)
	emp := make([]string, 0)
	desc := make([]string, 0)
	doc.Find(".job-item").Each(func(i int, s *goquery.Selection) {
		title := s.Find("a").Text()
		title = strings.Replace(title, ",", " ", -1)
		pos = append(pos, title)

	})

	url := "https://www.governmentjobs.com"
	doc.Find(".job-details-link").Each(func(i int, s *goquery.Selection) {
		title, _ := s.Attr("href")
		app = append(app, url+title)

	})

	doc.Find(".job-item .primaryInfo").Each(func(i int, s *goquery.Selection) {
		if i%2 == 0 {
			title := s.Text()
			title = strings.TrimSpace(title)
			title = strings.Replace(title, ",", " ", -1)
			loc = append(loc, title)

		} else {
			title := s.Text()
			title = strings.TrimSpace(title)
			title = strings.Replace(title, ",", ".", -1)
			emp = append(emp, title)

		}

	})
	doc.Find(".job-item .description").Each(func(i int, s *goquery.Selection) {
		title := s.Text()
		title = strings.Replace(title, ",", " ", -1)
		desc = append(desc, title)
	})

	for i := range pos {
		fmt.Printf("%s , %s , %s , %s , %s \n", pos[i], app[i], loc[i], emp[i], desc[i])
	}

}

func ExampleQuery() {
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
func main() {

	scrape()
	//getLinks()

}
