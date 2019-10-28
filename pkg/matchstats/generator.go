package matchstats

import (
	"github.com/RoboCup-SSL/ssl-go-tools/pkg/persistence"
	"github.com/RoboCup-SSL/ssl-go-tools/pkg/sslproto"
	"github.com/pkg/errors"
	"log"
	"path/filepath"
	"time"
)

type Generator struct {
	metaDataProcessor MetaDataProcessor
	gamePhaseDetector GamePhaseDetector
}

func NewGenerator() *Generator {
	generator := new(Generator)
	generator.metaDataProcessor = MetaDataProcessor{}
	generator.gamePhaseDetector = GamePhaseDetector{}

	generator.metaDataProcessor.timeoutTimeNormal = 300_000_000
	generator.metaDataProcessor.timeoutTimeExtra = 300_000_000 // should be 150_000_000 for all except 2019 (wrong config in new GC)
	generator.metaDataProcessor.timeoutsNormal = 4
	generator.metaDataProcessor.timeoutsExtra = 2

	return generator
}

func (m *Generator) Process(filename string) (*MatchStats, error) {

	logReader, err := persistence.NewReader(filename)
	if err != nil {
		return nil, errors.Wrap(err, "Could not read file")
	}

	matchStats := new(MatchStats)
	matchStats.Name = filepath.Base(filename)
	var lastRefereeMsg *sslproto.Referee

	channel := logReader.CreateChannel()
	for c := range channel {
		if c.MessageType.Id != persistence.MessageSslRefbox2013 {
			continue
		}
		r, err := c.ParseReferee()
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

	Aggregate(matchStats)

	return matchStats, logReader.Close()
}

func (m *Generator) OnNewStage(matchStats *MatchStats, referee *sslproto.Referee) {
	m.metaDataProcessor.OnNewStage(matchStats, referee)
	m.gamePhaseDetector.OnNewStage(matchStats, referee)
}

func (m *Generator) OnNewCommand(matchStats *MatchStats, referee *sslproto.Referee) {
	m.metaDataProcessor.OnNewCommand(matchStats, referee)
	m.gamePhaseDetector.OnNewCommand(matchStats, referee)
}

func (m *Generator) OnFirstRefereeMessage(matchStats *MatchStats, referee *sslproto.Referee) {
	m.metaDataProcessor.OnFirstRefereeMessage(matchStats, referee)
}

func (m *Generator) OnLastRefereeMessage(matchStats *MatchStats, referee *sslproto.Referee) {
	m.metaDataProcessor.OnLastRefereeMessage(matchStats, referee)
	m.gamePhaseDetector.OnLastRefereeMessage(matchStats, referee)
}

func (m *Generator) OnNewRefereeMessage(matchStats *MatchStats, referee *sslproto.Referee) {
	m.metaDataProcessor.OnNewRefereeMessage(matchStats, referee)
}

func packetTimeStampToTime(packetTimestamp uint64) time.Time {
	seconds := int64(packetTimestamp / 1_000_000)
	nanoSeconds := int64(packetTimestamp-uint64(seconds*1_000_000)) * 1000
	return time.Unix(seconds, nanoSeconds)
}

func Aggregate(matchStats *MatchStats) {

	matchStats.GamePhaseStats = map[string]*GamePhaseStats{}

	for _, p := range GamePhaseType_name {
		matchStats.GamePhaseStats[p] = new(GamePhaseStats)
	}

	for _, p := range matchStats.GamePhases {
		matchStats.GamePhaseStats[(*p).Type.String()].Duration += p.Duration
	}

	checkSum := uint32(0)
	for _, p := range GamePhaseType_name {
		checkSum += matchStats.GamePhaseStats[p].Duration
		matchStats.GamePhaseStats[p].DurationRelative = float32(matchStats.GamePhaseStats[p].Duration) / float32(matchStats.MatchDuration)
	}

	if matchStats.MatchDuration != checkSum {
		log.Printf("Match duration mismatch. Total: %v, Sum of phases: %v, Diff: %v", matchStats.MatchDuration, checkSum, matchStats.MatchDuration-checkSum)
	}
}
