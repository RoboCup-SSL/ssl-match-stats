package matchstats

import (
	"github.com/RoboCup-SSL/ssl-go-tools/pkg/sslproto"
	"log"
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

func (m *MetaDataProcessor) OnNewStage(matchStats *MatchStats, referee *sslproto.Referee) {
	if *referee.Stage == sslproto.Referee_EXTRA_TIME_BREAK {
		addTimeout(matchStats.TeamStatsYellow, referee.Yellow, m.timeoutTimeNormal, m.timeoutsNormal)
		addTimeout(matchStats.TeamStatsBlue, referee.Blue, m.timeoutTimeNormal, m.timeoutsNormal)
	}
	if *referee.Stage == sslproto.Referee_PENALTY_SHOOTOUT {
		matchStats.Shootout = true
	}
}

func (m *MetaDataProcessor) OnNewCommand(matchStats *MatchStats, referee *sslproto.Referee) {
	switch *referee.Command {
	case sslproto.Referee_PREPARE_PENALTY_BLUE:
		matchStats.TeamStatsBlue.PenaltyShotsTotal += 1
	case sslproto.Referee_PREPARE_PENALTY_YELLOW:
		matchStats.TeamStatsYellow.PenaltyShotsTotal += 1
	}
}

func (m *MetaDataProcessor) OnFirstRefereeMessage(matchStats *MatchStats, referee *sslproto.Referee) {
	if matchStats.TeamStatsBlue == nil {
		matchStats.TeamStatsBlue = new(TeamStats)
	}
	if matchStats.TeamStatsYellow == nil {
		matchStats.TeamStatsYellow = new(TeamStats)
	}

	matchStats.Shootout = false
	m.startTime = packetTimeStampToTime(*referee.PacketTimestamp)
	matchStats.StartTime = uint64(m.startTime.UnixNano() / 1000)
}

func (m *MetaDataProcessor) OnLastRefereeMessage(matchStats *MatchStats, referee *sslproto.Referee) {
	processTeam(matchStats.TeamStatsBlue, referee.Blue, referee.Yellow)
	processTeam(matchStats.TeamStatsYellow, referee.Yellow, referee.Blue)
	endTime := packetTimeStampToTime(*referee.PacketTimestamp)
	matchStats.MatchDuration = uint32(endTime.Sub(m.startTime).Microseconds())

	if uint32(*referee.Stage) <= uint32(sslproto.Referee_NORMAL_SECOND_HALF) {
		addTimeout(matchStats.TeamStatsYellow, referee.Yellow, m.timeoutTimeNormal, m.timeoutsNormal)
		addTimeout(matchStats.TeamStatsBlue, referee.Blue, m.timeoutTimeNormal, m.timeoutsNormal)
		matchStats.ExtraTime = false
	} else {
		addTimeout(matchStats.TeamStatsYellow, referee.Yellow, m.timeoutTimeExtra, m.timeoutsExtra)
		addTimeout(matchStats.TeamStatsBlue, referee.Blue, m.timeoutTimeExtra, m.timeoutsExtra)
		matchStats.ExtraTime = true
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
		}
	}
}

func addTimeout(teamStats *TeamStats, teamInfo *sslproto.Referee_TeamInfo, availableTime uint32, availableTimeouts uint32) {
	if *teamInfo.TimeoutTime > availableTime {
		log.Printf("Timeout time > available: %v > %v", *teamInfo.TimeoutTime, availableTime)
		if availableTime == 150_000_000 {
			availableTime = 300_000_000
			log.Println("Fixing known bug: In 2019, timeout time in extra halves was 5min instead of 2.5min")
		}
	}

	var timeouts int32
	if *teamInfo.Timeouts > 100 {
		// workaround for integer overflow
		timeouts = int32(int64(*teamInfo.Timeouts) - math.MaxUint32)
	} else {
		timeouts = int32(*teamInfo.Timeouts)
	}

	teamStats.TimeoutTime += availableTime - *teamInfo.TimeoutTime
	teamStats.TimeoutsTaken += int32(availableTimeouts) - timeouts

	teamStats.TimeoutsLeft += timeouts
}

func (m *MetaDataProcessor) OnNewRefereeMessage(matchStats *MatchStats, referee *sslproto.Referee) {
	m.updateMaxActiveYellowCards(referee.Blue, matchStats.TeamStatsBlue)
	m.updateMaxActiveYellowCards(referee.Yellow, matchStats.TeamStatsYellow)
}

func (m *MetaDataProcessor) updateMaxActiveYellowCards(teamInfo *sslproto.Referee_TeamInfo, teamStats *TeamStats) {
	activeCards := int32(len(teamInfo.YellowCardTimes))
	if teamStats.MaxActiveYellowCards < activeCards {
		teamStats.MaxActiveYellowCards = activeCards
	}
}

func processTeam(stats *TeamStats, team *sslproto.Referee_TeamInfo, otherTeam *sslproto.Referee_TeamInfo) {
	stats.Name = *team.Name
	stats.Goals = int32(*team.Score)
	stats.ConcededGoals = int32(*otherTeam.Score)
	stats.YellowCards = int32(*team.YellowCards)
	stats.RedCards = int32(*team.RedCards)
	if team.FoulCounter != nil {
		stats.Fouls = int32(*team.FoulCounter)
	}
}
