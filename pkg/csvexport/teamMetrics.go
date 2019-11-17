package csvexport

import (
	"github.com/RoboCup-SSL/ssl-match-stats/pkg/csvexport/aggregator"
	"github.com/RoboCup-SSL/ssl-match-stats/pkg/matchstats"
	"sort"
)

func WriteTeamMetricsPerGame(matchStatsCollection *matchstats.MatchStatsCollection, filename string) error {

	header := []string{"File", "Team", "Scored Goals", "Conceded Goals", "Fouls", "Yellow Cards", "Red Cards", "Timeout Time", "Timeouts", "Penalty Shots", "Ball Placement Time", "Ball Placements", "Max active Yellow Cards"}

	var records [][]string
	for _, matchStats := range matchStatsCollection.MatchStats {
		recordYellow := []string{matchStats.Name}
		recordYellow = append(recordYellow, teamNumbers(matchStats.TeamStatsYellow)...)
		records = append(records, recordYellow)
		recordBlue := []string{matchStats.Name}
		recordBlue = append(recordBlue, teamNumbers(matchStats.TeamStatsBlue)...)
		records = append(records, recordBlue)
	}

	return writeCsv(header, records, filename)
}

func WriteTeamMetricsSum(aggregator *aggregator.Aggregator, filename string) error {

	header := []string{"Team", "Scored Goals", "Conceded Goals", "Fouls", "Yellow Cards", "Red Cards", "Timeout Time", "Timeouts", "Penalty Shots", "Ball Placement Time", "Ball Placements", "Max active Yellow Cards"}

	var teamNamesSorted []string
	for k := range aggregator.TeamStats {
		teamNamesSorted = append(teamNamesSorted, k)
	}
	sort.Strings(teamNamesSorted)

	var records [][]string
	for _, teamName := range teamNamesSorted {
		teamStats := aggregator.TeamStats[teamName]
		records = append(records, teamNumbers(teamStats))
	}

	return writeCsv(header, records, filename)
}

func teamNumbers(stats *matchstats.TeamStats) []string {
	return []string{
		stats.Name,
		uintToStr(stats.Goals),
		uintToStr(stats.ConcededGoals),
		uintToStr(stats.Fouls),
		uintToStr(stats.YellowCards),
		uintToStr(stats.RedCards),
		uintToStr(stats.TimeoutTime),
		uintToStr(stats.TimeoutsTaken),
		uintToStr(stats.PenaltyShotsTotal),
		uintToStr(stats.BallPlacementTime),
		uintToStr(stats.BallPlacements),
		uintToStr(stats.MaxActiveYellowCards),
	}
}
