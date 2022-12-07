package connect

import (
	"fmt"
	"net"
	"os"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

// CmdConnect is used by agnhost Cobra.
var CmdConnect = &cobra.Command{
	Use:   "connect [host:port]",
	Short: "Attempts a TCP connection and returns useful errors",
	Long: `Tries to open a TCP connection to the given host and port. On error it prints an error message prefixed with a specific fixed string that test cases can check for:

* UNKNOWN - Generic/unknown (non-network) error (eg, bad arguments)
* TIMEOUT - The connection attempt timed out
* DNS - An error in DNS resolution
* REFUSED - Connection refused
* OTHER - Other networking error (eg, "no route to host")`,
	Args: cobra.ExactArgs(1),
	Run:  main,
}

var timeout time.Duration

func init() {
	CmdConnect.Flags().DurationVar(&timeout, "timeout", time.Duration(0), "Maximum time before returning an error")
}

func main(cmd *cobra.Command, args []string) {
	dest := args[0]

	// Redundantly parse and resolve the destination so we can return the correct
	// errors if there's a problem.
	if _, _, err := net.SplitHostPort(dest); err != nil {
		fmt.Fprintf(os.Stderr, "UNKNOWN: %v\n", err)
		os.Exit(1)
	}
	if _, err := net.ResolveTCPAddr("tcp", dest); err != nil {
		fmt.Fprintf(os.Stderr, "DNS: %v\n", err)
		os.Exit(1)
	}

	conn, err := net.DialTimeout("tcp", dest, timeout)
	if err == nil {
		conn.Close()
		os.Exit(0)
	}
	if opErr, ok := err.(*net.OpError); ok {
		if opErr.Timeout() {
			fmt.Fprintf(os.Stderr, "TIMEOUT\n")
			os.Exit(1)
		} else if syscallErr, ok := opErr.Err.(*os.SyscallError); ok {
			if syscallErr.Err == syscall.ECONNREFUSED {
				fmt.Fprintf(os.Stderr, "REFUSED\n")
				os.Exit(1)
			}
		}
	}

	fmt.Fprintf(os.Stderr, "OTHER: %v\n", err)
	os.Exit(1)
}
