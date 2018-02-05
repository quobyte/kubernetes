package web

import (
	"fmt"
	"net/http"
	"operator/pkg/resourcehandler"
)

var (
	urlHandler = map[string]fn{
		"nodeList": nodeListHandler,
		"menu":     menuHandler,
	}
)

type fn func(w http.ResponseWriter, r *http.Request)

const (
	//port to run webserver for quobyte-operator
	port string = ":7878"
)

func nodeListHandler(w http.ResponseWriter, r *http.Request) {
	nodes, err := resourcehandler.GetQuobyteNodes()
	if err != nil {
		panic(err)
	}
	for _, node := range nodes.Items {
		fmt.Fprintf(w, "Node: %s\n", node.Name)
	}
}

func menuHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "hello menu")
}

func getHandlerFunc(path string) fn {
	v := urlHandler[path]
	return v
}

//StartWebServer Starts webserver on port
func StartWebServer() {
	fmt.Printf("Starting server on port: %s", port)
	http.HandleFunc("/listNodes", nodeListHandler)
	http.HandleFunc("/menu", menuHandler)
	http.ListenAndServe(port, nil)
}
