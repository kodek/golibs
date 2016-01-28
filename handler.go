package healthz

import (
	"fmt"
	"net/http"
)

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}

func init() {
	http.HandleFunc("/healthz", healthzHandler)
}
