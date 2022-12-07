package fakegitserver

import (
	"io"
	"net/http"

	"github.com/spf13/cobra"
)

// CmdFakeGitServer is used by agnhost Cobra.
var CmdFakeGitServer = &cobra.Command{
	Use:   "fake-gitserver",
	Short: "Fakes a git server",
	Long: `When doing "git clone http://localhost:8000", you will clone an empty git repo named "localhost" on local.
You can also use "git clone http://localhost:8000 my-repo-name" to rename that repo.`,
	Args: cobra.MaximumNArgs(0),
	Run:  main,
}

func hello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "I am a fake git server")

}

// When doing `git clone localhost:8000`, you will clone an empty git repo named "8000" on local.
// You can also use `git clone localhost:8000 my-repo-name` to rename that repo.
func main(cmd *cobra.Command, args []string) {
	http.HandleFunc("/", hello)
	http.ListenAndServe(":8000", nil)
}
