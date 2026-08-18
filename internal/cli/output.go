package cli

import (
	"fmt"
	"io"

	"cpm-float/internal/network"
	"cpm-float/internal/pert"
)

func writePlan(w io.Writer, plan *network.PlanResult) {
	fmt.Fprint(w, plan.FormatTable())
	fmt.Fprintf(w, "project_duration %.2f\n", plan.MaxDuration())
	fmt.Fprintf(w, "critical_path    %s\n", joinPath(plan.PathActivities))
}

func writeMetrics(w io.Writer, m network.Metrics) {
	fmt.Fprintf(w, "activities   %d\n", m.ActivityCount)
	fmt.Fprintf(w, "edges        %d\n", m.EdgeCount)
	fmt.Fprintf(w, "critical     %d (%.2f%%)\n", m.CriticalCount, m.CriticalRatio*100)
	fmt.Fprintf(w, "max_parallel %d\n", m.MaxParallel)
	fmt.Fprintf(w, "longest      %.2f\n", m.LongestPath)
}

func writePert(w io.Writer, proj *pert.Project) {
	fmt.Fprint(w, pert.FormatStats(proj))
	fmt.Fprintf(w, "critical_path  %s\n", pert.FormatCriticalPath(proj))
}

func writeRisk(w io.Writer, proj *pert.Project) {
	report := pert.RiskReportFor(proj)
	for _, a := range report.Activities {
		fmt.Fprintf(w, "risk %s mean=%.2f var=%.4f sigma=%.2f crit=%v\n",
			a.ID, a.Mean, a.Variance, a.StdDev, a.Critical)
	}
}
