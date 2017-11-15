package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	spinner "github.com/janeczku/go-spinner"
	"github.com/urfave/cli"
)

func scrape() {
	doc, err := goquery.NewDocument("https://www.governmentjobs.com/jobs?page=1&keyword=cyber+security+&location=")
	if err != nil {
		log.Fatal(err)
	}
	pos := make([]string, 0)  // position title
	app := make([]string, 0)  // application link
	loc := make([]string, 0)  // location
	emp := make([]string, 0)  // employment type and pay i.e. full time / part time
	desc := make([]string, 0) // job description
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

func query() {
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

	app := cli.NewApp()
	app.Name = "autohunt"
	app.Usage = "Automatically search for jobs/internships in cyber security field"
	app.Action = func(c *cli.Context) error {
		fmt.Printf("searching for... %q\n", c.Args().Get(0))
		s := spinner.StartNew("This may take some while...")
		time.Sleep(5 * time.Second) // something more productive here
		s.Stop()
		return nil
	}

	app.Run(os.Args)
	//scrape()
	//query s

}
