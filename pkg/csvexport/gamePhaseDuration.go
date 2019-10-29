package csvexport

import (
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
		for _, phaseName := range matchstats.GamePhaseType_name {
			record = append(record, uintToStr(matchStats.GamePhaseDurations[phaseName].Duration))
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

			record := []string{
				matchStats.Name,
				phaseName[6:],
				strconv.FormatBool(matchStats.ExtraTime),
				strconv.FormatBool(matchStats.Shootout),
				uintToStr(matchStats.GamePhaseDurations[phaseName].Duration),
				uintToStr(matchStats.GamePhaseDurations[phaseName].DurationMin),
				uintToStr(matchStats.GamePhaseDurations[phaseName].DurationMedian),
				uintToStr(uint32(math.Round(float64(matchStats.GamePhaseDurations[phaseName].DurationAvg)))),
				uintToStr(matchStats.GamePhaseDurations[phaseName].DurationMax),
				uintToStr(matchStats.GamePhaseDurations[phaseName].Count),
			}
			records = append(records, record)
		}
	}

	return writeCsv(header, records, filename)
}

func WriteGamePhaseDurationStatsAggregated(matchStatsCollection *matchstats.MatchStatsCollection, filename string) error {

	header := []string{"Phase", "Duration Sum", "Min Duration", "Median Duration", "Avg Duration", "Max Duration", "Count"}

	var records [][]string
	for _, phaseName := range matchstats.GamePhaseType_name {
		record := []string{
			phaseName,
			uintToStr(matchStatsCollection.GamePhaseDurations[phaseName].Duration),
			uintToStr(matchStatsCollection.GamePhaseDurations[phaseName].DurationMin),
			uintToStr(matchStatsCollection.GamePhaseDurations[phaseName].DurationMedian),
			uintToStr(uint32(math.Round(float64(matchStatsCollection.GamePhaseDurations[phaseName].DurationAvg)))),
			uintToStr(matchStatsCollection.GamePhaseDurations[phaseName].DurationMax),
			uintToStr(matchStatsCollection.GamePhaseDurations[phaseName].Count),
		}
		records = append(records, record)
	}

	return writeCsv(header, records, filename)
}
