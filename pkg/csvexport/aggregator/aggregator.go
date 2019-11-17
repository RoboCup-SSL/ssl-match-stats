package aggregator

import (
	"github.com/RoboCup-SSL/ssl-go-tools/pkg/sslproto"
	"github.com/RoboCup-SSL/ssl-match-stats/pkg/matchstats"
	"log"
	"math"
	"sort"
)

type Aggregator struct {
	Collection         matchstats.MatchStatsCollection
	GamePhaseDurations map[string]*DurationStats
	GameEventDurations map[string]*DurationStats
	TeamStats          map[string]*matchstats.TeamStats
}

func NewAggregator(collection matchstats.MatchStatsCollection) *Aggregator {
	generator := new(Aggregator)
	generator.Collection = collection
	return generator
}

func (a *Aggregator) Aggregate() error {
	if err := a.AggregateGamePhases(); err != nil {
		return err
	}
	if err := a.AggregateGameEvents(); err != nil {
		return err
	}
	if err := a.AggregateTeamMetrics(); err != nil {
		return err
	}

	return nil
}

func AggregateGamePhaseDurations(matchStats *matchstats.MatchStats) map[string]*DurationStats {

	gamePhaseDurations := map[string]*DurationStats{}
	durations := map[string][]int{}

	for _, phaseName := range matchstats.GamePhaseType_name {
		gamePhaseDurations[phaseName] = new(DurationStats)
	}

	for _, phase := range matchStats.GamePhases {
		phaseName := (*phase).Type.String()
		gamePhaseDurations[phaseName].Duration += phase.Duration
		gamePhaseDurations[phaseName].Count += 1
		durations[phaseName] = append(durations[phaseName], int(phase.Duration))
	}

	for _, phaseName := range matchstats.GamePhaseType_name {
		stats := gamePhaseDurations[phaseName]
		phaseDurations := durations[phaseName]
		if len(phaseDurations) > 0 {
			sort.Ints(phaseDurations)
			stats.DurationMin = uint32(phaseDurations[0])
			stats.DurationMax = uint32(phaseDurations[len(phaseDurations)-1])
			stats.DurationMedian = uint32(phaseDurations[len(phaseDurations)/2])
			stats.DurationAvg = uint32(math.Round(float64(stats.Duration) / float64(len(phaseDurations))))
		}
	}

	checkSum := uint32(0)
	for _, phaseName := range matchstats.GamePhaseType_name {
		checkSum += gamePhaseDurations[phaseName].Duration
		gamePhaseDurations[phaseName].DurationRelative = float32(gamePhaseDurations[phaseName].Duration) / float32(matchStats.MatchDuration)
	}

	if matchStats.MatchDuration != checkSum {
		log.Printf("Match duration mismatch. Total: %v, Sum of phases: %v, Diff: %v", matchStats.MatchDuration, checkSum, matchStats.MatchDuration-checkSum)
	}

	return gamePhaseDurations
}

func AggregateGameEventDurations(matchStats *matchstats.MatchStats) map[string]*DurationStats {

	gameEventDurations := map[string]*DurationStats{}
	durations := map[string][]int{}

	for _, p := range sslproto.GameEventType_name {
		gameEventDurations[p] = new(DurationStats)
		gameEventDurations[p].Duration = 0
	}

	for _, phase := range matchStats.GamePhases {
		if len(phase.GameEventsEntry) == 0 {
			continue
		}

		primaryEvent := phase.GameEventsEntry[0]
		eventName := primaryEvent.Type.String()
		gameEventDurations[eventName].Duration += phase.Duration
		gameEventDurations[eventName].Count += 1
		durations[eventName] = append(durations[eventName], int(phase.Duration))
	}

	for _, eventName := range sslproto.GameEventType_name {
		stats := gameEventDurations[eventName]
		eventDurations := durations[eventName]
		if len(eventDurations) > 0 {
			sort.Ints(eventDurations)
			stats.DurationMin = uint32(eventDurations[0])
			stats.DurationMax = uint32(eventDurations[len(eventDurations)-1])
			stats.DurationMedian = uint32(eventDurations[len(eventDurations)/2])
			stats.DurationAvg = uint32(math.Round(float64(stats.Duration) / float64(len(eventDurations))))
		}
	}

	return gameEventDurations
}
