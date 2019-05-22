package helm

import (
	"strings"
	"testing"
)

func TestShouldNotify(t *testing.T) {
	values1 := []byte{1, 2, 3}
	values2 := []byte{2, 3}

	if shouldNotify(values1, values2) == false {
		t.Errorf("Expected shouldNotify to return true when values are different")
	}

	values1 = []byte{1, 2, 3}
	values2 = []byte{1, 2, 3}

	if shouldNotify(values1, values2) == true {
		t.Errorf("Expected shouldNotify to return false when values are the same")
	}
}

func TestGetReleaseValues(t *testing.T) {
	Setup(false)
	commandToRun := "echo"
	actualOutput := string(getReleaseValues(commandToRun, "fancy_release"))
	expectedOutput := "get values fancy_release"

	if strings.TrimRight(actualOutput, "\n") != expectedOutput {
		error := `Expected getReleaseValues to execute the command to run with correct args and return the output of the command. 
				Expected output: %v
				Actual output: %v`
		t.Errorf(error, expectedOutput, actualOutput)
	}
}
