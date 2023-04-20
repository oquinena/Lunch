// testing testing
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	c "github.com/TreyBastian/colourize"
	"github.com/getlantern/systray"
	"github.com/mmcdole/gofeed"
	"github.com/oquinena/lunch/icon"
)

type lunchalt struct {
	vanlig string
	veg    string
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

func runTray() {
   	birka, err := lunch("birkaskolan")
    if err != nil {
        fmt.Println(err)
    }
	ekebyhov, err := lunch("ekebyhovskolan")
    if err != nil {
        fmt.Println(err)
    }

    systray.SetIcon(icon.Data)
    systray.SetTitle("Lunch")
    systray.SetTooltip("Lunch")
    systray.AddMenuItem("Birkaskolan:", "Birka")
    mBirkaVanlig := systray.AddMenuItem(fmt.Sprintf("Vanlig: %s", birka.vanlig), "Birka")
    mBirkaVeg := systray.AddMenuItem(fmt.Sprintf("Vegetarisk: %s", birka.veg), "Birka")
    systray.AddMenuItem("Ekebyhovskolan:", "Ekebyhov")
    mEkebyhovVanlig := systray.AddMenuItem(fmt.Sprintf("Vanlig: %s", ekebyhov.vanlig), "Ekebyhov")
    mEkebyhovVeg := systray.AddMenuItem(fmt.Sprintf("Vegetarisk: %s", ekebyhov.veg), "Ekebyhov")
    mQuit := systray.AddMenuItem("Quit", "Quit the app")
    go func() {
        for {
            select {
            case <-mBirkaVanlig.ClickedCh:
                fmt.Println("Birka vanlig")
            case <-mBirkaVeg.ClickedCh:
                fmt.Println("Birka veg")
            case <-mEkebyhovVanlig.ClickedCh:
                fmt.Println("Ekebyhov vanlig")
            case <-mEkebyhovVeg.ClickedCh:
                fmt.Println("Ekebyhov veg")
            case <-mQuit.ClickedCh:
                systray.Quit()
            }
        }
    }()

    ticker := time.NewTicker(time.Minute * 5)
    tQuit := make(chan bool)
    go func() {
        for {
            select {
            case <-ticker.C:
                birka, err := lunch("birkaskolan")
                if err != nil {
                    fmt.Println(err)
                }
                ekebyhov, err := lunch("ekebyhovskolan")
                if err != nil {
                    fmt.Println(err)
                }
                mBirkaVanlig.SetTitle(fmt.Sprintf("Vanlig: %s", birka.vanlig))
                mBirkaVeg.SetTitle(fmt.Sprintf("Vegetarisk: %s", birka.veg))
                mEkebyhovVanlig.SetTitle(fmt.Sprintf("Vanlig: %s", ekebyhov.vanlig))
                mEkebyhovVeg.SetTitle(fmt.Sprintf("Vegetarisk: %s", ekebyhov.veg))
            case <-tQuit:
                ticker.Stop()
                return
            }
        }
    }()
}
    

func main() {
	all := flag.Bool("a", false, "Visar även vegetariska rätter")
	b := flag.Bool("b", false, "Endast Birkaskolans lunch")
	e := flag.Bool("e", false, "Visa endast Ekebyhovskolans lunch")
    tray := flag.Bool("t", false, "Visa lunchen i systemfältet")
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

	birka, err := lunch("birkaskolan")
    if err != nil {
        fmt.Println(err)
    }
	ekebyhov, err := lunch("ekebyhovskolan")
    if err != nil {
        fmt.Println(err)
    }

    if *tray {
        stop := func() {
            now := time.Now()
            fmt.Println("Bye bye", now)
        }
        systray.Run(runTray, stop)
    }

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

