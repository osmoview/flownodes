package flownodes

import "testing"

func TestReNodeClassNameValid(t *testing.T) {
	names := []string{
		"TestNode",
		"testnode",
		"node",
		"test",
		"filter",
	}

	for _, s := range names {
		if !ReNodeClassName.MatchString(s) {
			t.Errorf("%s should be valid name", s)
		}
	}
}

func TestReNodeClassNameInvalid(t *testing.T) {
	names := []string{
		"Test-Node",
		"test_node",
		"node_",
		"_test",
		"filter01",
	}

	for _, s := range names {
		if ReNodeClassName.MatchString(s) {
			t.Errorf("%s should be invalid name", s)
		}
	}
}
