package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/henkman/vrn"
)

func main() {
	var opts struct {
		Find        string
		Origin      string
		Destination string
	}
	flag.StringVar(&opts.Find, "f", "", "find a station")
	flag.StringVar(&opts.Origin, "o", "", "origin")
	flag.StringVar(&opts.Destination, "d", "", "destination")
	flag.Parse()
	if opts.Find != "" {
		var s vrn.Session
		if err := s.Init(); err != nil {
			log.Fatal(err)
		}
		stops, err := s.FindStop(opts.Find)
		if err != nil {
			log.Fatal(err)
		}
		for _, stop := range stops {
			raw, err := json.MarshalIndent(stop, "", "\t")
			if err != nil {
				panic(err)
			}
			fmt.Println(string(raw))
		}
		return
	}
	if opts.Origin != "" && opts.Destination != "" {
		var s vrn.Session
		if err := s.Init(); err != nil {
			log.Fatal(err)
		}
		origin, err := s.FindStop(opts.Origin)
		if err != nil {
			log.Fatal(err)
		}
		if len(origin) == 0 {
			fmt.Println("origin not found")
			return
		}
		dest, err := s.FindStop(opts.Destination)
		if err != nil {
			log.Fatal(err)
		}
		if len(dest) == 0 {
			fmt.Println("destination not found")
			return
		}
		trips, err := s.FindTrips(origin[0].Ref.Gid, dest[0].Ref.Gid)
		if err != nil {
			log.Fatal(err)
		}
		if len(trips) == 0 {
			fmt.Println("no trips found")
			return
		}
		for _, trip := range trips {
			fmt.Printf("%s: \n", trip.Duration)
			for _, leg := range trip.Legs {
				var name string
				if leg.Mode.Name != "" {
					name = leg.Mode.Name
				} else {
					name = leg.Mode.Product
				}
				fmt.Printf("\t%s\n", name)
				for _, point := range leg.Points {
					var platform string
					if point.PlatformName != "" {
						platform = "platform " + point.PlatformName
					}
					fmt.Printf("\t\t%s: '%s' %s %s\n",
						point.Usage,
						point.Name,
						platform,
						point.DateTime.Time)
				}
			}
		}
		return
	}
	flag.Usage()
}
