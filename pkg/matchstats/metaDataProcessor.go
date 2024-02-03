package matchstats

import (
	"github.com/RoboCup-SSL/ssl-match-stats/internal/referee"
	"math"
)

// map of team names that are not correct in the log files
var teamNameFixes = map[string]string{
	"Ri-One":           "Ri-one",
	"nAMeC":            "NAMeC",
	"Op-AmP":           "OP-AmP",
	"Robobulls":        "RoboBulls",
	"RobôCIn":          "RobôCin",
	"STOx's":           "STOx’s",
	"Tigers Mannheim":  "TIGERs Mannheim",
	"UMass MinuteBots": "UMass Minutebots",
}

func (g *Generator) handleNewStageForMetaData(ref *referee.Referee) {
	switch *ref.Stage {
	case referee.Referee_EXTRA_TIME_BREAK,
		referee.Referee_POST_GAME:
		addTimeout(g.matchStats.TeamStatsYellow, ref.Yellow)
		addTimeout(g.matchStats.TeamStatsBlue, ref.Blue)
	}

	if *ref.Stage == referee.Referee_EXTRA_FIRST_HALF {
		g.matchStats.ExtraTime = true
	}

	if *ref.Stage == referee.Referee_PENALTY_SHOOTOUT {
		g.matchStats.Shootout = true
	}
}

func (g *Generator) handlePenaltyKick(ref *referee.Referee) {
	switch *ref.Command {
	case referee.Referee_PREPARE_PENALTY_BLUE:
		g.penaltyKickTeam = TeamColor_TEAM_BLUE
	case referee.Referee_PREPARE_PENALTY_YELLOW:
		g.penaltyKickTeam = TeamColor_TEAM_YELLOW
	case referee.Referee_NORMAL_START:
		if g.penaltyKickTeam == TeamColor_TEAM_BLUE {
			g.matchStats.TeamStatsBlue.PenaltyShotsTotal += 1
		} else if g.penaltyKickTeam == TeamColor_TEAM_YELLOW {
			g.matchStats.TeamStatsYellow.PenaltyShotsTotal += 1
		}
		g.penaltyKickTeam = TeamColor_TEAM_NONE
	default:
		g.penaltyKickTeam = TeamColor_TEAM_NONE
	}
}

func (g *Generator) OnFirstRefereeMessageMeta(ref *referee.Referee) {
	if g.matchStats.TeamStatsBlue == nil {
		g.matchStats.TeamStatsBlue = new(TeamStats)
	}
	if g.matchStats.TeamStatsYellow == nil {
		g.matchStats.TeamStatsYellow = new(TeamStats)
	}

	g.matchStats.Shootout = false
	g.startTime = packetTimestampToTime(*ref.PacketTimestamp)
	g.matchStats.StartTime = uint64(g.startTime.UnixNano() / 1000)
}

func (g *Generator) finalizeMatchStats(ref *referee.Referee) {
	processTeam(g.matchStats.TeamStatsBlue, ref.Blue, ref.Yellow)
	processTeam(g.matchStats.TeamStatsYellow, ref.Yellow, ref.Blue)
	endTime := packetTimestampToTime(*ref.PacketTimestamp)
	g.matchStats.MatchDuration = endTime.Sub(g.startTime).Microseconds()
	g.matchStats.Type = mapMatchType(ref.MatchType)

	// if the log file does not end with POST_GAME, we have to keep track of the remaining timeouts
	if uint32(*ref.Stage) != uint32(referee.Referee_POST_GAME) {
		addTimeout(g.matchStats.TeamStatsYellow, ref.Yellow)
		addTimeout(g.matchStats.TeamStatsBlue, ref.Blue)
	}

	for _, gamePhase := range g.matchStats.GamePhases {
		if gamePhase.Type == GamePhaseType_PHASE_BALL_PLACEMENT {
			if gamePhase.ForTeam == TeamColor_TEAM_BLUE {
				g.matchStats.TeamStatsBlue.BallPlacementTime += gamePhase.Duration
				g.matchStats.TeamStatsBlue.BallPlacements++
			} else if gamePhase.ForTeam == TeamColor_TEAM_YELLOW {
				g.matchStats.TeamStatsYellow.BallPlacementTime += gamePhase.Duration
				g.matchStats.TeamStatsYellow.BallPlacements++
			}
		} else if gamePhase.Type == GamePhaseType_PHASE_TIMEOUT {
			if gamePhase.ForTeam == TeamColor_TEAM_BLUE {
				g.matchStats.TeamStatsBlue.TimeoutsTaken++
				g.matchStats.TeamStatsBlue.TimeoutTime += gamePhase.Duration
			} else if gamePhase.ForTeam == TeamColor_TEAM_YELLOW {
				g.matchStats.TeamStatsYellow.TimeoutsTaken++
				g.matchStats.TeamStatsYellow.TimeoutTime += gamePhase.Duration
			}
		}
	}
}

func mapMatchType(matchType *referee.MatchType) StatsMatchType {
	if matchType == nil {
		return StatsMatchType_MATCH_UNKNOWN
	}
	switch *matchType {
	case referee.MatchType_UNKNOWN_MATCH:
		return StatsMatchType_MATCH_UNKNOWN
	case referee.MatchType_GROUP_PHASE:
		return StatsMatchType_MATCH_GROUP_PHASE
	case referee.MatchType_ELIMINATION_PHASE:
		return StatsMatchType_MATCH_ELIMINATION_PHASE
	case referee.MatchType_FRIENDLY:
		return StatsMatchType_MATCH_FRIENDLY
	}
	return StatsMatchType_MATCH_UNKNOWN
}

func addTimeout(teamStats *TeamStats, teamInfo *referee.Referee_TeamInfo) {
	var timeoutsLeft int32
	if *teamInfo.Timeouts > 100 {
		// workaround for integer overflow
		timeoutsLeft = int32(int64(*teamInfo.Timeouts) - math.MaxUint32)
	} else {
		timeoutsLeft = int32(*teamInfo.Timeouts)
	}

	teamStats.TimeoutsLeft += timeoutsLeft
}

func (g *Generator) updateMaxActiveYellowCards(teamInfo *referee.Referee_TeamInfo, teamStats *TeamStats) {
	activeCards := int32(len(teamInfo.YellowCardTimes))
	if teamStats.MaxActiveYellowCards < activeCards {
		teamStats.MaxActiveYellowCards = activeCards
	}
}

func processTeam(stats *TeamStats, team *referee.Referee_TeamInfo, otherTeam *referee.Referee_TeamInfo) {
	stats.Name = fixTeamName(*team.Name)
	stats.Goals = int32(*team.Score)
	stats.ConcededGoals = int32(*otherTeam.Score)
	stats.YellowCards = int32(*team.YellowCards)
	stats.RedCards = int32(*team.RedCards)
	if team.FoulCounter != nil {
		stats.Fouls = int32(*team.FoulCounter)
	}
}

func fixTeamName(teamName string) string {
	if fixedName, ok := teamNameFixes[teamName]; ok {
		return fixedName
	}
	return teamName
}
