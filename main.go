package main

import (
	"fmt"
	"os"

	"cpm-float/internal/cli"
)

const usage = `cpm-float — critical path method / PERT network planner

usage:
  cpm-float plan [flags] [file]
    compute time parameters, floats and the critical path

flags:
  -f, --file <path>    read network JSON from a file (default: stdin)
  --target <days>      also print PERT completion probability
  -h, --help           show this help

input JSON:
  { "activities": [ { "id", "duration", "predecessors": ["a"] } ] }
  or PERT: { "id", "a", "m", "b", "predecessors": [...] }

examples:
  cpm-float plan example/house.json
  cpm-float plan example/house.json --target 30
  cat example/house.json | cpm-float plan`

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(2)
	}
	cmd, args := os.Args[1], os.Args[2:]
	switch cmd {
	case "plan":
		os.Exit(cli.RunPlan(args, os.Stdout, os.Stderr))
	case "-h", "--help", "help":
		fmt.Println(usage)
		os.Exit(0)
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n\n%s\n", cmd, usage)
		os.Exit(2)
	}
}
