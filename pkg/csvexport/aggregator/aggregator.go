package aggregator

import (
	"github.com/RoboCup-SSL/ssl-match-stats/internal/referee"
	"github.com/RoboCup-SSL/ssl-match-stats/pkg/matchstats"
	"log"
	"math"
	"sort"
)

type Aggregator struct {
	Collection         *matchstats.MatchStatsCollection
	GamePhaseDurations map[string]*DurationStats
	GameEventDurations map[string]*DurationStats
	TeamStats          map[string]*matchstats.TeamStats
}

func NewAggregator(collection *matchstats.MatchStatsCollection) *Aggregator {
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
	durations := map[string][]int64{}

	for _, phaseName := range matchstats.GamePhaseType_name {
		gamePhaseDurations[phaseName] = new(DurationStats)
	}

	for _, phase := range matchStats.GamePhases {
		phaseName := (*phase).Type.String()
		gamePhaseDurations[phaseName].Duration += phase.Duration
		gamePhaseDurations[phaseName].Count += 1
		durations[phaseName] = append(durations[phaseName], phase.Duration)
	}

	for _, phaseName := range matchstats.GamePhaseType_name {
		stats := gamePhaseDurations[phaseName]
		phaseDurations := durations[phaseName]
		if len(phaseDurations) > 0 {
			sort.Slice(phaseDurations, func(i, j int) bool { return phaseDurations[i] < phaseDurations[j] })
			stats.DurationMin = phaseDurations[0]
			stats.DurationMax = phaseDurations[len(phaseDurations)-1]
			stats.DurationMedian = phaseDurations[len(phaseDurations)/2]
			stats.DurationAvg = int64(math.Round(float64(stats.Duration) / float64(len(phaseDurations))))
		}
	}

	checkSum := int64(0)
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
	durations := map[string][]int64{}

	for _, p := range referee.GameEvent_Type_name {
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
		durations[eventName] = append(durations[eventName], phase.Duration)
	}

	for _, eventName := range referee.GameEvent_Type_name {
		stats := gameEventDurations[eventName]
		eventDurations := durations[eventName]
		if len(eventDurations) > 0 {
			sort.Slice(eventDurations, func(i, j int) bool { return eventDurations[i] < eventDurations[j] })
			stats.DurationMin = eventDurations[0]
			stats.DurationMax = eventDurations[len(eventDurations)-1]
			stats.DurationMedian = eventDurations[len(eventDurations)/2]
			stats.DurationAvg = int64(math.Round(float64(stats.Duration) / float64(len(eventDurations))))
		}
	}

	return gameEventDurations
}
