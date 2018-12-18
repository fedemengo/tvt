package main

import (
	"fmt"
	"os"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	app = kingpin.New("tvt", "tvt - Reserve ticket for www.tvtickets.com")

	ls = app.Command("ls", "List all available tv shows")
	rs = app.Command("rs", "Reserve ticket")

	force   = app.Flag("force", "Force reserving/creating a ticket").Short('f').Bool()
	verbose = app.Flag("verbose", "Verbose output of what is happening").Short('v').Bool()

	showName = rs.Arg("show", "TV show name").String()
	first    = rs.Arg("first", "First name").String()
	last     = rs.Arg("last", "Last name").String()
	number   = rs.Arg("number", "Number of tickets to reserve").Int()
	phone    = rs.Arg("phone", "Phone number").Int()
	email    = rs.Arg("email", "Email address").String()
	config   = rs.Arg("config", "Configuration file").String()
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
		if ReserveTicket(*showName, TicketData{*number, *first, *last, *phone, *email}, *force, *verbose) {
			fmt.Println("Successfully reserved", *number, "ticket/s for", *showName)
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
