package cli

import (
	"encoding/json"
	"fmt"
	"io"

	"cpm-float/internal/network"
)

func FormatJSON(v interface{}) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	return string(b)
}

func PlanToJSON(plan *network.PlanResult) map[string]interface{} {
	nodes := []map[string]interface{}{}
	for _, ref := range plan.Activities {
		n := ref.Node
		nodes = append(nodes, map[string]interface{}{
			"id":       ref.ID,
			"es":       n.ES,
			"ef":       n.EF,
			"ls":       n.LS,
			"lf":       n.LF,
			"tf":       n.TF,
			"ff":       ref.Free,
			"critical": n.Critical,
		})
	}
	return map[string]interface{}{
		"duration":       plan.MaxDuration(),
		"critical_path":  plan.PathActivities,
		"activities":     nodes,
	}
}

func WritePlanJSON(w io.Writer, plan *network.PlanResult) {
	fmt.Fprintln(w, FormatJSON(PlanToJSON(plan)))
}

func WriteMetricsJSON(w io.Writer, m network.Metrics) {
	fmt.Fprintln(w, FormatJSON(map[string]interface{}{
		"activities":   m.ActivityCount,
		"edges":        m.EdgeCount,
		"critical":     m.CriticalCount,
		"critical_pct": m.CriticalRatio,
		"max_parallel": m.MaxParallel,
		"longest":      m.LongestPath,
	}))
}

func ReadPathOnly(path string) (network.Input, error) {
	return ReadFile(path)
}

func writeError(w io.Writer, err error) int {
	fmt.Fprintf(w, "error: %v\n", err)
	return 1
}
