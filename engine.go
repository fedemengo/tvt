package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const delay = 3

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
	soldOut    bool
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

// reserveURL is the endpoint for reserving a ticket
const reserveURL = "http://www.tvtickets.com/fmi/tickets/confirmation.php"

// GetAvailableShows return a slice of available TV shows
func GetAvailableShows() []ShowInfo {
	res, err := http.Get(listURL)
	if err != nil {
		panic("The HTTP request failed: cannot retrieve shows list")
	}

	//cookies := res.Cookies()
	//fmt.Println("Cookies")
	//for _, cookie := range cookies {
	//	fmt.Println("\t", cookie.Name, "=", cookie.Value)
	//}
	//fmt.Println()

	//header := res.Header
	//fmt.Println("Header")
	//for k := range header {
	//	for v := range header[k] {
	//		fmt.Println("\t", k, ":", v)
	//	}
	//}
	//fmt.Println()

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

	// TODO: Complete here
	for _, s := range shows {
		s.soldOut = false
	}

	return shows
}

// ReserveTicket attempts to reserve a ticket, and return true in case of success
func ReserveTicket(showName string, data TicketData, forced bool) (succeed bool) {
	shows := GetAvailableShows()

	succeed = false
	for id, s := range shows {
		// If the current show is the correct show AND the show is not soldout OR forced is set to true, attempt to reserve ticket
		if strings.Index(s.ShowName, showName) != -1 && (s.soldOut || forced) {
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
				ioutil.WriteFile("./"+ticketName+".html", body, 0644)
			}
		}
	}

	return
}

func ticketIsValid(body []byte) bool {
	return true
}
