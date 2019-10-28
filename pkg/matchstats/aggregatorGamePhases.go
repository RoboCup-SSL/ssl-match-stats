package matchstats

import (
	"math"
	"sort"
)

func (a *Aggregator) AggregateGamePhases() error {

	a.Collection.GamePhaseDurations = map[string]*DurationStats{}
	durations := map[string][]int{}

	for _, phaseName := range GamePhaseType_name {
		a.Collection.GamePhaseDurations[phaseName] = new(DurationStats)
	}

	for _, matchStats := range a.Collection.MatchStats {
		for _, phase := range matchStats.GamePhases {
			phaseName := (*phase).Type.String()
			a.Collection.GamePhaseDurations[phaseName].Duration += phase.Duration
			durations[phaseName] = append(durations[phaseName], int(phase.Duration))
		}

	}

	for _, phaseName := range GamePhaseType_name {
		stats := a.Collection.GamePhaseDurations[phaseName]
		phaseDurations := durations[phaseName]
		if len(phaseDurations) > 0 {
			sort.Ints(phaseDurations)
			stats.DurationMin = uint32(phaseDurations[0])
			stats.DurationMax = uint32(phaseDurations[len(phaseDurations)-1])
			stats.DurationMedian = uint32(phaseDurations[len(phaseDurations)/2])
			stats.DurationAvg = uint32(math.Round(float64(stats.Duration) / float64(len(phaseDurations))))
		}
	}

	return nil
}
