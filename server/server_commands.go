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
	w.Header().Set("Server", "unctl")
	w.WriteHeader(200)
}

var CmdZeroTrustServer = &cobra.Command{
	Use:   "server",
	Short: "Runs the unctl server",
	Long: `
	  Runs the unctl server, exposing an endpoint where deploys defined in the unctl server config can be triggered.
	`,
	Args: cobra.MinimumNArgs(0),
	Run: func(_ *cobra.Command, _ []string) {
		ZeroTrustServer(true)
	},
}
