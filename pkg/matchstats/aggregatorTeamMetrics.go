package matchstats

func (a *Aggregator) AggregateTeamMetrics() error {
	a.Collection.TeamStats = map[string]*TeamStats{}
	for _, matchStats := range a.Collection.MatchStats {
		a.Collection.TeamStats[matchStats.TeamStatsYellow.Name] = &TeamStats{Name: matchStats.TeamStatsYellow.Name}
		a.Collection.TeamStats[matchStats.TeamStatsBlue.Name] = &TeamStats{Name: matchStats.TeamStatsBlue.Name}
	}

	for _, matchStats := range a.Collection.MatchStats {
		addTeamStats(a.Collection.TeamStats[matchStats.TeamStatsYellow.Name], matchStats.TeamStatsYellow)
		addTeamStats(a.Collection.TeamStats[matchStats.TeamStatsBlue.Name], matchStats.TeamStatsBlue)
	}

	return nil
}

func addTeamStats(to *TeamStats, team *TeamStats) {
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
