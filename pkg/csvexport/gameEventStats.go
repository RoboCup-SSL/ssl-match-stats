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
				uintToStr(durations[eventName].Duration),
				uintToStr(durations[eventName].DurationMin),
				uintToStr(durations[eventName].DurationMedian),
				uintToStr(uint32(math.Round(float64(durations[eventName].DurationAvg)))),
				uintToStr(durations[eventName].DurationMax),
				uintToStr(durations[eventName].Count),
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
			uintToStr(aggregator.GameEventDurations[eventName].Duration),
			uintToStr(aggregator.GameEventDurations[eventName].DurationMin),
			uintToStr(aggregator.GameEventDurations[eventName].DurationMedian),
			uintToStr(uint32(math.Round(float64(aggregator.GameEventDurations[eventName].DurationAvg)))),
			uintToStr(aggregator.GameEventDurations[eventName].DurationMax),
			uintToStr(aggregator.GameEventDurations[eventName].Count),
		}
		records = append(records, record)
	}

	return writeCsv(header, records, filename)
}
