package jtltojson

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

//Result struct: a data structure to hold the cleaned data of a raw JTL file
type Result struct {
	TestName            string           `json:"testname"`
	DateTime            string           `json:"datetime"`
	NintiethPercentiles map[string]int   `json:"ninetiethpercentiles"`
	ResponseTimes       map[string][]int `json:"responsetimes"`
}

//Ptor = Protocol-to-Result | similar shorthand to Itoa or Atoi
//Protocol level user (JMeter-http/api) responses to Result struct
func Ptor(file string) *Result {
	data := readJTL(file)
	rt := mapResponseTimes(readJTL(file))
	np := findNintiethPercentiles(rt)
	return &Result{
		TestName:            file,
		DateTime:            formatDateTime(data[0][0]),
		NintiethPercentiles: np,
		ResponseTimes:       rt,
	}
}

//Btor = Browser-to-Result | similar shorthand to Itoa or Atoi
//Browser level user (JMeter-Selenium WebDriver) responses to Result struct
func Btor(file string) *Result {
	data := readJTL(file)
	rt := mapResponseTimes(data)
	np := findNintiethPercentiles(rt)
	return &Result{
		TestName:            file,
		DateTime:            formatDateTime(data[0][0]),
		NintiethPercentiles: np,
		ResponseTimes:       rt,
	}
}

//Calculates the 90th percentile response time for each label's list of response times
func findNintiethPercentiles(rt map[string][]int) map[string]int {
	var np map[string]int
	for label, slice := range rt {
		sort.Ints(slice)
		np[label] = slice[int(float64(len(slice))*.9)]
	}
	return np
}

//Traverses data to create a map of labels with a corresponding list of responses times
func mapResponseTimes(data [][]string) map[string][]int {
	var rt map[string][]int
	for i := 0; i < len(data); i++ {
		label := data[i][2]
		elapsed, err := strconv.Atoi(data[i][1])
		checkErr(err)
		if (strings.EqualFold(label, "Setup Sampler") == false) &&
			(strings.EqualFold(label, "Setup Request") == false) &&
			(strings.EqualFold(label, "label") == false) {
			rt[label] = append(rt[label], elapsed)
		}
	}
	return rt
}

//Formats an epoch date/time string into the format defined within the function
func formatDateTime(dt string) string {
	epoch, err := strconv.ParseInt(dt, 10, 64)
	checkErr(err)
	dt = time.Unix(epoch, 0).Format("Jan-02-06 3:04pm")
	return dt
}

//Reads a JTL result file (formatted as CSV) into a list of records each containing a slice of fields from the CSV
func readJTL(file string) [][]string {
	csvFile, err := os.Open(file)
	checkErr(err)
	data, err := csv.NewReader(csvFile).ReadAll()
	checkErr(err)
	return data
}

//JSON encodes a Result struct
func (r *Result) JSON() []byte {
	json, err := json.Marshal(r)
	checkErr(err)
	return json
}

//Checks an error and then logs and prints accordingly
func checkErr(err error) {
	if err != nil {
		pc, file, line, _ := runtime.Caller(1)
		function := strings.TrimPrefix(filepath.Ext(runtime.FuncForPC(pc).Name()), ".")
		fmt.Println("[" + time.Now().Format("Jan-02-06 3:04pm") + "] Error Warning:" + file + " " + function + "() line:" + strconv.Itoa(line) + " " + err.Error())
		log.Fatal(err)
	}
}
