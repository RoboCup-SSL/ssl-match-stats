package csvexport

import (
	"github.com/RoboCup-SSL/ssl-match-stats/internal/referee"
	"github.com/RoboCup-SSL/ssl-match-stats/pkg/csvexport/aggregator"
	"github.com/RoboCup-SSL/ssl-match-stats/pkg/matchstats"
	"math"
)

func WriteGameEventDurationStats(matchStatsCollection *matchstats.MatchStatsCollection, filename string) error {

	header := []string{"File", "Game Event", "Duration Sum", "Min Duration", "Median Duration", "Avg Duration", "Max Duration", "Count"}

	var records [][]string
	for _, matchStats := range matchStatsCollection.MatchStats {
		for _, eventName := range referee.GameEvent_Type_name {
			durations := aggregator.AggregateGameEventDurations(matchStats)
			record := []string{
				matchStats.Name,
				eventName,
				int64ToStr(durations[eventName].Duration),
				int64ToStr(durations[eventName].DurationMin),
				int64ToStr(durations[eventName].DurationMedian),
				int64ToStr(int64(math.Round(float64(durations[eventName].DurationAvg)))),
				int64ToStr(durations[eventName].DurationMax),
				int64ToStr(durations[eventName].Count),
			}
			records = append(records, record)
		}
	}

	return writeCsv(header, records, filename)
}

func WriteGameEventDurationStatsAggregated(aggregator *aggregator.Aggregator, filename string) error {

	header := []string{"Game Event", "Duration Sum", "Min Duration", "Median Duration", "Avg Duration", "Max Duration", "Count"}

	var records [][]string
	for _, eventName := range referee.GameEvent_Type_name {
		record := []string{
			eventName,
			int64ToStr(aggregator.GameEventDurations[eventName].Duration),
			int64ToStr(aggregator.GameEventDurations[eventName].DurationMin),
			int64ToStr(aggregator.GameEventDurations[eventName].DurationMedian),
			int64ToStr(int64(math.Round(float64(aggregator.GameEventDurations[eventName].DurationAvg)))),
			int64ToStr(aggregator.GameEventDurations[eventName].DurationMax),
			int64ToStr(aggregator.GameEventDurations[eventName].Count),
		}
		records = append(records, record)
	}

	return writeCsv(header, records, filename)
}
