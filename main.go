package main

import (
	"fmt"
	"os"
)

// INFO represent the info I'll use to reserve the tickets
var INFO = TicketData{1, "F", "M", 1111111111, "fm@fm.io"}

func main() {

	if len(os.Args) != 0 {
		os.Args = append(os.Args, "")
	}

	switch os.Args[1] {
	case "ls":
		shows := GetAvailableShows()
		fmt.Println("Available shows")
		for _, s := range clear(shows) {
			fmt.Println("\t\"" + s + "\"")
		}
	case "rs":
		showName := getShowName()
		if ReserveTicket(showName, INFO, forced()) {
			fmt.Println("Successfully reserved", INFO.Number, "ticket/s for", showName)
		} else {
			fmt.Println("Couldn't reserve ticket/s")
		}
	default:
		help()
	}

	return

}

func getShowName() string {
	name := os.Getenv("TVT_SHOW")
	if len(name) == 0 {
		panic("Show name is not set")
	}

	return name
}

func forced() bool {
	forced := os.Getenv("FORCE")
	if forced == "true" {
		return true
	}
	return false
}

func help() {
	fmt.Println("tvt - Reserve ticket for www.tvtickets.com")
	fmt.Println("Usage:\t tvt option")
	fmt.Println()
	fmt.Println("Options:\tls - List all available tv shows")
	fmt.Println("        \trs - Reserve ticket")
	fmt.Println()
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
