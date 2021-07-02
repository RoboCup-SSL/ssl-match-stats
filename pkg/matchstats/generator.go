package matchstats

import (
	"github.com/RoboCup-SSL/ssl-match-stats/internal/referee"
	"github.com/RoboCup-SSL/ssl-match-stats/pkg/persistence"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
	"log"
	"path/filepath"
	"time"
)

type Generator struct {
	metaDataProcessor *MetaDataProcessor
	gamePhaseDetector *GamePhaseDetector
}

func NewGenerator() *Generator {
	generator := new(Generator)
	generator.metaDataProcessor = NewMetaDataProcessor()
	generator.gamePhaseDetector = NewGamePhaseDetector()
	return generator
}

func (m *Generator) Process(filename string) (*MatchStats, error) {

	logReader, err := persistence.NewReader(filename)
	if err != nil {
		return nil, errors.Wrap(err, "Could not read file")
	}

	matchStats := new(MatchStats)
	matchStats.Name = filepath.Base(filename)
	var lastRefereeMsg *referee.Referee

	channel := logReader.CreateChannel()
	for c := range channel {
		if c.MessageType.Id != persistence.MessageSslRefbox2013 {
			continue
		}
		r, err := getRefereeMsg(c)
		if err != nil {
			log.Println("Could not parse referee message: ", err)
			continue
		}

		if lastRefereeMsg == nil {
			m.OnFirstRefereeMessage(matchStats, r)
		} else if *r.CommandCounter < *lastRefereeMsg.CommandCounter {
			log.Printf("Ignoring possible foreign referee command. Command counter %v < %v", *r.CommandCounter, *lastRefereeMsg.CommandCounter)
			continue
		}

		if lastRefereeMsg == nil || *r.Stage != *lastRefereeMsg.Stage {
			m.OnNewStage(matchStats, r)
		}

		if lastRefereeMsg == nil || *r.Command != *lastRefereeMsg.Command {
			m.OnNewCommand(matchStats, r)
		}

		m.OnNewRefereeMessage(matchStats, r)

		lastRefereeMsg = r
	}

	m.OnLastRefereeMessage(matchStats, lastRefereeMsg)

	return matchStats, logReader.Close()
}

func getRefereeMsg(logMessage *persistence.Message) (refereeMsg *referee.Referee, err error) {
	if logMessage.MessageType.Id != persistence.MessageSslRefbox2013 {
		return
	}

	refereeMsg = new(referee.Referee)
	if err := proto.Unmarshal(logMessage.Message, refereeMsg); err != nil {
		err = errors.Wrap(err, "Could not parse referee message")
	}
	return
}

func (m *Generator) OnNewStage(matchStats *MatchStats, referee *referee.Referee) {
	m.metaDataProcessor.OnNewStage(matchStats, referee)
	m.gamePhaseDetector.OnNewStage(matchStats, referee)
}

func (m *Generator) OnNewCommand(matchStats *MatchStats, referee *referee.Referee) {
	m.metaDataProcessor.OnNewCommand(matchStats, referee)
	m.gamePhaseDetector.OnNewCommand(matchStats, referee)
}

func (m *Generator) OnFirstRefereeMessage(matchStats *MatchStats, referee *referee.Referee) {
	m.metaDataProcessor.OnFirstRefereeMessage(matchStats, referee)
}

func (m *Generator) OnLastRefereeMessage(matchStats *MatchStats, referee *referee.Referee) {
	m.metaDataProcessor.OnLastRefereeMessage(matchStats, referee)
	m.gamePhaseDetector.OnLastRefereeMessage(matchStats, referee)
}

func (m *Generator) OnNewRefereeMessage(matchStats *MatchStats, referee *referee.Referee) {
	m.metaDataProcessor.OnNewRefereeMessage(matchStats, referee)
	m.gamePhaseDetector.OnNewRefereeMessage(matchStats, referee)
}

func packetTimeStampToTime(packetTimestamp uint64) time.Time {
	seconds := int64(packetTimestamp / 1_000_000)
	nanoSeconds := int64(packetTimestamp-uint64(seconds*1_000_000)) * 1000
	return time.Unix(seconds, nanoSeconds)
}
