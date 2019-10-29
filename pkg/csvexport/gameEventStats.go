package csvexport

import (
	"github.com/RoboCup-SSL/ssl-go-tools/pkg/sslproto"
	"github.com/RoboCup-SSL/ssl-match-stats/pkg/matchstats"
	"math"
)

func WriteGameEventDurationStats(matchStatsCollection *matchstats.MatchStatsCollection, filename string) error {

	header := []string{"File", "Game Event", "Duration Sum", "Min Duration", "Median Duration", "Avg Duration", "Max Duration", "Count"}

	var records [][]string
	for _, matchStats := range matchStatsCollection.MatchStats {
		for _, eventName := range sslproto.GameEventType_name {
			record := []string{
				matchStats.Name,
				eventName,
				uintToStr(matchStats.GameEventDurations[eventName].Duration),
				uintToStr(matchStats.GameEventDurations[eventName].DurationMin),
				uintToStr(matchStats.GameEventDurations[eventName].DurationMedian),
				uintToStr(uint32(math.Round(float64(matchStats.GameEventDurations[eventName].DurationAvg)))),
				uintToStr(matchStats.GameEventDurations[eventName].DurationMax),
				uintToStr(matchStats.GameEventDurations[eventName].Count),
			}
			records = append(records, record)
		}
	}

	return writeCsv(header, records, filename)
}

func WriteGameEventDurationStatsAggregated(matchStatsCollection *matchstats.MatchStatsCollection, filename string) error {

	header := []string{"Game Event", "Duration Sum", "Min Duration", "Median Duration", "Avg Duration", "Max Duration", "Count"}

	var records [][]string
	for _, eventName := range sslproto.GameEventType_name {
		record := []string{
			eventName,
			uintToStr(matchStatsCollection.GameEventDurations[eventName].Duration),
			uintToStr(matchStatsCollection.GameEventDurations[eventName].DurationMin),
			uintToStr(matchStatsCollection.GameEventDurations[eventName].DurationMedian),
			uintToStr(uint32(math.Round(float64(matchStatsCollection.GameEventDurations[eventName].DurationAvg)))),
			uintToStr(matchStatsCollection.GameEventDurations[eventName].DurationMax),
			uintToStr(matchStatsCollection.GameEventDurations[eventName].Count),
		}
		records = append(records, record)
	}

	return writeCsv(header, records, filename)
}
