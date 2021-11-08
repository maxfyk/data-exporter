package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type HttpHeader struct {
	Name  string
	Value string
}

type Configuration struct {
	BaseUrl                   string
	RegexpPattern             string
	RegexpPatternPages        string
	CsvHeaders                []string
	HttpHeaders               []HttpHeader
	RegexpIndexes             []uint8
	RegexpIndexVote           uint8
	RegexpConvertVoteToNumber bool
	FileName                  string
}

var (
	config Configuration
)

func main() {

	start := time.Now()

	baseUrl := flag.String("baseUrl", "", "base url")
	fileName := flag.String("fileName", "", "fileName.csv")
	configFileName := flag.String("config", "config.json", "config.json")

	flag.Parse()

	readConfig(*configFileName)

	if len(*baseUrl) != 0 {
		config.BaseUrl = *baseUrl
	}
	if len(*fileName) != 0 {
		config.FileName = *fileName
	}

	if len(config.BaseUrl) == 0 {
		log.Fatal("Error: baseUrl is required")
	}
	if config.BaseUrl == "https://www.kinopoisk.ru/user/XXX/votes/list/ord/date/page/%d/" {
		log.Fatal("Error: change XXX in baseUrl to your profile id")
	}
	if len(config.FileName) == 0 {
		log.Fatal("Error: fileName is required")
	}

	parseAllItmes()
	elapsed := time.Since(start)
	fmt.Println("Time took ", elapsed)

}

func readConfig(configFileName string) {
	file, err := os.Open(configFileName)
	defer file.Close()

	if err != nil {
		ex, err := os.Executable()
		if err != nil {
			log.Fatal(err)
		}
		exPath := filepath.Dir(ex)
		file, err = os.Open(filepath.Join(exPath, configFileName))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(filepath.Join(exPath, configFileName))
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	checkError("Error parsing config:", err)
}

func getMaxPage(html string) uint {
	re := regexp.MustCompile(config.RegexpPatternPages)
	a := re.FindAllStringSubmatch(html, -1)

	pageNumbers := make([]uint, 0)
	for _, value := range a {
		u64, err := strconv.ParseUint(value[1], 10, 32)
		if err != nil {
			fmt.Println(err)
		}
		pageNumbers = append(pageNumbers, uint(u64))
	}
	return MaxIntInSlice(pageNumbers)
}

func parseAllItmes() {

	re := regexp.MustCompile(config.RegexpPattern)

	//getFirstPage
	fmt.Print("getting page 1...")

	html, err := getHtmlByPage(1)
	checkError("Cannot get html", err)

	html = strings.ReplaceAll(html, "\n", "")
	a := re.FindAllStringSubmatch(html, -1)
	maxPage := getMaxPage(html)

	if config.RegexpConvertVoteToNumber {
		for j := 0; j < len(a); j++ {
			vote, _ := strconv.Atoi(a[j][config.RegexpIndexVote])
			a[j][config.RegexpIndexVote] = strconv.Itoa(vote)
		}
	}

	fmt.Println("finished")

	//write to csv boilerplate
	file, err := os.Create(config.FileName)
	checkError("Cannot create file", err)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.Write(config.CsvHeaders)
	checkError("Cannot write to file", err)
	//end write to csv boilerplate

	//write header to csv
	for _, value := range a {
		err := writer.Write(getElementsByIndexes(value, config.RegexpIndexes))
		checkError("Cannot write to file", err)
	}
	//end write header to csv

	//get All pages

	for i := uint(2); i <= maxPage; i++ {
		fmt.Print("getting page ", i, "...")
		html, err := getHtmlByPage(i)
		checkError("Cannot get html", err)

		html = strings.ReplaceAll(html, "\n", "")
		a = re.FindAllStringSubmatch(html, -1)

		if config.RegexpConvertVoteToNumber {
			for j := 0; j < len(a); j++ {
				vote, _ := strconv.Atoi(a[j][config.RegexpIndexVote])
				a[j][config.RegexpIndexVote] = strconv.Itoa(vote)
			}
		}

		//write to csv
		for _, value := range a {
			err := writer.Write(getElementsByIndexes(value, config.RegexpIndexes))
			checkError("Cannot write to file", err)
		}
		fmt.Println("finished")
		//end write to csv
	}

}

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}

func MaxIntInSlice(array []uint) uint {
	var max uint = array[0]
	for _, value := range array {
		if max < value {
			max = value
		}
	}
	return max
}

func getElementsByIndexes(data []string, indexes []uint8) []string {
	res := make([]string, 0)

	for _, v := range indexes {
		res = append(res, data[v])
	}
	return res
}

func getHtmlByPageLocal(page uint) (string, error) {
	content, err := ioutil.ReadFile("kp_with_seen.html")
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func getHtmlByPage(page uint) (string, error) {

	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, fmt.Sprintf(config.BaseUrl, page), nil)

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	req.Header.Add("cache-control", "max-age=0")
	req.Header.Add("upgrade-insecure-requests", "1")
	req.Header.Add("sec-gpc", "1")
	req.Header.Add("sec-fetch-site", "none")
	req.Header.Add("sec-fetch-mode", "navigate")
	req.Header.Add("sec-fetch-user", "?1")
	req.Header.Add("sec-fetch-dest", "document")

	for _, v := range config.HttpHeaders {
		req.Header.Add(v.Name, v.Value)
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return string(body), nil
}
