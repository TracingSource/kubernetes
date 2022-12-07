package entrypoint

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// CmdEntrypointTester is used by agnhost Cobra.
var CmdEntrypointTester = &cobra.Command{
	Use:   "entrypoint-tester",
	Short: "Prints the args it's passed and exits",
	Long:  "Prints the args it's passed and exits.",
	Run:   main,
}

// This program prints all the executable's arguments and exits.
func main(cmd *cobra.Command, args []string) {
	// Some of the entrypoint-tester related tests overrides agnhost's default entrypoint
	// with agnhost-2, and this function's args will only contain the subcommand's
	// args (./agnhost entrypoint-tester these args), but we need to print *all* the
	// args, which is why os.Args should be printed instead.
	fmt.Printf("%v\n", os.Args)
	os.Exit(0)
}
