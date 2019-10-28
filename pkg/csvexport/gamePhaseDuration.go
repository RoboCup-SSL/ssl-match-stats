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
		for i := 0; i < len(matchstats.GamePhaseType_name); i++ {
			name := matchstats.GamePhaseType_name[int32(i)]
			record = append(record, uintToStr(matchStats.GamePhaseStats[name].Duration))
		}
		records = append(records, record)
	}

	return writeCsv(header, records, filename)
}

func WriteGamePhaseDurationStats(matchStatsCollection *matchstats.MatchStatsCollection, filename string) error {

	header := []string{"File", "Phase", "Extra time", "Shootout", "Duration Sum", "Min Duration", "Median Duration", "Avg Duration", "Max Duration"}

	var records [][]string
	for _, matchStats := range matchStatsCollection.MatchStats {
		for i := 0; i < len(matchstats.GamePhaseType_name); i++ {
			phaseName := matchstats.GamePhaseType_name[int32(i)]

			record := []string{
				matchStats.Name,
				phaseName[6:],
				strconv.FormatBool(matchStats.ExtraTime),
				strconv.FormatBool(matchStats.Shootout),
				uintToStr(matchStats.GamePhaseStats[phaseName].Duration),
				uintToStr(matchStats.GamePhaseStats[phaseName].DurationMin),
				uintToStr(matchStats.GamePhaseStats[phaseName].DurationMedian),
				uintToStr(uint32(math.Round(float64(matchStats.GamePhaseStats[phaseName].DurationAvg)))),
				uintToStr(matchStats.GamePhaseStats[phaseName].DurationMax),
			}
			records = append(records, record)
		}
	}

	return writeCsv(header, records, filename)
}
