package pert

import (
	"errors"
	"math"

	"cpm-float/internal/network"
)

var (
	ErrBadTriple = errors.New("pert: a/m/b triple invalid")
	ErrNoCriticalPath = errors.New("pert: no critical path available")
)

type Triple struct {
	A float64
	M float64
	B float64
}

type Estimate struct {
	Mean float64
	Variance float64
}

func FromTriple(t Triple) (Estimate, error) {
	if t.A > t.M || t.M > t.B {
		return Estimate{}, ErrBadTriple
	}
	if t.A < 0 {
		return Estimate{}, ErrBadTriple
	}
	mean := (t.A + 4*t.M + t.B) / 6
	span := (t.B - t.A) / 6
	return Estimate{Mean: mean, Variance: span * span}, nil
}

func DurationFromTriple(t Triple) float64 {
	e, err := FromTriple(t)
	if err != nil {
		return 0
	}
	return e.Mean
}

func VarianceFromTriple(t Triple) float64 {
	e, err := FromTriple(t)
	if err != nil {
		return 0
	}
	return e.Variance
}

type PathStat struct {
	Path     []string
	Mean     float64
	Variance float64
}

func CriticalPathStats(net *network.Graph, triples map[string]Triple) (PathStat, error) {
	paths := allPaths(net, net.Source, net.Sink)
	if len(paths) == 0 {
		return PathStat{}, ErrNoCriticalPath
	}
	critical := []PathStat{}
	for _, p := range paths {
		mean, variance := statsOf(p, net, triples)
		critical = append(critical, PathStat{Path: p, Mean: mean, Variance: variance})
	}
	best := critical[0]
	bestMean := -1.0
	for _, c := range critical {
		if c.Mean > bestMean+1e-9 {
			best = c
			bestMean = c.Mean
		} else if math.Abs(c.Mean-bestMean) <= 1e-9 && c.Variance > best.Variance {
			best = c
		}
	}
	return best, nil
}

func statsOf(path []string, net *network.Graph, triples map[string]Triple) (float64, float64) {
	mean := 0.0
	variance := 0.0
	for _, id := range path {
		if id == net.Source || id == net.Sink {
			continue
		}
		t, ok := triples[id]
		if !ok {
			continue
		}
		e, err := FromTriple(t)
		if err != nil {
			continue
		}
		mean += e.Mean
		variance += e.Variance
	}
	return mean, variance
}

func allPaths(net *network.Graph, from, to string) [][]string {
	paths := [][]string{}
	seen := make(map[string]bool)
	var walk func(id string, cur []string)
	walk = func(id string, cur []string) {
		if id == to {
			path := append([]string(nil), cur...)
			path = append(path, id)
			paths = append(paths, path)
			return
		}
		if seen[id] {
			return
		}
		seen[id] = true
		next := append([]string(nil), cur...)
		next = append(next, id)
		for _, s := range net.Successors[id] {
			walk(s, next)
		}
		seen[id] = false
	}
	walk(from, []string{})
	return paths
}

func Phi(x float64) float64 {
	// Abramowitz & Stegun 7.1.26 rational approximation
	sign := 1.0
	if x < 0 {
		sign = -1.0
		x = -x
	}
	t := 1 / (1 + 0.2316419*x)
	a1 := 0.319381530
	a2 := -0.356563782
	a3 := 1.781477937
	a4 := -1.821255978
	a5 := 1.330274429
	poly := t * (a1 + t*(a2+t*(a3+t*(a4+t*a5))))
	phi := 1 - 0.3989422804014327*math.Exp(-x*x/2)*poly
	if sign < 0 {
		return 1 - phi
	}
	return phi
}

func CompletionProbability(mean, stddev, target float64) float64 {
	if stddev <= 0 {
		if target >= mean {
			return 1
		}
		return 0
	}
	z := (target - mean) / stddev
	return Phi(z)
}

func InversePhi(p float64) float64 {
	if p <= 0 {
		return math.Inf(-1)
	}
	if p >= 1 {
		return math.Inf(1)
	}
	// rational approximation (Acklam's algorithm)
	a := []float64{-3.969683028665376e+01, 2.209460984245205e+02,
		-2.759285104469687e+02, 1.383577518672690e+02,
		-3.066479806614716e+01, 2.506628277459239e+00}
	b := []float64{-5.447609879822406e+01, 1.615858368580409e+02,
		-1.556989798598866e+02, 6.680131188771972e+01,
		-1.328068155288572e+01}
	c := []float64{-7.784894002430293e-03, -3.223964580411365e-01,
		-2.400758277161838e+00, -2.549732539343734e+00,
		4.374664141464968e+00, 2.938163982698783e+00}
	d := []float64{7.784695709041462e-03, 3.224671290700398e-01,
		2.445134137142996e+00, 3.754408661907416e+00}
	plow := 0.02425
	phigh := 1 - plow
	var q, r float64
	var x float64
	if p < plow {
		q = math.Sqrt(-2 * math.Log(p))
		x = (((((c[0]*q+c[1])*q+c[2])*q+c[3])*q+c[4])*q + c[5]) /
			((((d[0]*q+d[1])*q+d[2])*q+d[3])*q + 1)
	} else if p <= phigh {
		q = p - 0.5
		r = q * q
		x = (((((a[0]*r+a[1])*r+a[2])*r+a[3])*r+a[4])*r + a[5]) * q /
			(((((b[0]*r+b[1])*r+b[2])*r+b[3])*r+b[4])*r + 1)
	} else {
		q = math.Sqrt(-2 * math.Log(1-p))
		x = -(((((c[0]*q+c[1])*q+c[2])*q+c[3])*q+c[4])*q + c[5]) /
			((((d[0]*q+d[1])*q+d[2])*q+d[3])*q + 1)
	}
	return x
}

func TargetForProbability(mean, stddev, probability float64) float64 {
	z := InversePhi(probability)
	return mean + z*stddev
}
