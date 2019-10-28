package csvexport

import (
	"github.com/RoboCup-SSL/ssl-match-stats/pkg/matchstats"
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
