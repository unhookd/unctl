package server

import (
	"net/http"

	"gopkg.in/yaml.v2"

	"github.com/spf13/cobra"
)

const MaxFileSize = 100000

func buildValuesYaml(project string, sha string) ([]byte, error) {
	var data = map[string]map[string]map[string]string{}

	data[project] = map[string]map[string]string{}
	data[project]["image"] = map[string]string{}
	data[project]["image"]["sha"] = sha

	return yaml.Marshal(data)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "Unhookd")
	w.WriteHeader(200)
}

var CmdInstastageServer = &cobra.Command{
	Use:   "instastage-server",
	Short: "Runs unhookd in insecure legacy instastage mode",
	Long: `
	  insecure legacy instastage mode allows for more direct manipulation of the values passed to a particular chart
	`,
	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		instastageServer(true)
	},
}

var CmdZeroTrustServer = &cobra.Command{
	Use:   "zero-trust-server",
	Short: "Runs unhookd in secure zero-trust model mode",
	Long: `
	  secure zero-trust model requires a known app to be present in a config, and performs additional parameter verification.
	`,
	Args: cobra.MinimumNArgs(0),
	Run: func(_ *cobra.Command, _ []string) {
		ZeroTrustServer(true)
	},
}
