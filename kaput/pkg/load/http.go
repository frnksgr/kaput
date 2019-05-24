package load

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/frnksgr/kaput/kaput/pkg/config"

	"github.com/frnksgr/kaput/kaput/pkg/help"
	"github.com/gorilla/mux"
)

const helpLoad = `
/load/[count] create binary tree of requests
Where command is:
	help        this page
	count       number of requests created

For example:
	TODO
`

func init() {
	help.Add("/load", helpLoad)
}

// spawn a new Child defined by node
func spawnChild(url string, node *Node) (int, error) {
	// create request
	var buffer bytes.Buffer
	if err := encodeNode(node, &buffer); err != nil {
		return 0, err
	}
	resp, err := http.Post(url, "application/json", &buffer)
	if err != nil {
		return 0, err
	}

	// handle response
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("Unexpected response, status code: %d", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	count, err := strconv.Atoi(strings.TrimSpace(string(body)))
	if err != nil {
		return 0, err
	}
	return count, nil
}

// get count from url /load/{count:\\d{0,4}}
func nodeCount(r *http.Request) int {
	value := mux.Vars(r)["count"]
	count, err := strconv.Atoi(value)
	if err != nil { // should not happen
		panic(err)
	}
	return count
}

// PostHandler handle post requests
// accepting payloads application/json (inner nodes)
// application/x-www-form-urlencoded (root node)
func PostHandler(w http.ResponseWriter, r *http.Request) {
	nodeCount := nodeCount(r)
	var node *Node
	var err error

	// construct node
	switch r.Header.Get("Content-Type") {
	case "application/json": // inner node
		node, err = decodeNode(r.Body)
		if err != nil {
			panic(err)
		}
	case "application/x-www-form-urlencoded": //root node
		load, err := decodeLoad(r)
		if err != nil {
			panic(err)
		}
		node = &Node{
			Index: 1,
			Load:  load,
		}
	}

	// get a node specific logger
	logger := node.logger()

	// execute workload
	node.Run(logger)

	// spawn children
	children := make([]*Node, 0, 2) // actual children to be spawned
	statusCode := http.StatusOK     // default
	childCount := 0                 // count successfully spawned children

	// figure out children to be spawned
	for i := 0; i < 2; i++ {
		child := node.child(i)
		if child.Index > nodeCount {
			break
		}
		children = append(children, child)
	}

	if len(children) > 0 { // some children to spawn
		url := fmt.Sprintf("%s://%s:%s/load/%d", config.Data.Calling.Protocol,
			config.Data.Calling.Domain, config.Data.Calling.Port, nodeCount)

		type childResult struct {
			count int
			err   error
		}

		channel := make(chan childResult, len(children))
		for _, child := range children {
			go func(child *Node) {
				count, err := spawnChild(url, child)
				channel <- childResult{count, err}
			}(child)
		}
		for i := 0; i < len(children); i++ {
			result := <-channel
			if result.err != nil {
				logger.Println(result.err)
				statusCode = http.StatusInternalServerError
			}
			childCount += result.count
		}
	}

	// create response
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, childCount)
}
