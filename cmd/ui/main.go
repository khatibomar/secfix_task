package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	. "modernc.org/tk9.0"
	_ "modernc.org/tk9.0/themes/azure"
)

type Response struct {
	OsVersions struct {
		OsQueryVersion string `json:"os_query_version"`
		OsVersion      string `json:"os_version"`
		SnapshotAt     string `json:"snapshot_at"`
	} `json:"os_versions"`
	Apps struct {
		Installed  []string `json:"installed"`
		SnapshotAt string   `json:"snapshot_at"`
	} `json:"apps"`
}

func main() {
	_ = ActivateTheme("azure light")
	var scroll *TScrollbarWidget

	t := Text(Wrap("none"), Setgrid(true), Yscrollcommand(func(e *Event) { e.ScrollSet(scroll) }))
	scroll = TScrollbar(Command(func(e *Event) { e.Yview(t) }))
	Grid(t, Sticky("news"), Pady("2m"), Ipadx("1m"), Ipady("1m"))
	Grid(scroll, Row(0), Column(1), Sticky("nes"), Pady("2m"))
	GridRowConfigure(App, 0, Weight(1))
	GridColumnConfigure(App, 0, Weight(1))
	Grid(TExit(), Padx("1m"), Pady("2m"), Ipadx("1m"), Ipady("1m"))

	resp, err := http.Get("http://localhost:4000/v1/latest_data")
	if err != nil {
		t.Replace("1.0", "end", fmt.Sprintf("Error reading response: %v\n", err))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Replace("1.0", "end", fmt.Sprintf("Error reading response: %v\n", err))
	}

	var data Response
	if err := json.Unmarshal(body, &data); err != nil {
		t.Replace("1.0", "end", fmt.Sprintf("Error parsing JSON: %v\n", err))
	}

	output := fmt.Sprintf("OS Query Version: %s\nOS Version: %s\nSnapshot Time: %s\n\nInstalled Apps at %s:\n",
		data.OsVersions.OsQueryVersion,
		data.OsVersions.OsVersion,
		data.OsVersions.SnapshotAt,
		data.Apps.SnapshotAt)

	for _, app := range data.Apps.Installed {
		output += fmt.Sprintf("- %s\n", app)
	}

	t.Insert("end", output)
	App.Center().Wait()
}
