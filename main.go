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
	pos  []string //`csv:"Position"`
	link []string //`csv:"Application Link"`
	loc  []string //`csv:"Location"`
	emp  []string //`csv:"Employment Type"`
	desc []string //`csv:"Description"`
}

func (j *jobs) MarshalCSV(id int) ([]string, error) {
	line := make([]string, 0)
	line = append(line, j.pos[id])
	line = append(line, j.link[id])
	line = append(line, j.loc[id])
	line = append(line, j.emp[id])
	line = append(line, j.desc[id])
	return line, nil
}

type jobDetails struct {
	Name string
}

type location struct {
	LocationName           string
	CountryCode            string
	CountrySubDivisionCode string
	CityName               string
}

type jobObject struct {
	PositionTitle        string
	PositionURI          string
	ApplyURI             []string
	QualificationSummary string
	OrganizationName     string
	DepartmentName       string
	PositionLocation     []location
	PositionOfferingType []jobDetails
}

type results struct {
	MatchedObjectDescriptor jobObject
}

type searchObject struct {
	SearchResultItems []results
}

type jsonObject struct {
	SearchResult searchObject
}

func (j *jsonObject) MarshalCSV(id int) ([]string, error) {
	line := make([]string, 0)
	line = append(line, j.SearchResult.SearchResultItems[id].MatchedObjectDescriptor.PositionTitle)
	line = append(line, j.SearchResult.SearchResultItems[id].MatchedObjectDescriptor.ApplyURI[0])
	line = append(line, j.SearchResult.SearchResultItems[id].MatchedObjectDescriptor.PositionLocation[0].CityName)
	line = append(line, j.SearchResult.SearchResultItems[id].MatchedObjectDescriptor.PositionOfferingType[0].Name)
	line = append(line, j.SearchResult.SearchResultItems[id].MatchedObjectDescriptor.QualificationSummary)
	return line, nil
}
func scrape(list jobs, out *csv.Writer) {
	doc, err := goquery.NewDocument("https://www.governmentjobs.com/jobs?page=1&keyword=cyber+security+&location=")
	if err != nil {
		log.Fatal(err)
	}

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

func query(out *csv.Writer) {

	req, err := http.NewRequest("GET", "https://data.usajobs.gov/api/search?JobCategoryCode=2210&Keyword=", nil)
	checkError("request error", err)
	req.Header.Set("Host", "data.usajobs.gov")
	//req.Header.Set("User-Agent", os.ExpandEnv("axp141330@utdallas.edu"))
	//req.Header.Set("Authorization-Key", os.ExpandEnv("YbqpIDaiMNCy7BnTbVVaSevXNgzsfl7S1KwA9beNtro="))

	req.Header.Set("User-Agent", "axp141330@utdallas.edu")
	req.Header.Set("Authorization-Key", "YbqpIDaiMNCy7BnTbVVaSevXNgzsfl7S1KwA9beNtro=") // I know Probably should hide this but ehh

	resp, err := http.DefaultClient.Do(req)
	checkError("response error", err)
	body, err := ioutil.ReadAll(resp.Body)
	checkError("read error", err)
	//fmt.Print(string(body))
	var j jsonObject
	err = json.Unmarshal(body, &j)
	checkError("json error", err)
	for i := range j.SearchResult.SearchResultItems {
		csvContent, err := j.MarshalCSV(i) // Get all clients as CSV string
		checkError("MarshalCSV error", err)
		err = out.Write(csvContent)
		checkError("Write to File error", err)
		/*b, err := json.Marshal(j)
		checkError("json error", err)
		fmt.Print(string(b))*/
	}
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
		s := spinner.StartNew("This wont take long ...")
		//time.Sleep(1 * time.Second)
		//fmt.Println("")

		//time.Sleep(5 * time.Second) // something more productive here
		scrape(nonFedJobs, writer)
		query(writer)

		s.Stop()
		return nil
	}

	app.Run(os.Args)

}
