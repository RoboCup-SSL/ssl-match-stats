package csvexport

import (
	"github.com/RoboCup-SSL/ssl-go-tools/pkg/sslproto"
	"github.com/RoboCup-SSL/ssl-match-stats/pkg/matchstats"
	"math"
)

func WriteGameEventDurationStats(matchStatsCollection *matchstats.MatchStatsCollection, filename string) error {

	header := []string{"File", "Game Event", "Duration Sum", "Min Duration", "Median Duration", "Avg Duration", "Max Duration"}

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
			}
			records = append(records, record)
		}
	}

	return writeCsv(header, records, filename)
}
