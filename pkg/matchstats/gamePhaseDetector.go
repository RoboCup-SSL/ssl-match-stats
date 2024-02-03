package matchstats

import (
	"crypto/md5"
	"fmt"
	"github.com/RoboCup-SSL/ssl-match-stats/internal/referee"
	"google.golang.org/protobuf/proto"
	"log"
	"time"
)

func (g *Generator) startNewGamePhase(ref *referee.Referee, phaseType GamePhaseType) {
	g.stopCurrentGamePhase(ref)
	prevPhase := g.currentPhase
	g.currentPhase = new(GamePhase)
	g.currentPhase.Type = phaseType
	g.currentPhase.StartTime = *ref.PacketTimestamp
	g.currentPhase.Stage = mapProtoStageToStageType(*ref.Stage)
	if ref.StageTimeLeft != nil {
		g.currentPhase.StageTimeLeftEntry = *ref.StageTimeLeft
	}
	g.currentPhase.CommandEntry = mapProtoCommandToCommand(*ref.Command, *ref.CommandTimestamp)
	g.currentPhase.ForTeam = mapProtoCommandToTeam(*ref.Command)
	g.currentPhase.GameEventsEntry = ref.GameEvents[:]

	if prevPhase != nil {
		g.currentPhase.CommandPrev = prevPhase.CommandEntry
	} else {
		g.currentPhase.CommandPrev = &Command{
			Type:      CommandType_COMMAND_UNKNOWN,
			ForTeam:   TeamColor_TEAM_UNKNOWN,
			Timestamp: *ref.CommandTimestamp,
		}
	}

	if g.currentPhase.CommandEntry.Type == CommandType_COMMAND_UNKNOWN {
		log.Println("Warn")
	}
}

func (g *Generator) stopCurrentGamePhase(ref *referee.Referee) {
	if g.currentPhase == nil {
		return
	}
	g.currentPhase.EndTime = *ref.PacketTimestamp
	start := packetTimestampToTime(g.currentPhase.StartTime)
	end := packetTimestampToTime(g.currentPhase.EndTime)
	g.currentPhase.Duration = end.Sub(start).Microseconds()
	g.matchStats.GamePhases = append(g.matchStats.GamePhases, g.currentPhase)
	g.currentPhase.CommandExit = mapProtoCommandToCommand(*ref.Command, *ref.CommandTimestamp)
	if ref.NextCommand != nil && int32(*ref.NextCommand) >= 0 {
		g.currentPhase.NextCommandProposed = mapProtoCommandToCommand(*ref.NextCommand, *ref.CommandTimestamp)
	}
	g.currentPhase.GameEventsExit = ref.GameEvents[:]
	if ref.StageTimeLeft != nil {
		g.currentPhase.StageTimeLeftExit = *ref.StageTimeLeft
	}
	g.finalizeRobotCount(ref, TeamColor_TEAM_YELLOW)
	g.finalizeRobotCount(ref, TeamColor_TEAM_BLUE)
}

func (g *Generator) handleNewStageForGamePhases(ref *referee.Referee) {
	switch *ref.Stage {
	case referee.Referee_NORMAL_FIRST_HALF_PRE,
		referee.Referee_NORMAL_FIRST_HALF,
		referee.Referee_NORMAL_SECOND_HALF_PRE,
		referee.Referee_NORMAL_SECOND_HALF,
		referee.Referee_EXTRA_FIRST_HALF_PRE,
		referee.Referee_EXTRA_SECOND_HALF_PRE,
		referee.Referee_EXTRA_FIRST_HALF,
		referee.Referee_EXTRA_SECOND_HALF,
		referee.Referee_PENALTY_SHOOTOUT:
		g.gamePaused = false
		break
	case referee.Referee_NORMAL_HALF_TIME,
		referee.Referee_EXTRA_TIME_BREAK,
		referee.Referee_EXTRA_HALF_TIME,
		referee.Referee_PENALTY_SHOOTOUT_BREAK:
		g.gamePaused = true
		g.startNewGamePhase(ref, GamePhaseType_PHASE_BREAK)
		break
	case referee.Referee_POST_GAME:
		g.startNewGamePhase(ref, GamePhaseType_PHASE_POST_GAME)
		break
	default:
		log.Println("Unknown stage: ", *ref.Stage)
	}
}

func (g *Generator) processGameEvents(ref *referee.Referee) {
	if g.currentPhase == nil {
		return
	}
	g.currentPhase.GameEventsApplied = processGameEvents(g.currentPhase.GameEventsApplied, ref.GameEvents, *ref.PacketTimestamp)

	var proposedGameEvents []*referee.GameEvent
	for _, gameEvent := range ref.GameEventProposals {
		proposedGameEvents = append(proposedGameEvents, gameEvent.GameEvents...)
	}

	g.currentPhase.GameEventsProposed = processGameEvents(g.currentPhase.GameEventsProposed, proposedGameEvents, *ref.PacketTimestamp)
}

func (g *Generator) processRobotCount(ref *referee.Referee, teamColor TeamColor) {
	curRobotCount := g.currentRobotCount[teamColor]
	if curRobotCount == nil || g.getCurrentRobotCount(teamColor) != curRobotCount.Count {
		g.finalizeRobotCount(ref, teamColor)
		g.currentRobotCount[teamColor] = &RobotCount{
			TeamColor: teamColor,
			Count:     g.getCurrentRobotCount(teamColor),
			StartTime: *ref.PacketTimestamp,
		}
	}
}

func (g *Generator) finalizeRobotCount(ref *referee.Referee, teamColor TeamColor) {
	curRobotCount := g.currentRobotCount[teamColor]
	if curRobotCount == nil {
		return
	}
	start := packetTimestampToTime(curRobotCount.StartTime)
	end := packetTimestampToTime(*ref.PacketTimestamp)
	duration := end.Sub(start)
	curRobotCount.Duration = duration.Microseconds()
	if duration > time.Millisecond*100 {
		g.currentPhase.RobotCount = append(g.currentPhase.RobotCount, curRobotCount)
	}
	g.currentRobotCount[teamColor] = nil
}

func (g *Generator) getCurrentRobotCount(teamColor TeamColor) int32 {
	count := int32(0)
	for id, last := range g.robotLastDetection[teamColor] {
		first := g.robotFirstDetection[teamColor][id]
		age := last - first
		if age > 0.5 {
			count++
		}
	}
	return count
}

func processGameEvents(gameEvents []*GameEventTimed, newGameEvents []*referee.GameEvent, timestamp uint64) []*GameEventTimed {

	// add game event id, if there is none yet
	for _, newGameEvent := range newGameEvents {
		upsertGameEventId(newGameEvent)
	}

	// mark game events as withdrawn, if they are not in the new game events anymore
	for _, gameEvent := range gameEvents {
		if !containsGameEvent(gameEvent.GameEvent, newGameEvents) {
			gameEvent.Withdrawn = true
		}
	}

	// add new game events
	for _, newGameEvent := range newGameEvents {
		if !containsGameEventTimed(newGameEvent, gameEvents) {
			gameEventTimed := GameEventTimed{
				GameEvent: newGameEvent,
				Timestamp: timestamp,
				Withdrawn: false,
			}
			gameEvents = append(gameEvents, &gameEventTimed)
		}
	}

	return gameEvents
}

func upsertGameEventId(gameEvent *referee.GameEvent) {
	if gameEvent.Id != nil {
		return
	}
	b, err := proto.Marshal(gameEvent)
	if err != nil {
		log.Fatal("Could not marshal game event: ", gameEvent, err)
	}
	h := md5.New()
	h.Write(b)
	gameEvent.Id = new(string)
	*gameEvent.Id = fmt.Sprintf("%x", h.Sum(nil))
}

func containsGameEventTimed(gameEvent *referee.GameEvent, gameEvents []*GameEventTimed) bool {
	for _, existingEvent := range gameEvents {
		if *existingEvent.GameEvent.Id == *gameEvent.Id {
			return true
		}
	}
	return false
}

func containsGameEvent(gameEvent *referee.GameEvent, gameEvents []*referee.GameEvent) bool {
	for _, existingEvent := range gameEvents {
		if *existingEvent.Id == *gameEvent.Id {
			return true
		}
	}
	return false
}
