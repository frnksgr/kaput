package load

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/frnksgr/kaput/kaput/pkg/config"
	"github.com/google/uuid"

	"github.com/frnksgr/kaput/kaput/pkg/help"
	"github.com/gorilla/mux"
)

const helpLoad = `
/load/[count] create binary tree of requests
Each request triggering a asynchronous workload execution
Workloads are specified by json
Where command is:
	count       number of requests created

For example:
	curl -H "Content-Type: application/json" -d '{"timeout":5000, "tasks": [{"cpu", 0.75]}' \
	  <base url>/load/100 
	Create 100 nodes/requests each running a cpu load of 75% for 5000 milliseconds
`

func init() {
	help.Add("/load", helpLoad)
}

func (n *node) spawn() (int, error) {
	// create URL
	url := fmt.Sprintf("%s://%s:%s/load/%d", config.InternalCallingProtocol(),
		config.InternalCallingDomain(), config.InternalCallingPort(), n.Count)
	// create body
	var reqBody bytes.Buffer
	if err := encodeNode(n, &reqBody); err != nil {
		return 0, err
	}
	// create request object
	req, err := http.NewRequest("POST", url, &reqBody)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Node-Type", "inner")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}

	// handle response
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("Unexpected response, status code: %d", resp.StatusCode)
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	count, err := strconv.Atoi(strings.TrimSpace(string(respBody)))
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

// Handler handle post requests
// requiring payloads application/json
func Handler(w http.ResponseWriter, r *http.Request) {
	var n *node
	var err error

	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	n, err = decodeNode(r.Body)
	if err != nil {
		panic(err)
	}

	if r.Header.Get("Node-Type") != "inner" {
		// root node
		n.UUID = uuid.New()
		n.Count = nodeCount(r)
		n.Index = 1
		n.Level = 0
	}

	// get a node specific logger
	logger := n.logger()

	// execute workload
	n.run(logger)

	// spawn children
	children := make([]*node, 0, 2) // actual children to be spawned
	statusCode := http.StatusOK     // default
	count := 1                      // count successfully spawned children

	// figure out children to be spawned
	for i := 0; i < 2; i++ {
		child := n.child(i)
		if child.Index > n.Count {
			break
		}
		children = append(children, child)
	}

	if len(children) > 0 { // some children to spawn
		type childResult struct {
			count int
			err   error
		}

		channel := make(chan childResult, len(children))
		for _, child := range children {
			go func(child *node) {
				count, err := child.spawn()
				channel <- childResult{count, err}
			}(child)
		}
		for i := 0; i < len(children); i++ {
			result := <-channel
			if result.err != nil {
				logger.Println(result.err)
				statusCode = http.StatusInternalServerError
			}
			count += result.count
		}
	}

	// create response
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintln(w, count)
}
