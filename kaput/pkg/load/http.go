package load

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

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

func decodeLoad(r io.Reader) (*Load, error) {
	var l Load
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&l); err != nil {
		return nil, err
	}
	return &l, nil
}

func encodeLoad(l *Load, w io.Writer) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(l)
}

func createChild(url string, n *Node) error {
	indexedURL := fmt.Sprintf("%s?index=%d", url, n.Index)
	var buffer bytes.Buffer
	if err := encodeLoad(n.Load, &buffer); err != nil {
		return err
	}
	resp, err := http.Post(indexedURL, "application/json", &buffer)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Unexpected response, status code: %d", resp.StatusCode)
	}
	return nil
}

// Handler server load requests
func Handler(w http.ResponseWriter, r *http.Request) {
	value := mux.Vars(r)["count"]
	count, err := strconv.Atoi(value)
	if err != nil { // should not happen
		log.Panic(err)
	}

	// only root request should have no index
	// and is defaulted to 1
	value = r.URL.Query().Get("index")
	index, err := strconv.Atoi(value)
	if err != nil {
		// defaults to 1
		index = 1
	}

	load, err := decodeLoad(r.Body)
	if err != nil {
		log.Panic(err)
	}
	node := &Node{
		Index: index,
		Load:  load,
	}

	// execute workload
	node.Run()

	statusCode := http.StatusOK
	url := fmt.Sprintf("%s://%s:%s/load/%d",
		config.Data.Calling.Protocol,
		config.Data.Calling.Domain,
		config.Data.Calling.Port,
		count)

	switch {
	case count >= 2*node.Index+1: // inner node with two children
		lc, rc := make(chan error), make(chan error)
		ln := &Node{
			Index: 2 * node.Index,
			Load:  node.Load,
		}
		rn := &Node{
			Index: 2*node.Index + 1,
			Load:  node.Load,
		}

		go func() {
			lc <- createChild(url, ln)
		}()
		go func() {
			rc <- createChild(url, rn)
		}()

		for i := 0; i < 2; i++ {
			select {
			case err := <-lc:
				if err != nil {
					log.Println(err)
					statusCode = http.StatusInternalServerError
				}
			case err := <-rc:
				if err != nil {
					log.Println(err)
					statusCode = http.StatusInternalServerError
				}
			}
		}

	case count == 2*node.Index: // inner node with one child
		lc := make(chan error)
		ln := &Node{
			Index: 2 * node.Index,
			Load:  node.Load,
		}

		go func() {
			lc <- createChild(url, ln)
		}()
		err := <-lc
		if err != nil {
			log.Println(err)
			statusCode = http.StatusInternalServerError
		}
	default: // leaf
		node.Println("is Leaf")
	}

	w.WriteHeader(statusCode)
}
