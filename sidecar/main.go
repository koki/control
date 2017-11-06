package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"

	"github.com/golang/glog"
	"github.com/koki/control/pkg/koki"
)

func writePodStatus(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		status := map[string]interface{}{
			"cpu_usage": 0.65,
			"mem_usage": 0.76,
		}
		s, err := json.Marshal(status)
		if err != nil {
			glog.Error(err)
			return
		}
		fmt.Fprintf(w, string(s))
	default:
		fmt.Fprintf(w, "unsupported verb %s", r.Method)
	}
}

func main() {
	var port = flag.Int("port", koki.DefaultSidecarPort, "which port should the sidecar server use?")
	flag.Parse()

	http.HandleFunc("/podstatus", writePodStatus)
	glog.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
