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
	stopsA, err := s.FindStop("Bettenbach (Bergstr), Haus Nr. 34")
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
	stopsB, err := s.FindStop("Bettenbach (Bergstr), Abzweigung")
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
	trips, err := s.FindTrips(stopsA[0], stopsB[0])
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
	for _, trip := range trips {
		fmt.Printf("%+v", trip)
	}
}
