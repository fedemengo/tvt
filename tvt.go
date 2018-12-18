package main

import (
	"fmt"
	"os"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	app = kingpin.New("tvt", "tvt - Reserve ticket for www.tvtickets.com")

	ls = app.Command("ls", "List all available tv shows")
	rs = app.Command("rs", "Reserve tickets for shows")

	force   = app.Flag("force", "Force reserving/creating a ticket").Short('f').Bool()
	verbose = app.Flag("verbose", "Verbose output of what is happening").Short('v').Bool()

	show   = rs.Flag("show", "TV show name").String()
	first  = rs.Flag("first", "First name").String()
	last   = rs.Flag("last", "Last name").String()
	number = rs.Flag("number", "Number of tickets to reserve").Short('n').String()
	phone  = rs.Flag("phone", "Phone number").Short('p').String()
	email  = rs.Flag("email", "Email address").Short('e').String()
	config = rs.Flag("config", "Configuration file").Short('c').String()
)

func main() {

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case ls.FullCommand():

		shows := GetAvailableShows(*force)
		fmt.Println("Available shows")
		for _, s := range clear(shows) {
			fmt.Println(" - \"" + s + "\"")
		}
	case rs.FullCommand():

		data := TicketData{*number, *first, *last, *phone, *email}
		if !cliIsValid(*show, data) && len(*config) == 0 {
			fmt.Println("Couldn't get required information from command line argument or file")
		} else if ok, n, show := ReserveTicket(*show, data, *force, *verbose, *config); ok {
			fmt.Println("Successfully reserved", n, "ticket/s for", show)
		} else {
			fmt.Println("Couldn't reserve ticket/s")
		}
	}

	return
}

func clear(shows []ShowInfo) []string {

	res := make([]string, 0)
	uniq := make(map[string]bool)
	for _, s := range shows {
		if ok := uniq[s.ShowName]; !ok {
			res = append(res, s.ShowName)
			uniq[s.ShowName] = true
		}
	}

	return res
}
