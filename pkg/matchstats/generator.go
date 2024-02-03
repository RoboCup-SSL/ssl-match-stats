package matchstats

import (
	"github.com/RoboCup-SSL/ssl-go-tools/pkg/persistence"
	"github.com/RoboCup-SSL/ssl-match-stats/internal/referee"
	"github.com/RoboCup-SSL/ssl-match-stats/internal/vision"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
	"log"
	"path/filepath"
	"time"
)

type Generator struct {
	matchStats *MatchStats

	currentPhase        *GamePhase
	currentRobotCount   map[TeamColor]*RobotCount
	robotFirstDetection map[TeamColor]map[uint32]float64
	robotLastDetection  map[TeamColor]map[uint32]float64
	gamePaused          bool
	startTime           time.Time
	penaltyKickTeam     TeamColor
}

func NewGenerator() *Generator {
	return &Generator{
		currentRobotCount: map[TeamColor]*RobotCount{
			TeamColor_TEAM_YELLOW: nil,
			TeamColor_TEAM_BLUE:   nil,
		},
		robotFirstDetection: map[TeamColor]map[uint32]float64{
			TeamColor_TEAM_YELLOW: make(map[uint32]float64),
			TeamColor_TEAM_BLUE:   make(map[uint32]float64),
		},
		robotLastDetection: map[TeamColor]map[uint32]float64{
			TeamColor_TEAM_YELLOW: make(map[uint32]float64),
			TeamColor_TEAM_BLUE:   make(map[uint32]float64),
		},
	}
}

func (g *Generator) Process(filename string) (*MatchStats, error) {

	logReader, err := persistence.NewReader(filename)
	if err != nil {
		return nil, errors.Wrap(err, "Could not read file")
	}

	g.matchStats = new(MatchStats)
	g.matchStats.Name = filepath.Base(filename)
	var lastRefereeMsg *referee.Referee

	channel := logReader.CreateChannel()
	for c := range channel {
		if c.MessageType.Id == persistence.MessageSslVision2014 {
			v, err := getVisionMsg(c)
			if err != nil {
				log.Println("Could not parse vision message: ", err)
				continue
			}
			g.OnNewVisionMessage(v)
			continue
		}

		if c.MessageType.Id != persistence.MessageSslRefbox2013 {
			continue
		}
		r, err := getRefereeMsg(c)
		if err != nil {
			log.Println("Could not parse referee message: ", err)
			continue
		}

		if lastRefereeMsg == nil {
			g.OnFirstRefereeMessageMeta(r)
		} else if *r.PacketTimestamp < *lastRefereeMsg.PacketTimestamp {
			log.Printf("Skip out of order referee message:\nPrev:%v\nCurr:%v\n", lastRefereeMsg, r)
			continue
		}

		if lastRefereeMsg == nil || *r.Stage != *lastRefereeMsg.Stage {
			g.OnNewStage(r)
		}

		if lastRefereeMsg == nil || *r.Command != *lastRefereeMsg.Command {
			g.OnNewCommand(r)
		}

		g.OnNewRefereeMessage(r)

		lastRefereeMsg = r
	}

	if lastRefereeMsg != nil {
		g.OnLastRefereeMessage(lastRefereeMsg)
	}

	return g.matchStats, logReader.Close()
}

func getRefereeMsg(logMessage *persistence.Message) (refereeMsg *referee.Referee, err error) {
	refereeMsg = new(referee.Referee)
	if err := proto.Unmarshal(logMessage.Message, refereeMsg); err != nil {
		err = errors.Wrap(err, "Could not parse referee message")
	}
	return
}

func getVisionMsg(logMessage *persistence.Message) (visionMsg *vision.SSL_WrapperPacket, err error) {
	visionMsg = new(vision.SSL_WrapperPacket)
	if err := proto.Unmarshal(logMessage.Message, visionMsg); err != nil {
		err = errors.Wrap(err, "Could not parse vision message")
	}
	return
}

func (g *Generator) OnNewStage(referee *referee.Referee) {
	g.handleNewStageForMetaData(referee)
	g.handleNewStageForGamePhases(referee)
}

func (g *Generator) OnNewCommand(ref *referee.Referee) {
	g.handlePenaltyKick(ref)

	if !g.gamePaused {

		phaseType := mapProtoCommandToGamePhaseType(*ref.Command)
		if phaseType != GamePhaseType_PHASE_UNKNOWN {
			g.startNewGamePhase(ref, phaseType)
		}
	}
}

func (g *Generator) OnLastRefereeMessage(ref *referee.Referee) {
	g.finalizeMatchStats(ref)
	g.stopCurrentGamePhase(ref)
}

func (g *Generator) OnNewRefereeMessage(ref *referee.Referee) {
	g.updateMaxActiveYellowCards(ref.Blue, g.matchStats.TeamStatsBlue)
	g.updateMaxActiveYellowCards(ref.Yellow, g.matchStats.TeamStatsYellow)
	g.processGameEvents(ref)
	g.processRobotCount(ref, TeamColor_TEAM_YELLOW)
	g.processRobotCount(ref, TeamColor_TEAM_BLUE)
}
