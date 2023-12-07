package matchstats

import (
	"github.com/RoboCup-SSL/ssl-match-stats/internal/referee"
	"math"
	"time"
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

type MetaDataProcessor struct {
	startTime time.Time
	// [microseconds]
	timeoutTimeNormal uint32
	timeoutsNormal    uint32
	// [microseconds]
	timeoutTimeExtra uint32
	timeoutsExtra    uint32

	penaltyKickTeam TeamColor
}

func NewMetaDataProcessor() *MetaDataProcessor {
	metaDataProcessor := new(MetaDataProcessor)
	// the referee messages only contain the remaining timeout time
	// as we can not guarantee that log files are complete, we do not know for sure, how much timeout time
	// was available initially, so we set it here explicitly based on the current rule set
	metaDataProcessor.timeoutTimeNormal = 300_000_000
	metaDataProcessor.timeoutTimeExtra = 150_000_000
	metaDataProcessor.timeoutsNormal = 4
	metaDataProcessor.timeoutsExtra = 2
	return metaDataProcessor
}

func (m *MetaDataProcessor) OnNewStage(matchStats *MatchStats, ref *referee.Referee) {
	switch *ref.Stage {
	case referee.Referee_EXTRA_TIME_BREAK,
		referee.Referee_POST_GAME:
		addTimeout(matchStats.TeamStatsYellow, ref.Yellow)
		addTimeout(matchStats.TeamStatsBlue, ref.Blue)
	}

	if *ref.Stage == referee.Referee_EXTRA_FIRST_HALF {
		matchStats.ExtraTime = true
	}

	if *ref.Stage == referee.Referee_PENALTY_SHOOTOUT {
		matchStats.Shootout = true
	}
}

func (m *MetaDataProcessor) OnNewCommand(matchStats *MatchStats, ref *referee.Referee) {
	switch *ref.Command {
	case referee.Referee_PREPARE_PENALTY_BLUE:
		m.penaltyKickTeam = TeamColor_TEAM_BLUE
	case referee.Referee_PREPARE_PENALTY_YELLOW:
		m.penaltyKickTeam = TeamColor_TEAM_YELLOW
	case referee.Referee_NORMAL_START:
		if m.penaltyKickTeam == TeamColor_TEAM_BLUE {
			matchStats.TeamStatsBlue.PenaltyShotsTotal += 1
		} else if m.penaltyKickTeam == TeamColor_TEAM_YELLOW {
			matchStats.TeamStatsYellow.PenaltyShotsTotal += 1
		}
		m.penaltyKickTeam = TeamColor_TEAM_NONE
	default:
		m.penaltyKickTeam = TeamColor_TEAM_NONE
	}
}

func (m *MetaDataProcessor) OnFirstRefereeMessage(matchStats *MatchStats, ref *referee.Referee) {
	if matchStats.TeamStatsBlue == nil {
		matchStats.TeamStatsBlue = new(TeamStats)
	}
	if matchStats.TeamStatsYellow == nil {
		matchStats.TeamStatsYellow = new(TeamStats)
	}

	matchStats.Shootout = false
	m.startTime = packetTimestampToTime(*ref.PacketTimestamp)
	matchStats.StartTime = uint64(m.startTime.UnixNano() / 1000)
}

func (m *MetaDataProcessor) OnLastRefereeMessage(matchStats *MatchStats, ref *referee.Referee) {
	processTeam(matchStats.TeamStatsBlue, ref.Blue, ref.Yellow)
	processTeam(matchStats.TeamStatsYellow, ref.Yellow, ref.Blue)
	endTime := packetTimestampToTime(*ref.PacketTimestamp)
	matchStats.MatchDuration = endTime.Sub(m.startTime).Microseconds()
	matchStats.Type = mapMatchType(ref.MatchType)

	// if the log file does not end with POST_GAME, we have to keep track of the remaining timeouts
	if uint32(*ref.Stage) != uint32(referee.Referee_POST_GAME) {
		addTimeout(matchStats.TeamStatsYellow, ref.Yellow)
		addTimeout(matchStats.TeamStatsBlue, ref.Blue)
	}

	for _, gamePhase := range matchStats.GamePhases {
		if gamePhase.Type == GamePhaseType_PHASE_BALL_PLACEMENT {
			if gamePhase.ForTeam == TeamColor_TEAM_BLUE {
				matchStats.TeamStatsBlue.BallPlacementTime += gamePhase.Duration
				matchStats.TeamStatsBlue.BallPlacements++
			} else if gamePhase.ForTeam == TeamColor_TEAM_YELLOW {
				matchStats.TeamStatsYellow.BallPlacementTime += gamePhase.Duration
				matchStats.TeamStatsYellow.BallPlacements++
			}
		} else if gamePhase.Type == GamePhaseType_PHASE_TIMEOUT {
			if gamePhase.ForTeam == TeamColor_TEAM_BLUE {
				matchStats.TeamStatsBlue.TimeoutsTaken++
				matchStats.TeamStatsBlue.TimeoutTime += gamePhase.Duration
			} else if gamePhase.ForTeam == TeamColor_TEAM_YELLOW {
				matchStats.TeamStatsYellow.TimeoutsTaken++
				matchStats.TeamStatsYellow.TimeoutTime += gamePhase.Duration
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

func (m *MetaDataProcessor) OnNewRefereeMessage(matchStats *MatchStats, referee *referee.Referee) {
	m.updateMaxActiveYellowCards(referee.Blue, matchStats.TeamStatsBlue)
	m.updateMaxActiveYellowCards(referee.Yellow, matchStats.TeamStatsYellow)
}

func (m *MetaDataProcessor) updateMaxActiveYellowCards(teamInfo *referee.Referee_TeamInfo, teamStats *TeamStats) {
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
