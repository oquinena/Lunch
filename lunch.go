package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	c "github.com/TreyBastian/colourize"
	"github.com/mmcdole/gofeed"
)

type lunchalt struct {
	vanlig string
	veg    string
}

func main() {
	all := flag.Bool("a", false, "Visar även vegetariska rätter")
	b := flag.Bool("b", false, "Endast Birkaskolans lunch")
	e := flag.Bool("e", false, "Visa endast Ekebyhovskolans lunch")
	flag.Parse()

	var dagar = map[string]time.Weekday{
		"Måndag":  time.Monday,
		"Tisdag":  time.Tuesday,
		"Onsdag":  time.Wednesday,
		"Torsdag": time.Thursday,
		"Fredag":  time.Friday,
		"Lördag":  time.Saturday,
		"Söndag":  time.Sunday,
	}

	dag := time.Now().Weekday()
	if dag == dagar["Lördag"] || dag == dagar["Söndag"] {
		fmt.Println("Helg")
		os.Exit(0)
	}

	birka, _ := lunch("birkaskolan")
	ekebyhov, _ := lunch("ekebyhovskolan")

	if *all {
		fmt.Printf(c.Colourize("Birkaskolan:", c.Green))
		fmt.Printf("\n%s\n%s\n\n", birka.vanlig, birka.veg)
		fmt.Printf(c.Colourize("Ekebyhovskolan:", c.Green))
		fmt.Printf("\n%s\n%s\n", ekebyhov.vanlig, ekebyhov.veg)
	} else if *b {
		fmt.Println(birka.vanlig)
	} else if *e {
		fmt.Println(ekebyhov.vanlig)
	} else {
		fmt.Println(c.Colourize("Birkaskolan:    ", c.Green), birka.vanlig)
		fmt.Println(c.Colourize("Ekebyhovskolan: ", c.Green), ekebyhov.vanlig)
	}
}

func lunch(skola string) (lunchalt, error) {
	var lv lunchalt
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL("https://skolmaten.se/" + skola + "/rss/days/")
	if err != nil {
		fmt.Errorf("Fel vid hämtning av lunchalternativ")
		return lv, err
	}
	l := feed.Items[0].Description
	lunch := strings.Split(l, "<br/>")
	lv.vanlig, lv.veg = lunch[0], lunch[1]

	return lv, err
}
