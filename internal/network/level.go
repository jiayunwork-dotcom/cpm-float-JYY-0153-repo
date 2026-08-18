package network

import "sort"

type LevelInfo struct {
	Level int
	IDs   []string
}

func Levelize(g *Graph) []LevelInfo {
	level := make(map[string]int)
	for _, id := range g.Order {
		lv := 0
		for _, p := range g.Predecessors[id] {
			if level[p]+1 > lv {
				lv = level[p] + 1
			}
		}
		level[id] = lv
	}
	maxLevel := 0
	for _, lv := range level {
		if lv > maxLevel {
			maxLevel = lv
		}
	}
	out := []LevelInfo{}
	for i := 0; i <= maxLevel; i++ {
		ids := []string{}
		for id, lv := range level {
			if lv == i {
				ids = append(ids, id)
			}
		}
		sort.Strings(ids)
		out = append(out, LevelInfo{Level: i, IDs: ids})
	}
	return out
}

func LevelOf(g *Graph, id string) int {
	levels := Levelize(g)
	for _, lv := range levels {
		for _, nid := range lv.IDs {
			if nid == id {
				return lv.Level
			}
		}
	}
	return -1
}

func ActivitiesInLevel(g *Graph, level int) []string {
	for _, lv := range Levelize(g) {
		if lv.Level == level {
			return lv.IDs
		}
	}
	return nil
}

func MergeIntoLevel(g *Graph, levels []LevelInfo) []LevelInfo {
	merged := []LevelInfo{}
	for _, lv := range levels {
		if len(lv.IDs) == 0 {
			continue
		}
		merged = append(merged, lv)
	}
	return merged
}

func Breadth(g *Graph) int {
	max := 0
	for _, lv := range Levelize(g) {
		if len(lv.IDs) > max {
			max = len(lv.IDs)
		}
	}
	return max
}

func Depth(g *Graph) int {
	return len(Levelize(g)) - 1
}
