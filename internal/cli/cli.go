package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"cpm-float/internal/network"
	"cpm-float/internal/pert"
)

var ErrEmptyInput = errors.New("cli: empty input")
var ErrInvalidJSON = errors.New("cli: invalid json")

var stdinOverride io.Reader

func ReadInput(r io.Reader) (network.Input, error) {
	var in network.Input
	dec := json.NewDecoder(r)
	if err := dec.Decode(&in); err != nil {
		if err == io.EOF {
			return in, ErrEmptyInput
		}
		return in, fmt.Errorf("%w: %v", ErrInvalidJSON, err)
	}
	return in, nil
}

func ReadFile(path string) (network.Input, error) {
	f, err := os.Open(path)
	if err != nil {
		return network.Input{}, err
	}
	defer f.Close()
	return ReadInput(f)
}

func ReadStdin() (network.Input, error) {
	if stdinOverride != nil {
		return ReadInput(stdinOverride)
	}
	return ReadInput(os.Stdin)
}

func RunPlan(args []string, stdout, stderr io.Writer) int {
	var path string
	target := 0.0
	hasTarget := false
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-f", "--file":
			if i+1 < len(args) {
				i++
				path = args[i]
			}
		case "--target":
			if i+1 < len(args) {
				i++
				_, err := fmt.Sscanf(args[i], "%f", &target)
				if err == nil {
					hasTarget = true
				}
			}
		case "-h", "--help":
			fmt.Fprintln(stdout, planUsage)
			return 0
		default:
			if path == "" {
				path = args[i]
			}
		}
	}

	var in network.Input
	var err error
	if path != "" {
		in, err = ReadFile(path)
	} else {
		in, err = ReadStdin()
	}
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}

	triples, err := pert.ExtractTriples(in)
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}

	if len(triples) > 0 {
		in = pert.FillDurations(in, triples)
		proj, err := pert.BuildProject(in, triples)
		if err != nil {
			fmt.Fprintf(stderr, "error: %v\n", err)
			return 1
		}
		fmt.Fprintln(stdout, pert.FormatStats(proj))
		if hasTarget {
			fmt.Fprintln(stdout, pert.FormatProbability(proj, target))
		}
		return 0
	}

	plan, err := network.Plan(in)
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}
	fmt.Fprint(stdout, plan.FormatTable())
	fmt.Fprintf(stdout, "project_duration %.2f\n", plan.MaxDuration())
	fmt.Fprintf(stdout, "critical_path    %s\n", joinPath(plan.PathActivities))
	if hasTarget {
		fmt.Fprintf(stdout, "completion       n/a (deterministic)\n")
	}
	return 0
}

func joinPath(ids []string) string {
	out := ""
	for i, id := range ids {
		if i > 0 {
			out += " -> "
		}
		out += id
	}
	return out
}

const planUsage = `usage: cpm-float plan [-f <file>] [--target <days>]

compute forward/backward pass, floats and the critical path of a CPM network
from JSON on stdin or from a file. activities carry id, duration (or a/m/b),
and predecessors. --target adds a PERT completion probability.`
