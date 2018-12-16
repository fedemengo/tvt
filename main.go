package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type result struct {
	SortDate string `json:"sort_date"`
	ShowTime string `json:"ShowTime"`
	Age      string `json:"Age"`
	Special  string `json:"Special"`
	Date     string `json:"ShowDate"`
	Day      string `json:"ShowDay"`
	Name     string `json:"ShowName"`
	RID      string `json:"rID"`
}

type myData struct {
	ShowDateTime int64  `json:"ShowDateTime"`
	Number       int    `json:"Number"`
	First        string `json:"First"`
	Last         string `json:"Last"`
	Phone        int    `json:"Phone"`
	Email        string `json:"Email"`
	FindUs       string `json:"FindUs"`
	New          string `json:"-new"`
}

// LIST is the endpoint for the list of shows
const LIST = "http://www.tvtickets.com/api/public/index.php/show-schedule"

// RESERVE is the endpoint for reserving a ticket
const RESERVE = "http://www.tvtickets.com/fmi/tickets/confirmation.php"

// INFO represent the info I'll use to reserve the tickets
var INFO = myData{-1, 5, "F", "M", 1111111111, "fm@fm.io", "x", "Submit"}

func main() {

	response, err := http.Get(LIST)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		return
	}

	var shows []result
	cookies := response.Cookies()
	header := response.Header
	data, _ := ioutil.ReadAll(response.Body)
	err = json.Unmarshal(data, &shows)

	fmt.Println("Cookies")
	for _, cookie := range cookies {
		fmt.Println("\t", cookie.Name, "=", cookie.Value)
	}
	fmt.Println()

	fmt.Println("Header")
	for k := range header {
		for v := range header[k] {
			fmt.Println("\t", k, ":", v)
		}
	}
	fmt.Println()

	if err != nil {
		fmt.Println("Can't unmarshal results")
		return
	}

	rIDs := make([]string, 0)
	for _, show := range shows {
		if strings.Index(show.Name, "The Big Bang Theory") != -1 {
			//fmt.Println(show.Name, "has ID", show.RID)
			//val, _ := strconv.ParseInt(show.RID, 10, 64)
			rIDs = append(rIDs, show.RID)
		}
	}

	fmt.Println("The Big Bang Theory IDs:", rIDs)

	formValues := url.Values{}
	formValues.Set("ShowDateTime", rIDs[0])
	formValues.Set("Number", "1")
	formValues.Set("First", "L")
	formValues.Set("Last", "P")
	formValues.Set("Phone", "21444112")
	formValues.Set("Email", "fkk25099@iencm.com")
	formValues.Set("FindUs", "...")
	formValues.Set("-new", "Submit")

	res, err := http.PostForm(RESERVE, formValues)
	if err == nil {
		body, _ := ioutil.ReadAll(res.Body)
		defer response.Body.Close()

		ioutil.WriteFile("./body.html", body, 0644)
		//fmt.Println(string(body))
	} else {
		fmt.Println("An error occoured:", err)
	}
}
