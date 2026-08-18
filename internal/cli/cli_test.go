package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunPlanExample(t *testing.T) {
	var out, errBuf bytes.Buffer
	code := RunPlan([]string{"../../example/house.json"}, &out, &errBuf)
	if code != 0 {
		t.Fatalf("exit=%d stderr=%s", code, errBuf.String())
	}
	if !strings.Contains(out.String(), "project_duration") {
		t.Fatalf("missing duration: %s", out.String())
	}
	if !strings.Contains(out.String(), "critical_path") {
		t.Fatalf("missing critical path: %s", out.String())
	}
}

func TestRunPlanTarget(t *testing.T) {
	var out, errBuf bytes.Buffer
	code := RunPlan([]string{"--file", "../../example/pert.json", "--target", "30"}, &out, &errBuf)
	if code != 0 {
		t.Fatalf("exit=%d stderr=%s", code, errBuf.String())
	}
	if !strings.Contains(out.String(), "project_mean") {
		t.Fatalf("missing pert stats: %s", out.String())
	}
	if !strings.Contains(out.String(), "P(duration") {
		t.Fatalf("missing probability: %s", out.String())
	}
}

func TestRunPlanStdin(t *testing.T) {
	input := `{"activities":[{"id":"a","duration":2},{"id":"b","duration":3,"predecessors":["a"]}]}`
	old := stdinOverride
	stdinOverride = strings.NewReader(input)
	defer func() { stdinOverride = old }()

	var out, errBuf bytes.Buffer
	code := RunPlan([]string{}, &out, &errBuf)
	if code != 0 {
		t.Fatalf("exit=%d stderr=%s", code, errBuf.String())
	}
	if !strings.Contains(out.String(), "project_duration 5.00") {
		t.Fatalf("want duration 5.00, got:\n%s", out.String())
	}
}

func TestRunPlanRejectsCycle(t *testing.T) {
	input := `{"activities":[
	  {"id":"a","duration":1,"predecessors":["c"]},
	  {"id":"b","duration":1,"predecessors":["a"]},
	  {"id":"c","duration":1,"predecessors":["b"]}]}`
	old := stdinOverride
	stdinOverride = strings.NewReader(input)
	defer func() { stdinOverride = old }()

	var out, errBuf bytes.Buffer
	code := RunPlan([]string{}, &out, &errBuf)
	if code == 0 {
		t.Fatal("cycle should exit non-zero")
	}
	if !strings.Contains(errBuf.String(), "cycle") {
		t.Fatalf("stderr should mention cycle, got %q", errBuf.String())
	}
}
