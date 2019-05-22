package lib

import (
	"testing"
)

func TestVerifyQueryStringParams(t *testing.T) {
	if VerifyQueryStringParams("foo") == false {
		t.Error("Expected query string param of non zero length")
	}
	if VerifyQueryStringParams("") == true {
		t.Error("Expected zero length param to return false")
	}
}

func TestGetEnv(t *testing.T) {
	if len(GetEnv("HOME", "/home/foo")) == 0 {
		t.Error("Expected to return env")
	}
}

func TestDifference(t *testing.T) {
	var found bool
	a := []string{"a", "b", "c"}
	b := []string{"a", "b", "c", "d"}
	ab := Difference(b, a)
	for _, c := range ab {
		if c == "d" {
			found = true
		}
	}
	if !found {
		t.Errorf("Expected slice contain difference in values but found %s", ab)
	}
}
