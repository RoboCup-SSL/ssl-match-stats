package aggregator

import "github.com/RoboCup-SSL/ssl-match-stats/pkg/matchstats"

func (a *Aggregator) AggregateTeamMetrics() error {
	a.TeamStats = map[string]*matchstats.TeamStats{}
	for _, matchStats := range a.Collection.MatchStats {
		a.TeamStats[matchStats.TeamStatsYellow.Name] = &matchstats.TeamStats{Name: matchStats.TeamStatsYellow.Name}
		a.TeamStats[matchStats.TeamStatsBlue.Name] = &matchstats.TeamStats{Name: matchStats.TeamStatsBlue.Name}
	}

	for _, matchStats := range a.Collection.MatchStats {
		addTeamStats(a.TeamStats[matchStats.TeamStatsYellow.Name], matchStats.TeamStatsYellow)
		addTeamStats(a.TeamStats[matchStats.TeamStatsBlue.Name], matchStats.TeamStatsBlue)
	}

	return nil
}

func addTeamStats(to *matchstats.TeamStats, team *matchstats.TeamStats) {
	to.Goals += team.Goals
	to.ConcededGoals += team.ConcededGoals
	to.Fouls += team.Fouls
	to.YellowCards += team.YellowCards
	to.RedCards += team.RedCards
	to.TimeoutTime += team.TimeoutTime
	to.TimeoutsLeft += team.TimeoutsLeft
	to.TimeoutsTaken += team.TimeoutsTaken
	to.PenaltyShotsTotal += team.PenaltyShotsTotal
	to.BallPlacementTime += team.BallPlacementTime
	to.BallPlacements += team.BallPlacements
	if to.MaxActiveYellowCards < team.MaxActiveYellowCards {
		to.MaxActiveYellowCards = team.MaxActiveYellowCards
	}
}
