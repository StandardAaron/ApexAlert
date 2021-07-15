// +build linux,amd64 windows

package main

/* to-do:
- Log to-do as proper issues and remove to-do :D
- test windows and other *nix wm/de (and eventually OSX)
- add config file and options for:
	- fave/exclude maps/modes
	- blackout period
	- afk-mode
	- throttling
	- user name/platform (for BatttlePass tracking)
- Ranked notifications (pull time from some other source)
- Keep a list of maps and auto-pop new ones
	(eg: alert on brand new arenas maps)
- News feed parser with alerts on new news items
*/

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/hashicorp/hcl/v2/hclsimple"
)

// Current and Next are reusable types in their parent structs
// which saves us some lines :D
type CurrentMap struct {
	Start             int64  `json:"start"`
	End               int64  `json:"end"`
	ReadableDateStart string `json:"readableDate_start"`
	ReadableDateEnd   string `json:"readableDate_end"`
	Map               string `json:"map"`
	DurationInSecs    int    `json:"DurationInSecs"`
	DurationInMinutes int    `json:"DurationInMinutes"`
	RemainingSecs     int    `json:"remainingSecs"`
	RemainingMins     int    `json:"remainingMins"`
	RemainingTimer    string `json:"remainingTimer"`
}

type NextMap struct {
	Start             int64  `json:"start"`
	End               int64  `json:"end"`
	ReadableDateStart string `json:"readableDate_start"`
	ReadableDateEnd   string `json:"readableDate_end"`
	Map               string `json:"map"`
	DurationInSecs    int    `json:"DurationInSecs"`
	DurationInMinutes int    `json:"DurationInMinutes"`
}

type GameModeState struct {
	CurrentMap CurrentMap `json:"current"`
	NextMap    NextMap    `json:"next"`
}

/*
	Ranked's structure is *technically* different (see end of comment),
	but Go's Unmarshal (conveiently) nils out missing feilds,
	allowing us to simply re-use the same Current/Next types
	that we already defined.
	Technically, it's:
		Current struct {
			Map string `json:"map"`
		} `json:"current"`
		Next struct {
			Map string `json:"map"`
		} `json:"next"`
*/
type BattleRoyale GameModeState
type Arenas GameModeState
type Ranked GameModeState

type ApexMaps struct {
	BattleRoyale BattleRoyale `json:"battle_royale"`
	Arenas       Arenas       `json:"arenas"`
	Ranked       Ranked       `json:"ranked"`
}

type Configuration struct {
	ApiKey string `hcl:"api_key"`
}

func readConfig() (config Configuration) {
	err := hclsimple.DecodeFile("config.hcl", nil, &config)
	if err != nil {
		panic(fmt.Sprintf("Failed to load configuration: %s", err))
	}
	return
}

func getApiUrl() (apiUrl string) {
	//Make sure you request a key from https://apexlegendsapi.com/
	c := readConfig()
	apiUrl += "https://api.mozambiquehe.re/maprotation?version=2&auth=" + string(c.ApiKey)
	return
}

func getMapRotationData() ApexMaps {
	//Hit the endpoint
	apiURL := getApiUrl()
	resp, err := http.Get(apiURL)
	if err != nil {
		fmt.Println("http get error: ", err)
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	jsonStream, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Got an error reading body bytestream", err)
	}
	var apexMapData ApexMaps
	jsonerr := json.Unmarshal(jsonStream, &apexMapData)
	if jsonerr != nil {
		fmt.Println("error Unmarshalling data stream: \n", err)
	}
	return apexMapData
}

func mapWatcher(mapType string) {
	for {
		apex := getMapRotationData()
		var current CurrentMap
		var next NextMap
		switch mapType {
		case "Arenas":
			current = apex.Arenas.CurrentMap
			next = apex.Arenas.NextMap
		case "Battle Royale":
			current = apex.BattleRoyale.CurrentMap
			next = apex.BattleRoyale.NextMap
		default:
			panic("no type given to mapWatcher.")
		}
		Head := mapType + " is now on " + current.Map + "!"
		Subhead := "Next up is " + next.Map + " in " + strconv.Itoa(current.RemainingMins) + " Mins."
		issueAlert(Head, Subhead, notificationAsset)
		//use seconds remainging to sleep, and add a buffer so we dont get caught in a scenario
		//where RemainingSecs is 0 and we call the API again immediately (and dupe alerts)
		time.Sleep(time.Duration(current.RemainingSecs+5) * time.Second)
	}
}

func issueAlert(head string, subhead string, icon string) {
	fmt.Println(head, subhead, icon) //debug
	err := beeep.Notify(head, subhead, icon)
	if err != nil {
		panic("failed to notify")
	}
}

var notificationAsset string = "assets/apexAlertIcon.jpg"

func main() {
	/*
		I'd like to clean this up,
		we want the mapWatcher goroutines to be long-lived,
		but there's currently no exit criteria so this is
		"runs till Ctrl-C behaviour."

		I have also explored using recievers/Go methods
		for calling {gameModeState}.mapWatcher(),
		but that would require a redesign using timers instead of sleep.
	*/
	var wg sync.WaitGroup
	var gameModes = []string{"Arenas", "Battle Royale"}

	for i := 0; i < len(gameModes); i++ {
		wg.Add(1)
		go mapWatcher(gameModes[i])
		//wait a couple seconds just to be nice on the API :D
		time.Sleep(2 * time.Second)
	}
	wg.Wait()
}