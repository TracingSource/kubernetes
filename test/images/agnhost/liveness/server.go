// A simple server that is alive for 10 seconds, then reports unhealthy for
// the rest of its (hopefully) short existence.

package liveness

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/spf13/cobra"
)

// CmdLiveness is used by agnhost Cobra.
var CmdLiveness = &cobra.Command{
	Use:   "liveness",
	Short: "Starts a server that is alive for 10 seconds",
	Long:  "A simple server that is alive for 10 seconds, then reports unhealthy for the rest of its (hopefully) short existence",
	Args:  cobra.MaximumNArgs(0),
	Run:   main,
}

func main(cmd *cobra.Command, args []string) {
	started := time.Now()
	http.HandleFunc("/started", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		data := (time.Since(started)).String()
		w.Write([]byte(data))
	})
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		duration := time.Since(started)
		if duration.Seconds() > 10 {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf("error: %v", duration.Seconds())))
		} else {
			w.WriteHeader(200)
			w.Write([]byte("ok"))
		}
	})
	http.HandleFunc("/redirect", func(w http.ResponseWriter, r *http.Request) {
		loc, err := url.QueryUnescape(r.URL.Query().Get("loc"))
		if err != nil {
			http.Error(w, fmt.Sprintf("invalid redirect: %q", r.URL.Query().Get("loc")), http.StatusBadRequest)
			return
		}
		http.Redirect(w, r, loc, http.StatusFound)
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
