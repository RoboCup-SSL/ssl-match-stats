package csvexport

import (
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

func WriteTeamMetricsSum(matchStatsCollection *matchstats.MatchStatsCollection, filename string) error {

	header := []string{"Team", "Scored Goals", "Conceded Goals", "Fouls", "Yellow Cards", "Red Cards", "Timeout Time", "Timeouts", "Penalty Shots", "Ball Placement Time", "Ball Placements", "Max active Yellow Cards"}

	teams := map[string]*matchstats.TeamStats{}
	for _, matchStats := range matchStatsCollection.MatchStats {
		teams[matchStats.TeamStatsYellow.Name] = &matchstats.TeamStats{Name: matchStats.TeamStatsYellow.Name}
		teams[matchStats.TeamStatsBlue.Name] = &matchstats.TeamStats{Name: matchStats.TeamStatsBlue.Name}
	}

	for _, matchStats := range matchStatsCollection.MatchStats {
		addTeamStats(teams[matchStats.TeamStatsYellow.Name], matchStats.TeamStatsYellow)
		addTeamStats(teams[matchStats.TeamStatsBlue.Name], matchStats.TeamStatsBlue)
	}

	var teamNamesSorted []string
	for k := range teams {
		teamNamesSorted = append(teamNamesSorted, k)
	}
	sort.Strings(teamNamesSorted)

	var records [][]string
	for _, teamName := range teamNamesSorted {
		teamStats := teams[teamName]
		records = append(records, teamNumbers(teamStats))
	}

	return writeCsv(header, records, filename)
}

func addTeamStats(to *matchstats.TeamStats, team *matchstats.TeamStats) {
	to.Goals += team.Goals
	to.ConcededGoals += team.ConcededGoals
	to.Fouls += team.Fouls
	to.YellowCards += team.YellowCards
	to.RedCards += team.RedCards
	to.TimeoutTime += team.TimeoutTime
	to.Timeouts += team.Timeouts
	to.PenaltyShotsTotal += team.PenaltyShotsTotal
	to.BallPlacementTime += team.BallPlacementTime
	to.BallPlacements += team.BallPlacements
	if to.MaxActiveYellowCards < team.MaxActiveYellowCards {
		to.MaxActiveYellowCards = team.MaxActiveYellowCards
	}
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
		uintToStr(stats.Timeouts),
		uintToStr(stats.PenaltyShotsTotal),
		uintToStr(stats.BallPlacementTime),
		uintToStr(stats.BallPlacements),
		uintToStr(stats.MaxActiveYellowCards),
	}
}
