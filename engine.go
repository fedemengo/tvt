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

// TicketData contains the necessary data to reserve a ticket
type TicketData struct {
	Number int
	First  string
	Last   string
	Phone  int
	Email  string
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

// ReserveTicket attempts to reserve a ticket, and return true in case of success
func ReserveTicket(showName string, data TicketData, forced, verbose bool) (succeed bool) {

	succeed = false
	for succeed == false {
		shows := GetAvailableShows(forced)
		for id, s := range shows {
			// If the current show is the correct show AND the show is not soldout OR forced is set to true, attempt to reserve ticket
			if strings.Index(s.ShowName, showName) != -1 {
				// Create form values for requesting the ticket
				formValues := url.Values{
					"ShowDateTime": {s.ID},
					"Number":       {strconv.Itoa(data.Number)},
					"First":        {data.First},
					"Last":         {data.Last},
					"Phone":        {strconv.Itoa(data.Phone)},
					"Email":        {data.Email},
					"FindUs":       {"..."},
					"-ne:":         {"Submit"},
				}

				// Send request with form data
				res, err := http.PostForm(reserveURL, formValues)
				if err != nil {
					panic("The HTTP request failed: cannot request ticket")
				}

				// If successfully reserved ticket, store it as plain html file
				body, _ := ioutil.ReadAll(res.Body)
				defer res.Body.Close()
				if ticketIsValid(body) {
					succeed = true
					ticketName := showName + "-" + data.First + data.Last + strconv.Itoa(id)
					ioutil.WriteFile("./t"+ticketName+".html", body, 0644)
				}
			}
		}
		if succeed == false && verbose {
			fmt.Println("Couldn't reserve ticket. Retring in", timeOut, "second.")
		}
		time.Sleep(timeOut * time.Second)
	}

	return
}

func ticketIsValid(body []byte) bool {
	html := string(body)
	text := strings.TrimSpace(html[strings.Index(html, "<title>")+len("<title>") : strings.Index(html, "</title>")])
	return text == "Confirmation"
}
