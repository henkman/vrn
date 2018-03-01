package vrn

import (
	"fmt"
	"testing"
)

func TestFindSingleStop(t *testing.T) {
	var s Session
	if err := s.Init(); err != nil {
		t.Error(err)
		t.Fail()
		return
	}
	stops, err := s.FindStop("Bettenbach (Bergstr), Haus Nr. 34")
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
	if len(stops) == 0 || stops[0].MainLoc != "Mörlenbach" {
		t.Error("wrong stop")
		t.Fail()
		return
	}
}

func TestFindMultipleStops(t *testing.T) {
	var s Session
	if err := s.Init(); err != nil {
		t.Error(err)
		t.Fail()
		return
	}
	stops, err := s.FindStop("Bettenbach (Bergstr)")
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
	if len(stops) <= 1 {
		t.Error("no stops")
		t.Fail()
		return
	}
}

func TestFindTrips(t *testing.T) {
	var s Session
	if err := s.Init(); err != nil {
		t.Error(err)
		t.Fail()
		return
	}
	origin := "de:07314:2066"
	dest := "de:08222:2417"
	trips, err := s.FindTrips(origin, dest)
	if err != nil {
		t.Error(err)
		t.Fail()
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
}
