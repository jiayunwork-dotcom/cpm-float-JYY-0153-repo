# cpm-float

Critical Path Method / PERT network planner. Given activity durations and
predecessor relations as JSON, it performs the forward and backward passes,
computes earliest/latest start and finish times, total and free float for every
activity, extracts the critical path (a zero-float chain from the virtual start
to the virtual end), and reports the project duration.

For PERT three-point estimates (a/m/b) it computes the mean `(a+4m+b)/6` and
variance `((b-a)/6)^2`, sums the variance along one selected critical path
(choosing the highest-variance one when several critical paths tie), and gives
the probability of finishing by a target date via a normal approximation
(standard normal CDF uses the Abramowitz & Stegun 7.1.26 rational approximation,
documented in `internal/pert/pert.go`).

## Build

```bash
go build .
go test ./...
```

## Usage

```bash
go run . plan example/house.json
go run . plan example/house.json --target 25
go run . plan example/pert.json --target 30
cat example/simple.json | go run . plan
```

Input JSON:

```json
{
  "activities": [
    { "id": "foundation", "duration": 4, "predecessors": ["site-prep"] },
    { "id": "site-prep", "duration": 2 }
  ]
}
```

PERT variant (no `duration`; `a/m/b` required):

```json
{ "id": "pour", "a": 3, "m": 4, "b": 8, "predecessors": ["excavate"] }
```

The graph auto-adds `START` and `END` dummy nodes: activities without
predecessors start from `START`, activities without successors feed `END`.

Output lists per-activity ES/EF/LS/LF/TF/FF with a `*` marker on critical
activities, then the project duration and the critical path. With `--target`,
the PERT completion probability `P(duration <= target)` is printed.

## Validation

Unknown predecessors, self loops, directed cycles, non-positive durations,
duplicate ids and invalid a/m/b triples (`a > m` or `m > b`) are rejected with
a message on stderr and a non-zero exit code. Total float satisfies both
`TF = LS - ES` and `TF = LF - EF`; free float never exceeds total float.
