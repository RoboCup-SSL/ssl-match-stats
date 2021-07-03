package matchstats

import (
	"github.com/RoboCup-SSL/ssl-match-stats/internal/referee"
	"math"
	"time"
)

type MetaDataProcessor struct {
	startTime         time.Time
	timeoutTimeNormal uint32
	timeoutsNormal    uint32
	timeoutTimeExtra  uint32
	timeoutsExtra     uint32
}

func NewMetaDataProcessor() *MetaDataProcessor {
	metaDataProcessor := new(MetaDataProcessor)
	// the referee messages only contain the remaining timeout time
	// as we can not guarantee that log files a complete, we do not know for sure, how much timeout time
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
		matchStats.TeamStatsBlue.PenaltyShotsTotal += 1
	case referee.Referee_PREPARE_PENALTY_YELLOW:
		matchStats.TeamStatsYellow.PenaltyShotsTotal += 1
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
	m.startTime = packetTimeStampToTime(*ref.PacketTimestamp)
	matchStats.StartTime = uint64(m.startTime.UnixNano() / 1000)
}

func (m *MetaDataProcessor) OnLastRefereeMessage(matchStats *MatchStats, ref *referee.Referee) {
	processTeam(matchStats.TeamStatsBlue, ref.Blue, ref.Yellow)
	processTeam(matchStats.TeamStatsYellow, ref.Yellow, ref.Blue)
	endTime := packetTimeStampToTime(*ref.PacketTimestamp)
	matchStats.MatchDuration = uint32(endTime.Sub(m.startTime).Microseconds())

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
	stats.Name = *team.Name
	stats.Goals = int32(*team.Score)
	stats.ConcededGoals = int32(*otherTeam.Score)
	stats.YellowCards = int32(*team.YellowCards)
	stats.RedCards = int32(*team.RedCards)
	if team.FoulCounter != nil {
		stats.Fouls = int32(*team.FoulCounter)
	}
}
