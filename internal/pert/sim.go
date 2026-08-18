package pert

import (
	"math"
	"math/rand"
)

type SimResult struct {
	Samples    []float64
	Mean       float64
	StdDev     float64
	P25        float64
	P50        float64
	P75        float64
}

func SampleTriple(t Triple, rng *rand.Rand) float64 {
	e, err := FromTriple(t)
	if err != nil {
		return 0
	}
	return e.Mean + math.Sqrt(e.Variance)*rng.NormFloat64()
}

func SimulateDuration(p *Project, n int, seed int64) SimResult {
	rng := rand.New(rand.NewSource(seed))
	samples := make([]float64, n)
	mean := 0.0
	for i := 0; i < n; i++ {
		total := 0.0
		for _, id := range p.Critical.Path {
			t, ok := p.Triples[id]
			if !ok {
				continue
			}
			total += SampleTriple(t, rng)
		}
		samples[i] = total
		mean += total
	}
	mean /= float64(n)
	variance := 0.0
	for _, s := range samples {
		d := s - mean
		variance += d * d
	}
	variance /= float64(n)
	return SimResult{
		Samples: samples,
		Mean:    mean,
		StdDev:  math.Sqrt(variance),
		P25:     percentile(samples, 0.25),
		P50:     percentile(samples, 0.50),
		P75:     percentile(samples, 0.75),
	}
}

func percentile(sorted []float64, p float64) float64 {
	if len(sorted) == 0 {
		return 0
	}
	cp := append([]float64(nil), sorted...)
	sortFloats(cp)
	idx := int(p * float64(len(cp)-1))
	return cp[idx]
}

func sortFloats(a []float64) {
	for i := 1; i < len(a); i++ {
		for j := i; j > 0 && a[j] < a[j-1]; j-- {
			a[j], a[j-1] = a[j-1], a[j]
		}
	}
}

func SimulatedProbability(s SimResult, target float64) float64 {
	if len(s.Samples) == 0 {
		return 0
	}
	count := 0
	for _, v := range s.Samples {
		if v <= target {
			count++
		}
	}
	return float64(count) / float64(len(s.Samples))
}

func ConfidenceInterval(s SimResult) (float64, float64) {
	lo := s.Mean - 1.96*s.StdDev
	hi := s.Mean + 1.96*s.StdDev
	return lo, hi
}
