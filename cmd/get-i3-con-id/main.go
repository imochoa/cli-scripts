package main

// https://eager.io/blog/go-and-json/

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
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
	Y      int `json:"x"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

type containerNode struct {
	Id                 float64          `json:"id"`
	Type               string           `json:"type"`
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
	//  1) run i3-msg -t get_tree
	out, err := exec.Command("/usr/bin/i3-msg", "-t", "get_tree").Output()
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(out, rootContainer)
	if err != nil {
		log.Fatal(err)
	}

}

type nodeFilter func(*containerNode) bool

func noFilter(c *containerNode) bool {
	return true
}

func recursiveI3ContainerSearch(c containerNode, filter nodeFilter, childConts *[]containerNode) {
	// Better names!

	if filter(&c) {
		*childConts = append(*childConts, c)
	} else {
		for _, n := range c.Nodes {
			recursiveI3ContainerSearch(n, filter, childConts)
		}
	}
}

func main() {

	// Get i3 tree
	var rootContainer containerNode
	saveI3Tree(&rootContainer)

	// find matches

	// //  1) run i3-msg -t get_tree
	// out, err := exec.Command("/usr/bin/i3-msg", "-t", "get_tree").Output()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// // fmt.Printf("The data is %s\n", out)

	// var rootContainer containerNode

	// //  2) parse JSON output looking for a class

	// err = json.Unmarshal(out, &rootContainer)

	//  3) return ID so that we can run `i3-msg move '[con_id="xxx"] workspace 4'`
	var lenGuess int = 20
	childConts := make([]containerNode, 0, lenGuess)
	// var childConts []containerNode

	// func filterLeafContainers(c *containerNode) bool {
	// 	return len(c.Nodes) == 0 && c.Type == "con"
	// }

	filterContClass := func(c *containerNode) bool {
		return len(c.Nodes) == 0 && c.Type == "con" && c.WindowProperties.Class == "Gnome-terminal"
	}

	recursiveI3ContainerSearch(rootContainer, filterContClass, &childConts)

	fmt.Println(childConts)
}
