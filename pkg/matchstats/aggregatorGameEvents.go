package matchstats

import (
	"github.com/RoboCup-SSL/ssl-go-tools/pkg/sslproto"
	"math"
	"sort"
)

func (a *Aggregator) AggregateGameEvents() error {

	a.Collection.GameEventDurations = map[string]*DurationStats{}
	durations := map[string][]int{}

	for _, name := range sslproto.GameEventType_name {
		a.Collection.GameEventDurations[name] = new(DurationStats)
	}

	for _, matchStats := range a.Collection.MatchStats {
		for _, phase := range matchStats.GamePhases {
			if len(phase.GameEventsEntry) == 0 {
				continue
			}

			primaryEvent := phase.GameEventsEntry[0]
			eventName := primaryEvent.Type.String()
			a.Collection.GameEventDurations[eventName].Duration += phase.Duration
			durations[eventName] = append(durations[eventName], int(phase.Duration))
		}
	}

	for _, eventName := range sslproto.GameEventType_name {
		stats := a.Collection.GameEventDurations[eventName]
		eventDurations := durations[eventName]
		if len(eventDurations) > 0 {
			sort.Ints(eventDurations)
			stats.DurationMin = uint32(eventDurations[0])
			stats.DurationMax = uint32(eventDurations[len(eventDurations)-1])
			stats.DurationMedian = uint32(eventDurations[len(eventDurations)/2])
			stats.DurationAvg = uint32(math.Round(float64(stats.Duration) / float64(len(eventDurations))))
		}
	}

	return nil
}
