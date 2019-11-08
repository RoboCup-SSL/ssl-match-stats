package matchstats

import (
	"github.com/RoboCup-SSL/ssl-go-tools/pkg/persistence"
	"github.com/RoboCup-SSL/ssl-go-tools/pkg/sslproto"
	"github.com/pkg/errors"
	"log"
	"math"
	"path/filepath"
	"sort"
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

	AggregateGamePhaseStats(matchStats)
	AggregateGameEvents(matchStats)

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
	m.gamePhaseDetector.OnNewRefereeMessage(matchStats, referee)
}

func packetTimeStampToTime(packetTimestamp uint64) time.Time {
	seconds := int64(packetTimestamp / 1_000_000)
	nanoSeconds := int64(packetTimestamp-uint64(seconds*1_000_000)) * 1000
	return time.Unix(seconds, nanoSeconds)
}

func AggregateGamePhaseStats(matchStats *MatchStats) {

	matchStats.GamePhaseDurations = map[string]*DurationStats{}
	durations := map[string][]int{}

	for _, phaseName := range GamePhaseType_name {
		matchStats.GamePhaseDurations[phaseName] = new(DurationStats)
	}

	for _, phase := range matchStats.GamePhases {
		phaseName := (*phase).Type.String()
		matchStats.GamePhaseDurations[phaseName].Duration += phase.Duration
		matchStats.GamePhaseDurations[phaseName].Count += 1
		durations[phaseName] = append(durations[phaseName], int(phase.Duration))
	}

	for _, phaseName := range GamePhaseType_name {
		stats := matchStats.GamePhaseDurations[phaseName]
		phaseDurations := durations[phaseName]
		if len(phaseDurations) > 0 {
			sort.Ints(phaseDurations)
			stats.DurationMin = uint32(phaseDurations[0])
			stats.DurationMax = uint32(phaseDurations[len(phaseDurations)-1])
			stats.DurationMedian = uint32(phaseDurations[len(phaseDurations)/2])
			stats.DurationAvg = uint32(math.Round(float64(stats.Duration) / float64(len(phaseDurations))))
		}
	}

	checkSum := uint32(0)
	for _, phaseName := range GamePhaseType_name {
		checkSum += matchStats.GamePhaseDurations[phaseName].Duration
		matchStats.GamePhaseDurations[phaseName].DurationRelative = float32(matchStats.GamePhaseDurations[phaseName].Duration) / float32(matchStats.MatchDuration)
	}

	if matchStats.MatchDuration != checkSum {
		log.Printf("Match duration mismatch. Total: %v, Sum of phases: %v, Diff: %v", matchStats.MatchDuration, checkSum, matchStats.MatchDuration-checkSum)
	}
}

func AggregateGameEvents(matchStats *MatchStats) {

	matchStats.GameEventDurations = map[string]*DurationStats{}
	durations := map[string][]int{}

	for _, p := range sslproto.GameEventType_name {
		matchStats.GameEventDurations[p] = new(DurationStats)
		matchStats.GameEventDurations[p].Duration = 0
	}

	for _, phase := range matchStats.GamePhases {
		if len(phase.GameEventsEntry) == 0 {
			continue
		}

		primaryEvent := phase.GameEventsEntry[0]
		eventName := primaryEvent.Type.String()
		matchStats.GameEventDurations[eventName].Duration += phase.Duration
		matchStats.GameEventDurations[eventName].Count += 1
		durations[eventName] = append(durations[eventName], int(phase.Duration))
	}

	for _, eventName := range sslproto.GameEventType_name {
		stats := matchStats.GameEventDurations[eventName]
		eventDurations := durations[eventName]
		if len(eventDurations) > 0 {
			sort.Ints(eventDurations)
			stats.DurationMin = uint32(eventDurations[0])
			stats.DurationMax = uint32(eventDurations[len(eventDurations)-1])
			stats.DurationMedian = uint32(eventDurations[len(eventDurations)/2])
			stats.DurationAvg = uint32(math.Round(float64(stats.Duration) / float64(len(eventDurations))))
		}
	}
}
