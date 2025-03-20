package cmd

import (
	"strings"

	v1 "k8s.io/api/core/v1"
)

// selectCOntainers is copied from
// https://github.com/kubernetes/kubernetes/blob/aa35eff1b636f587f418f9cc16a020353735d125/staging/src/k8s.io/kubectl/pkg/cmd/set/helper.go#L29-L41

// selectContainers allows one or more containers to be matched against a string or wildcard
func selectContainers(containers []v1.Container, spec string) ([]*v1.Container, []*v1.Container) {
	out := []*v1.Container{}
	skipped := []*v1.Container{}
	for i, c := range containers {
		if selectString(c.Name, spec) {
			out = append(out, &containers[i])
		} else {
			skipped = append(skipped, &containers[i])
		}
	}
	return out, skipped
}


// selectString is copied from
// https://github.com/kubernetes/kubernetes/blob/aa35eff1b636f587f418f9cc16a020353735d125/staging/src/k8s.io/kubectl/pkg/cmd/set/helper.go#L43-L79

// selectString returns true if the provided string matches spec, where spec is a string with
// a non-greedy '*' wildcard operator.
// TODO: turn into a regex and handle greedy matches and backtracking.
func selectString(s, spec string) bool {
	if spec == "*" {
		return true
	}
	if !strings.Contains(spec, "*") {
		return s == spec
	}

	pos := 0
	match := true
	parts := strings.Split(spec, "*")
Loop:
	for i, part := range parts {
		if len(part) == 0 {
			continue
		}
		next := strings.Index(s[pos:], part)
		switch {
		// next part not in string
		case next < pos:
			fallthrough
		// first part does not match start of string
		case i == 0 && pos != 0:
			fallthrough
		// last part does not exactly match remaining part of string
		case i == (len(parts)-1) && len(s) != (len(part)+next):
			match = false
			break Loop
		default:
			pos = next
		}
	}
	return match
}
