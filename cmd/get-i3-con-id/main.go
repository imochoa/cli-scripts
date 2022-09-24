package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"
)

type windowProperties struct {
	Class    string `json:"class"`
	Instance string `json:"instance"`
	Machine  string `json:"machine"`
	Title    string `json:"title"`
	// 		"transient_for": null
}

type windowPosition struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

type containerNode struct {
	Id                 int              `json:"id"`
	Type               string           `json:"type"` // "con" "workspace"
	Name               string           `json:"name"`
	Nodes              []containerNode  `json:"nodes"`
	WindowProperties   windowProperties `json:"window_properties"`
	Orientation        string           `json:"orientation"`
	ScratchpadState    string           `json:"scratchpad_state"`
	Percent            float64          `json:"percent"`
	Urgent             bool             `json:"urgent"`
	Focused            bool             `json:"focused"`
	Output             string           `json:"output"`
	Layout             string           `json:"layout"`
	WorkspaceLayout    string           `json:"workspace_layout"`
	LastSplitLayout    string           `json:"last_split_layout"`
	Border             string           `json:"border"`
	CurrentBorderWidth int              `json:"current_border_width"`
	Window             int              `json:"window"`
	WindowType         string           `json:"window_type"`
	Sticky             bool             `json:"sticky"`
	Rect               windowPosition   `json:"rect"`
	DecoRect           windowPosition   `json:"deco_rect"`
	Geometry           windowPosition   `json:"geometry"`
	Floating           string           `json:"floating"`
	WindowIconPadding  int              `json:"window_icon_padding"`
	// FullscreenMode     bool   *not always!*          `json:"fullscreen_mode"`

	// "marks": [],
	// "floating_nodes": [],
	// "focus": [],
	// "swallows": []
}

func saveI3Tree(rootContainer *containerNode) {
	// run "i3-msg -t get_tree" and save the results in `rootContainer`
	out, err := exec.Command("/usr/bin/i3-msg", "-t", "get_tree").Output()
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(out, rootContainer)
	if err != nil {
		log.Fatal(err)
	}

}

func moveContainerToWorkspace(containerId int, workspace int) error {
	// run "i3-msg -t get_tree" and save the results in `rootContainer`

	out, err := exec.Command(
		"/usr/bin/i3-msg",
		fmt.Sprintf("[con_id=\"%d\"]", containerId),
		"move",
		"workspace",
		strconv.Itoa(workspace),
	).Output()

	stdout := string(out)
	fmt.Println(stdout)
	return err
}

type nodeFilter func(*containerNode) bool

// func noFilter(c *containerNode) bool {
// 	return true
// }

func recursiveI3ContainerSearch(c containerNode, filter nodeFilter, childConts *[]containerNode) {
	// Look at each

	if filter(&c) {
		*childConts = append(*childConts, c)
	} else {
		for _, n := range c.Nodes {
			recursiveI3ContainerSearch(n, filter, childConts)
		}
	}
}

func main() {
	// Move the first matching container to workspace `workspace`
	appClass := flag.String("class", "", "Application's WM_CLASS")
	workspaceNumber := flag.Int("workspace", 1, "Workspace to move the application to")
	flag.Parse()

	if *appClass == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	fmt.Printf("Attempting to move WM_CLASS %s to workspace %d ...\n", *appClass, *workspaceNumber)

	// var wIdx string = strconv.Itoa(*workspaceNumber)

	filterContClass := func(c *containerNode) bool {
		return (len(c.Nodes) == 0 &&
			c.Type == "con" &&
			c.WindowProperties.Class == *appClass)
	}

	var rootContainer containerNode
	//  3) return ID so that we can run
	var lenGuess int = 20
	childConts := make([]containerNode, 0, lenGuess)

	var maxStartupSeconds float64 = 2.0
	pollInterval := 100 * time.Millisecond
	start := time.Now()
	for time.Since(start).Seconds() < maxStartupSeconds {

		// Get i3 tree
		saveI3Tree(&rootContainer)

		childConts = childConts[:0]
		recursiveI3ContainerSearch(rootContainer, filterContClass, &childConts)

		if len(childConts) > 0 {

			err := moveContainerToWorkspace(childConts[0].Id, *workspaceNumber)
			if err != nil {
				log.Fatal(err)
			}
			os.Exit(0)
		}

		fmt.Print('.')
		time.Sleep(pollInterval)
		// fmt.Println(time.Now())
	}

	fmt.Printf("Did not find a WM_CLASS %s container...\n", *appClass)
	os.Exit(1)

}
