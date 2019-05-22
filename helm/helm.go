package helm

import (
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"

	"github.com/org/unhookd/lib"
)

//helmOut, shouldNotify, err := helm.HelmBinUpgradeInstall(release, chart, sha, string(values), shouldWait, shouldDryRun)

var usrBinHelm string
var helmCmdArgs []string

func HelmBinUpgradeInstall(release string, namespace string, chart string, version string, values string, shouldWait bool, shouldDryRun bool) ([]byte, bool, error) {
	// Update helm repo to ensure we have the latest chart update
	helmRepoUpdate := exec.Command(usrBinHelm, "repo", "update", "org-s3")
	o, err := helmRepoUpdate.CombinedOutput()
	if err != nil {
		log.Printf("Repo Update Error: %v :: %s", err, o)
		return nil, false, err
	}

	if version == "instastage" {
		helmCmdArgs = []string{"upgrade", "--debug", "--install", release, "--namespace", namespace, chart, "-f", "-"}
	} else {
		helmCmdArgs = []string{"upgrade", "--debug", "--install", release, "--namespace", namespace, "--version", version, chart, "-f", "-"}
	}

	if shouldWait {
		helmCmdArgs = append(helmCmdArgs, "--wait", "--timeout", "600")
	}

	if shouldDryRun {
		helmCmdArgs = append(helmCmdArgs, "--dry-run")
	}

	helmCmd := exec.Command(usrBinHelm, helmCmdArgs...)
	log.Printf(" helm command args: %v", helmCmdArgs)

	stdin, err := helmCmd.StdinPipe()
	if err != nil {
		fmt.Println(err)
	}

	go func() {
		defer stdin.Close()
		if len(values) > 0 {
			io.WriteString(stdin, values)
		}
	}()

	valuesBeforeUpgrade := getReleaseValues(usrBinHelm, release)
	helmOutput, err := helmCmd.CombinedOutput()
	valuesAfterUpgrade := getReleaseValues(usrBinHelm, release)

	shouldNotify := shouldNotify(valuesBeforeUpgrade, valuesAfterUpgrade)
	return helmOutput, shouldNotify, err
}

func getReleaseValues(helmCommand, release string) []byte {
	helmGetValuesArgs := []string{"get", "values", release}
	helmGetValuesCmd := exec.Command(helmCommand, helmGetValuesArgs...)
	values, err := helmGetValuesCmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}
	return values
}

func shouldNotify(valuesBeforeUpgrade, valuesAfterUpgrade []byte) bool {
	if string(valuesBeforeUpgrade) != string(valuesAfterUpgrade) {
		return true
	}
	return false
}

func Setup(runHelm bool) {
	//TODO: this is where the s3 url goes
	chartsRepoUrl := lib.GetEnv("S3_CHARTS_REPO_URL", "s3://org-charts-repo")
	var output []byte
	var err error

	if runHelm {
		helmVersion := exec.Command("helm", "version", "-s", "--short", "--template", "{{ .Server.SemVer }}")
		output, err := helmVersion.CombinedOutput()
		if err != nil {
			log.Fatal(fmt.Sprintf("Helm Version Detect failure: %v %v", string(output), err))
			return
		}

		usrBinHelm = fmt.Sprintf("helm-%s", string(output))
	} else {
		usrBinHelm = "echo"
	}

	log.Println("initializing helm client")
	helmInit := exec.Command(usrBinHelm, "init", "--client-only")
	output, err = helmInit.CombinedOutput()
	if err != nil {
		log.Fatal(fmt.Sprintf("Helm Repo init failure: %v %v", string(output), err))
		return
	}

	helmRepoRemoveStable := exec.Command(usrBinHelm, "repo", "remove", "stable")
	output, err = helmRepoRemoveStable.CombinedOutput()
	if err != nil && !strings.Contains(string(output), "no repo named \"stable\" found") {
		log.Fatal(fmt.Sprintf("Helm Repo remove stable failure: %v %v", string(output), err))
		return
	}

	//TODO: put helm s3 repo add??
	log.Println("adding in-cluster charts repo")
	helmRepoAdd := exec.Command(usrBinHelm, "repo", "add", "org-s3", chartsRepoUrl)
	output, err = helmRepoAdd.CombinedOutput()
	if err != nil {
		log.Fatal(fmt.Sprintf("Helm Repo init failure: %v %v", string(output), err))
		return
	}
}
