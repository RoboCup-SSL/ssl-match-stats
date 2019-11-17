package csvexport

import (
	"github.com/RoboCup-SSL/ssl-match-stats/pkg/csvexport/aggregator"
	"github.com/RoboCup-SSL/ssl-match-stats/pkg/matchstats"
	"math"
	"strconv"
)

func WriteGamePhaseDurations(matchStatsCollection *matchstats.MatchStatsCollection, filename string) error {

	header := []string{"File", "Extra time", "Shootout"}
	for i := 0; i < len(matchstats.GamePhaseType_name); i++ {
		header = append(header, matchstats.GamePhaseType_name[int32(i)][6:])
	}

	var records [][]string
	for _, matchStats := range matchStatsCollection.MatchStats {
		record := []string{matchStats.Name, strconv.FormatBool(matchStats.ExtraTime), strconv.FormatBool(matchStats.Shootout)}
		durations := aggregator.AggregateGamePhaseDurations(matchStats)
		for _, phaseName := range matchstats.GamePhaseType_name {
			record = append(record, uintToStr(durations[phaseName].Duration))
		}
		records = append(records, record)
	}

	return writeCsv(header, records, filename)
}

func WriteGamePhaseDurationStats(matchStatsCollection *matchstats.MatchStatsCollection, filename string) error {

	header := []string{"File", "Phase", "Extra time", "Shootout", "Duration Sum", "Min Duration", "Median Duration", "Avg Duration", "Max Duration", "Count"}

	var records [][]string
	for _, matchStats := range matchStatsCollection.MatchStats {
		for i := 0; i < len(matchstats.GamePhaseType_name); i++ {
			phaseName := matchstats.GamePhaseType_name[int32(i)]
			durations := aggregator.AggregateGamePhaseDurations(matchStats)

			record := []string{
				matchStats.Name,
				phaseName[6:],
				strconv.FormatBool(matchStats.ExtraTime),
				strconv.FormatBool(matchStats.Shootout),
				uintToStr(durations[phaseName].Duration),
				uintToStr(durations[phaseName].DurationMin),
				uintToStr(durations[phaseName].DurationMedian),
				uintToStr(uint32(math.Round(float64(durations[phaseName].DurationAvg)))),
				uintToStr(durations[phaseName].DurationMax),
				uintToStr(durations[phaseName].Count),
			}
			records = append(records, record)
		}
	}

	return writeCsv(header, records, filename)
}

func WriteGamePhaseDurationStatsAggregated(aggregator *aggregator.Aggregator, filename string) error {

	header := []string{"Phase", "Duration Sum", "Min Duration", "Median Duration", "Avg Duration", "Max Duration", "Count"}

	var records [][]string
	for _, phaseName := range matchstats.GamePhaseType_name {
		record := []string{
			phaseName,
			uintToStr(aggregator.GamePhaseDurations[phaseName].Duration),
			uintToStr(aggregator.GamePhaseDurations[phaseName].DurationMin),
			uintToStr(aggregator.GamePhaseDurations[phaseName].DurationMedian),
			uintToStr(uint32(math.Round(float64(aggregator.GamePhaseDurations[phaseName].DurationAvg)))),
			uintToStr(aggregator.GamePhaseDurations[phaseName].DurationMax),
			uintToStr(aggregator.GamePhaseDurations[phaseName].Count),
		}
		records = append(records, record)
	}

	return writeCsv(header, records, filename)
}
