package vrn

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

const API = "https://www.vrn.de/mngvrn"

type Gid = string

type Stop struct {
	Name    string `json:"name"`
	AnyType string `json:"anyType"`
	Object  string `json:"object"`
	MainLoc string `json:"mainLoc,omitempty"`
	Ref     struct {
		Gid    Gid    `json:"gid"`
		Place  string `json:"place"`
		Coords string `json:"coords"`
	} `json:"ref"`
}

type Trip struct {
	Duration string `json: "duration"`
	Legs     []struct {
		Points []struct {
			Name                string `json: "name"`
			NameWO              string `json: "nameWO"`
			PlatformName        string `json: "platformName"`
			PlannedPlatformName string `json: "plannedPlatformName"`
			Place               string `json: "place"`
			NameWithPlace       string `json: "nameWithPlace"`
			Usage               string `json: "usage"`
			Desc                string `json: "desc"`
			DateTime            struct {
				Date   string `json: "date"`
				Time   string `json: "time"`
				RtDate string `json: "rtDate"`
				RtTime string `json: "rtTime"`
			} `json: "dateTime"`
		} `json: "points"`
		Mode struct {
			Name        string `json: "name"`
			Number      string `json: "number"`
			Symbol      string `json: "symbol"`
			Product     string `json: "product"`
			Destination string `json: "destination"`
			Desc        string `json: "desc"`
		} `json: "mode"`
	} `json: "legs"`
}

type Session struct {
	cli http.Client
}

func (s *Session) Init() error {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return err
	}
	s.cli = http.Client{
		Jar:     jar,
		Timeout: time.Second * 10,
	}
	return nil
}

func (s *Session) FindStop(query string) ([]Stop, error) {
	const endpoint = API + "/XML_STOPFINDER_REQUEST"
	vals := url.Values{
		"anyObjFilter_sf":      []string{"0"},
		"coordOutputFormat":    []string{"EPSG:4326"},
		"locationServerActive": []string{"1"},
		"name_sf":              []string{query},
		"outputFormat":         []string{"json"},
		"type_sf":              []string{"any"},
		"vrnsuggestMacro":      []string{"vrn_suggest"},
		"w_regPrefAl":          []string{"5"},
	}.Encode()
	res, err := s.cli.Get(endpoint + "?" + vals)
	if err != nil {
		return nil, err
	}
	raw, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, err
	}
	var stopsjson struct {
		StopFinder struct {
			Points json.RawMessage `json:"points"`
		} `json:"stopFinder"`
	}
	err = json.Unmarshal(raw, &stopsjson)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err := json.Compact(&buf, stopsjson.StopFinder.Points); err != nil {
		return nil, err
	}
	points := buf.Bytes()
	var stops []Stop
	if bytes.Contains(points, []byte(`"point":{`)) {
		var rm struct {
			Stop Stop `json:"point"`
		}
		if err := json.Unmarshal(points, &rm); err != nil {
			return nil, err
		}
		stops = []Stop{rm.Stop}
	} else {
		if err := json.Unmarshal(points, &stops); err != nil {
			return nil, err
		}
	}
	return stops, nil
}

func (s *Session) FindTrips(origin, dest Gid) ([]Trip, error) {
	const endpoint = API + "/XML_TRIP_REQUEST2"
	vals := url.Values{
		"changeSpeed":          []string{"normal"},
		"coordOutputFormat":    []string{"EPSG:4326"},
		"cycleSpeed":           []string{"14"},
		"deleteITPTWalk":       []string{"0"},
		"exclMOT_15":           []string{"1"},
		"exclMOT_16":           []string{"1"},
		"excludedMeans":        []string{"checkbox"},
		"itOptionsActive":      []string{"1"},
		"itPathListActive":     []string{"1"},
		"itdTime":              []string{time.Now().Format("1504")},
		"lineRestriction":      []string{"0400"},
		"locationServerActive": []string{"1"},
		"name_destination":     []string{dest},
		"name_origin":          []string{origin},
		"outputFormat":         []string{"json"},
		"ptMacro":              []string{"true"},
		"ptOptionsActive":      []string{"1"},
		"routeType":            []string{"leasttime"},
		"trITMOTvalue":         []string{"15"},
		"type_destination":     []string{"any"},
		"type_origin":          []string{"any"},
		"useElevationData":     []string{"1"},
		"useRealtime":          []string{"1"},
		"useUT":                []string{"1"},
		"useUnifiedTickets":    []string{"1"},
	}.Encode()
	res, err := s.cli.Get(endpoint + "?" + vals)
	if err != nil {
		return nil, err
	}

	raw, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, err
	}
	var tripsjson struct {
		Trips json.RawMessage `json: "trips"`
	}
	err = json.Unmarshal(raw, &tripsjson)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err := json.Compact(&buf, tripsjson.Trips); err != nil {
		return nil, err
	}
	ctrips := buf.Bytes()
	var trips []Trip
	if bytes.Contains(ctrips, []byte(`"trip":{`)) {
		var rm struct {
			Trip Trip `json:"trip"`
		}
		if err := json.Unmarshal(ctrips, &rm); err != nil {
			return nil, err
		}
		trips = []Trip{rm.Trip}
	} else {
		if err := json.Unmarshal(ctrips, &trips); err != nil {
			return nil, err
		}
	}
	return trips, nil
}
