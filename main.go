package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	spinner "github.com/janeczku/go-spinner"
	"github.com/urfave/cli"
)

func writeCSVHeader(out *csv.Writer) error {
	line := []string{"Position", "Application Link", "Location", "Employment Type", "Description"}
	err := out.Write(line)
	return err
}

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}

type jobs struct {
	pos  []string `csv:"Position"`
	link []string `csv:"Application Link"`
	loc  []string `csv:"Location"`
	emp  []string `csv:"Employment Type"`
	desc []string `csv:"Description"`
}

func (j *jobs) MarshalCSV(id int) ([]string, error) {
	line := make([]string, 0)
	line = append(line, j.pos[id])
	line = append(line, j.link[id])
	line = append(line, j.loc[id])
	line = append(line, j.emp[id])
	line = append(line, j.desc[id])
	//line := fmt.Sprintf("%s , %s , %s , %s , %s", j.pos[id], j.link[id], j.loc[id], j.emp[id], j.desc[id])
	return line, nil
}

func scrape(list jobs, out *csv.Writer) {
	doc, err := goquery.NewDocument("https://www.governmentjobs.com/jobs?page=1&keyword=cyber+security+&location=")
	if err != nil {
		log.Fatal(err)
	}
	/*pos := make([]string, 0)  // position title
	app := make([]string, 0)  // application link
	loc := make([]string, 0)  // location
	emp := make([]string, 0)  // employment type and pay i.e. full time / part time
	desc := make([]string, 0) // job description*/

	doc.Find(".job-item").Each(func(i int, s *goquery.Selection) {
		title := s.Find("a").Text()
		title = strings.Replace(title, ",", " ", -1)
		list.pos = append(list.pos, title)

	})

	url := "https://www.governmentjobs.com"
	doc.Find(".job-details-link").Each(func(i int, s *goquery.Selection) {
		title, _ := s.Attr("href")
		list.link = append(list.link, url+title)

	})

	doc.Find(".job-item .primaryInfo").Each(func(i int, s *goquery.Selection) {
		if i%2 == 0 {
			title := s.Text()
			title = strings.TrimSpace(title)
			title = strings.Replace(title, ",", " ", -1)
			list.loc = append(list.loc, title)

		} else {
			title := s.Text()
			title = strings.TrimSpace(title)
			title = strings.Replace(title, ",", ".", -1)
			list.emp = append(list.emp, title)

		}

	})
	doc.Find(".job-item .description").Each(func(i int, s *goquery.Selection) {
		title := s.Text()
		title = strings.Replace(title, ",", " ", -1)
		list.desc = append(list.desc, title)
	})

	for i := range list.pos {
		csvContent, err := list.MarshalCSV(i) // Get all clients as CSV string
		checkError("MarshalCSV error", err)
		err = out.Write(csvContent)
		checkError("Write to File error", err)
		//fmt.Println(csvContent) // Display all clients as CSV string
		//fmt.Printf("%s , %s , %s , %s , %s \n", list.pos[i], list.link[i], list.loc[i], list.emp[i], list.desc[i])
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

	nonFedJobs := jobs{} // hold jobs found for non federal jobs

	file, err := os.Create("autohunt_results.csv")
	checkError("Cannot write to file", err)
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	err = writeCSVHeader(writer)
	checkError("Header Write Error", err)
	app := cli.NewApp()
	app.Name = "autohunt"
	app.Usage = "Automatically search for jobs/internships in cyber security field"
	app.Action = func(c *cli.Context) error {
		fmt.Printf("searching for... %q\n", c.Args().Get(0))
		s := spinner.StartNew("This wont take long...")
		fmt.Println()
		scrape(nonFedJobs, writer)

		//time.Sleep(5 * time.Second) // something more productive here
		s.Stop()
		return nil
	}

	app.Run(os.Args)

	//query s

}
