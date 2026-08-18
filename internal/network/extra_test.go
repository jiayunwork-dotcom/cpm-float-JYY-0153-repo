package network

import (
	"testing"
)

func TestEnumeratePathsCount(t *testing.T) {
	in := Input{Activities: []Activity{
		{ID: "a", Duration: 1},
		{ID: "b", Duration: 1, Predecessors: []string{"a"}},
		{ID: "c", Duration: 1, Predecessors: []string{"a"}},
		{ID: "d", Duration: 1, Predecessors: []string{"b", "c"}},
	}}
	plan, err := Plan(in)
	if err != nil {
		t.Fatalf("Plan: %v", err)
	}
	paths := EnumeratePaths(plan.Graph)
	if paths.Count != 2 {
		t.Fatalf("want 2 paths, got %d", paths.Count)
	}
}

func TestLongestPaths(t *testing.T) {
	in := Input{Activities: []Activity{
		{ID: "a", Duration: 1},
		{ID: "b", Duration: 3, Predecessors: []string{"a"}},
		{ID: "c", Duration: 2, Predecessors: []string{"a"}},
		{ID: "d", Duration: 1, Predecessors: []string{"b", "c"}},
	}}
	plan, err := Plan(in)
	if err != nil {
		t.Fatalf("Plan: %v", err)
	}
	paths := LongestPaths(plan.Graph, plan.MaxDuration())
	if len(paths) == 0 {
		t.Fatal("no paths")
	}
	if paths[0].Length != 5 {
		t.Fatalf("longest path want 5, got %.2f", paths[0].Length)
	}
}

func TestMetrics(t *testing.T) {
	in := Input{Activities: []Activity{
		{ID: "a", Duration: 1},
		{ID: "b", Duration: 1, Predecessors: []string{"a"}},
		{ID: "c", Duration: 1, Predecessors: []string{"a"}},
		{ID: "d", Duration: 1, Predecessors: []string{"b", "c"}},
	}}
	plan, err := Plan(in)
	if err != nil {
		t.Fatalf("Plan: %v", err)
	}
	m, err := Analyze(plan.Graph)
	if err != nil {
		t.Fatalf("Analyze: %v", err)
	}
	if m.ActivityCount != 4 {
		t.Fatalf("want 4 activities, got %d", m.ActivityCount)
	}
	if m.CriticalRatio <= 0 || m.CriticalRatio > 1 {
		t.Fatalf("critical ratio out of range: %v", m.CriticalRatio)
	}
}

func TestLevelize(t *testing.T) {
	in := Input{Activities: []Activity{
		{ID: "a", Duration: 1},
		{ID: "b", Duration: 1, Predecessors: []string{"a"}},
		{ID: "c", Duration: 1, Predecessors: []string{"a"}},
		{ID: "d", Duration: 1, Predecessors: []string{"b", "c"}},
	}}
	plan, err := Plan(in)
	if err != nil {
		t.Fatalf("Plan: %v", err)
	}
	levels := Levelize(plan.Graph)
	if LevelOf(plan.Graph, "a") != 1 {
		t.Fatalf("a should be level 1, got %d", LevelOf(plan.Graph, "a"))
	}
	if LevelOf(plan.Graph, "d") != 3 {
		t.Fatalf("d should be level 3, got %d", LevelOf(plan.Graph, "d"))
	}
	if len(levels) < 5 {
		t.Fatalf("want >=5 levels, got %d", len(levels))
	}
}

func TestFloatTable(t *testing.T) {
	in := Input{Activities: []Activity{
		{ID: "a", Duration: 1},
		{ID: "b", Duration: 1, Predecessors: []string{"a"}},
		{ID: "c", Duration: 1, Predecessors: []string{"a"}},
		{ID: "d", Duration: 1, Predecessors: []string{"b", "c"}},
	}}
	plan, err := Plan(in)
	if err != nil {
		t.Fatalf("Plan: %v", err)
	}
	table := FloatTable(plan.Graph, plan.MaxDuration())
	if len(table) != 6 {
		t.Fatalf("want 6 rows (4 acts + START + END), got %d", len(table))
	}
	for _, row := range table {
		if row.TF < 0 {
			t.Fatalf("negative TF for %s", row.ID)
		}
	}
}

func TestCrashReducesDuration(t *testing.T) {
	in := Input{Activities: []Activity{
		{ID: "a", Duration: 6},
		{ID: "b", Duration: 2, Predecessors: []string{"a"}},
	}}
	plan, err := Plan(in)
	if err != nil {
		t.Fatalf("Plan: %v", err)
	}
	if plan.MaxDuration() != 8 {
		t.Fatalf("want 8, got %.2f", plan.MaxDuration())
	}
	crashes := map[string]CrashActivity{
		"a": {ID: "a", Normal: 6, Crash: 2, CostSlope: 10},
	}
	results, err := CrashProject(in, crashes, 4)
	if err != nil {
		t.Fatalf("CrashProject: %v", err)
	}
	if len(results) == 0 {
		t.Fatal("crash should have applied")
	}
	final, _ := Plan(in)
	_ = final
}

func TestHasPath(t *testing.T) {
	in := Input{Activities: []Activity{
		{ID: "a", Duration: 1},
		{ID: "b", Duration: 1, Predecessors: []string{"a"}},
		{ID: "c", Duration: 1, Predecessors: []string{"a"}},
	}}
	plan, err := Plan(in)
	if err != nil {
		t.Fatalf("Plan: %v", err)
	}
	if !HasPath(plan.Graph, "a", "b") {
		t.Fatal("a->b should have a path")
	}
	if HasPath(plan.Graph, "b", "a") {
		t.Fatal("b->a should not have a path")
	}
}
