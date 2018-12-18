package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const timeOut = 5

// ShowInfo contains the metada of each show
type ShowInfo struct {
	SortDate   string `json:"sort_date"`
	Time       string `json:"ShowTime"`
	MinimumAge string `json:"Age"`
	Special    string `json:"Special"`
	Date       string `json:"ShowDate"`
	Day        string `json:"ShowDay"`
	ShowName   string `json:"ShowName"`
	ID         string `json:"rID"`
}

type showFileInfo struct {
	ShowName string       `json:"Show"`
	Data     []TicketData `json:"Data"`
}

// TicketData contains the necessary data to reserve a ticket
type TicketData struct {
	Number string `json:"Number"`
	First  string `json:"First"`
	Last   string `json:"Last"`
	Phone  string `json:"Phone"`
	Email  string `json:"Email"`
}

// listURL is the endpoint for the list of shows
const listURL = "http://www.tvtickets.com/api/public/index.php/show-schedule"

// listAvailableURL is the endpoint for the list of available shows
const listAvailableURL = "https://www.tvtickets.com/api/public/index.php/available-shows"

// reserveURL is the endpoint for reserving a ticket
const reserveURL = "http://www.tvtickets.com/fmi/tickets/confirmation.php"

// GetAvailableShows return a slice of available TV shows
func GetAvailableShows(forced bool) []ShowInfo {
	url := listAvailableURL
	if forced {
		url = listURL
	}

	res, err := http.Get(url)
	if err != nil {
		panic("The HTTP request failed: cannot retrieve shows list")
	}

	var shows []ShowInfo
	data, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		panic("Read error: cannot read response body")
	}

	err = json.Unmarshal(data, &shows)
	if err != nil {
		panic("Parsing error: cannot create JSON from data")
	}

	return shows
}

func readFile(configFile string) (bool, map[string][]TicketData) {
	if len(configFile) == 0 {
		return false, nil
	}

	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic("Couldn't open file")
	}

	var showData []showFileInfo
	err = json.Unmarshal(data, &showData)
	if err != nil {
		panic("Parsing error: cannot create JSON from data")
	}

	result := make(map[string][]TicketData)
	for _, e := range showData {
		result[e.ShowName] = e.Data
	}

	return true, result
}

func commandLineValid(show string, data TicketData) bool {
	return len(show) != 0 && len(data.Number) != 0 && len(data.First) != 0 && len(data.Last) != 0 && len(data.Phone) != 0 && len(data.Email) != 0
}

func sendRequest(id int, showName string, values url.Values) bool {

	values.Set("FindUs", "...")
	values.Set("-ne:", "Submit")

	// Send request with form data
	res, err := http.PostForm(reserveURL, values)
	if err != nil {
		panic("The HTTP request failed: cannot request ticket")
	}

	// If successfully reserved ticket, store it as plain html file
	body, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if ticketIsValid(body) {
		ticketName := showName + "-" + values.Get("First") + values.Get("Last") + strconv.Itoa(id)
		ioutil.WriteFile("./t"+ticketName+".html", body, 0644)
		return true
	}

	return false
}

// ReserveTicket attempts to reserve a ticket, and return true in case of success
func ReserveTicket(showName string, data TicketData, forced, verbose bool, configFile string) (bool, int, []string) {

	file, fileData := readFile(configFile)
	cli := commandLineValid(showName, data)

	succeed := false
	nTickets := 0
	nShows := make([]string, 0)
	for succeed == false {
		shows := GetAvailableShows(forced)
		for id, s := range shows {
			// Attempt to reserve ticket with options from file
			if entry, ok := fileData[s.ShowName]; file && ok {
				for iid, e := range entry {
					formValues := url.Values{
						"ShowDateTime": {s.ID},
						"Number":       {e.Number},
						"First":        {e.First},
						"Last":         {e.Last},
						"Phone":        {e.Phone},
						"Email":        {e.Email},
					}

					if sendRequest(id*10+iid, s.ShowName, formValues) {
						succeed = true
						nTickets++
						nShows = append(nShows, s.ShowName)
					}
				}
			}

			// Attempt to reserve ticket with options from command line
			if cli && strings.Index(s.ShowName, showName) != -1 {
				// Create form values for requesting the ticket
				formValues := url.Values{
					"ShowDateTime": {s.ID},
					"Number":       {data.Number},
					"First":        {data.First},
					"Last":         {data.Last},
					"Phone":        {data.Phone},
					"Email":        {data.Email},
				}

				if sendRequest(id, s.ShowName, formValues) {
					succeed = true
					nTickets++
					nShows = append(nShows, s.ShowName)
				}
			}
		}

		if succeed == false && verbose {
			fmt.Println("Couldn't reserve ticket. Retring in", timeOut, "second.")
		}
		time.Sleep(timeOut * time.Second)
	}

	return succeed, nTickets, nShows
}

func ticketIsValid(body []byte) bool {
	html := string(body)
	text := strings.TrimSpace(html[strings.Index(html, "<title>")+len("<title>") : strings.Index(html, "</title>")])
	return text == "Confirmation"
}
