package nosnatproxy

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"k8s.io/component-base/logs"
)

// CmdNoSnatTestProxy is used by agnhost Cobra.
var CmdNoSnatTestProxy = &cobra.Command{
	Use:   "no-snat-test-proxy",
	Short: "Creates a proxy for the /checknosnat endpoint",
	Long:  `Creates the /checknosnat endpoint which proxies the request to the given target (/checknosnat?target=target_ip&ips=ip1,ip2) and returns its response, or a 500 response on error.`,
	Args:  cobra.MaximumNArgs(0),
	Run:   main,
}

var port string

func init() {
	CmdNoSnatTestProxy.Flags().StringVar(&port, "port", "31235", "The port to serve /checknosnat endpoint on.")
}

// This Pod's /checknosnat takes `target` and `ips` arguments, and queries {target}/checknosnat?ips={ips}

type masqTestProxy struct {
	Port string
}

func main(cmd *cobra.Command, args []string) {
	m := &masqTestProxy{
		Port: port,
	}

	logs.InitLogs()
	defer logs.FlushLogs()

	if err := m.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func (m *masqTestProxy) Run() error {
	// register handler
	http.HandleFunc("/checknosnat", checknosnat)

	// spin up the server
	return http.ListenAndServe(":"+m.Port, nil)
}

func checknosnatURL(pip, ips string) string {
	return fmt.Sprintf("http://%s/checknosnat?ips=%s", pip, ips)
}

func checknosnat(w http.ResponseWriter, req *http.Request) {
	url := checknosnatURL(req.URL.Query().Get("target"), req.URL.Query().Get("ips"))
	resp, err := http.Get(url)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "error querying %q, err: %v", url, err)
	} else {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "error reading body of response from %q, err: %v", url, err)
		} else {
			// Respond the same status code and body as /checknosnat on the internal Pod
			w.WriteHeader(resp.StatusCode)
			w.Write(body)
		}
	}
	resp.Body.Close()
}
