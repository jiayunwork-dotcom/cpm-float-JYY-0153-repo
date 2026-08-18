package network

import (
	"math"
	"sort"
)

type CrashActivity struct {
	ID       string
	Normal   float64
	Crash    float64
	CostSlope float64
}

type CrashResult struct {
	Activity    string
	Duration    float64
	ProjectTime float64
	Cost        float64
	Applied     bool
}

func CrashProject(in Input, crashes map[string]CrashActivity, target float64) ([]CrashResult, error) {
	if err := Validate(in); err != nil {
		return nil, err
	}
	current := in
	results := []CrashResult{}
	iter := 0
	for iter < 1000 {
		plan, err := Plan(current)
		if err != nil {
			return nil, err
		}
		time := plan.MaxDuration()
		if time <= target+1e-9 {
			break
		}
		best := ""
		bestSlope := math.Inf(1)
		for id, c := range crashes {
			act := activityByID(plan.Graph, id)
			if act.Duration <= c.Crash+1e-9 {
				continue
			}
			if c.CostSlope < bestSlope {
				bestSlope = c.CostSlope
				best = id
			}
		}
		if best == "" {
			break
		}
		plan2, _ := Plan(current)
		if plan2 == nil {
			break
		}
		c := crashes[best]
		step := (activityByID(plan2.Graph, best).Duration - c.Crash) / 2
		if step < 0.5 {
			step = 0.5
		}
		act := activityByID(plan2.Graph, best)
		newDur := act.Duration - step
		if newDur < c.Crash {
			newDur = c.Crash
		}
		for i := range current.Activities {
			if current.Activities[i].ID == best {
				current.Activities[i].Duration = newDur
			}
		}
		results = append(results, CrashResult{
			Activity:    best,
			Duration:    newDur,
			ProjectTime: time,
			Cost:        bestSlope * step,
			Applied:     true,
		})
		iter++
	}
	final, _ := Plan(current)
	if final == nil {
		return results, nil
	}
	return results, nil
}

func CrashOptions(in Input, crashDays map[string]float64) []CrashActivity {
	out := []CrashActivity{}
	for _, a := range in.Activities {
		if crash, ok := crashDays[a.ID]; ok && crash < a.Duration {
			out = append(out, CrashActivity{
				ID:        a.ID,
				Normal:    a.Duration,
				Crash:     crash,
				CostSlope: (a.Duration - crash) * 10,
			})
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].ID < out[j].ID
	})
	return out
}

func activityByID2(acts []Activity, id string) Activity {
	for _, a := range acts {
		if a.ID == id {
			return a
		}
	}
	return Activity{ID: id}
}
