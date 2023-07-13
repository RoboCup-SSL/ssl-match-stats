package matchstats

import (
	"github.com/RoboCup-SSL/ssl-match-stats/internal/referee"
	"log"
	"reflect"
)

type GamePhaseDetector struct {
	currentPhase *GamePhase
	gamePaused   bool
}

func NewGamePhaseDetector() *GamePhaseDetector {
	return &GamePhaseDetector{gamePaused: false}
}

func (d *GamePhaseDetector) startNewPhase(matchStats *MatchStats, ref *referee.Referee, phaseType GamePhaseType) {
	d.stopCurrentPhase(matchStats, ref)
	prevPhase := d.currentPhase
	d.currentPhase = new(GamePhase)
	d.currentPhase.Type = phaseType
	d.currentPhase.StartTime = *ref.PacketTimestamp
	d.currentPhase.Stage = mapProtoStageToStageType(*ref.Stage)
	if ref.StageTimeLeft != nil {
		d.currentPhase.StageTimeLeftEntry = int64(*ref.StageTimeLeft)
	}
	d.currentPhase.CommandEntry = mapProtoCommandToCommand(*ref.Command)
	d.currentPhase.ForTeam = mapProtoCommandToTeam(*ref.Command)
	d.currentPhase.GameEventsEntry = ref.GameEvents[:]

	if prevPhase != nil {
		d.currentPhase.CommandPrev = prevPhase.CommandEntry
	} else {
		d.currentPhase.CommandPrev = &Command{Type: CommandType_COMMAND_UNKNOWN, ForTeam: TeamColor_TEAM_UNKNOWN}
	}

	if d.currentPhase.CommandEntry.Type == CommandType_COMMAND_UNKNOWN {
		log.Println("Warn")
	}
}

func (d *GamePhaseDetector) stopCurrentPhase(matchStats *MatchStats, ref *referee.Referee) {
	if d.currentPhase == nil {
		return
	}
	d.currentPhase.EndTime = *ref.PacketTimestamp
	start := packetTimeStampToTime(d.currentPhase.StartTime)
	end := packetTimeStampToTime(d.currentPhase.EndTime)
	d.currentPhase.Duration = end.Sub(start).Microseconds()
	matchStats.GamePhases = append(matchStats.GamePhases, d.currentPhase)
	d.currentPhase.CommandExit = mapProtoCommandToCommand(*ref.Command)
	if ref.NextCommand != nil && int32(*ref.NextCommand) >= 0 {
		d.currentPhase.NextCommandProposed = mapProtoCommandToCommand(*ref.NextCommand)
	}
	d.currentPhase.GameEventsExit = ref.GameEvents[:]
	if ref.StageTimeLeft != nil {
		d.currentPhase.StageTimeLeftExit = int64(*ref.StageTimeLeft)
	}
}

func (d *GamePhaseDetector) OnNewStage(matchStats *MatchStats, ref *referee.Referee) {
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
		d.gamePaused = false
		break
	case referee.Referee_NORMAL_HALF_TIME,
		referee.Referee_EXTRA_TIME_BREAK,
		referee.Referee_EXTRA_HALF_TIME,
		referee.Referee_PENALTY_SHOOTOUT_BREAK:
		d.gamePaused = true
		d.startNewPhase(matchStats, ref, GamePhaseType_PHASE_BREAK)
		break
	case referee.Referee_POST_GAME:
		d.startNewPhase(matchStats, ref, GamePhaseType_PHASE_POST_GAME)
		break
	default:
		log.Println("Unknown stage: ", *ref.Stage)
	}
}

func (d *GamePhaseDetector) OnNewCommand(matchStats *MatchStats, ref *referee.Referee) {
	if d.gamePaused {
		return
	}

	phaseType := mapProtoCommandToGamePhaseType(*ref.Command)
	if phaseType != GamePhaseType_PHASE_UNKNOWN {
		d.startNewPhase(matchStats, ref, phaseType)
	}
}

func (d *GamePhaseDetector) OnLastRefereeMessage(matchStats *MatchStats, ref *referee.Referee) {
	d.stopCurrentPhase(matchStats, ref)
}

func (d *GamePhaseDetector) OnNewRefereeMessage(_ *MatchStats, ref *referee.Referee) {
	if d.currentPhase == nil {
		return
	}
	d.currentPhase.GameEventsApplied = processGameEvents(d.currentPhase.GameEventsApplied, ref.GameEvents, *ref.PacketTimestamp)

	var proposedGameEvents []*referee.GameEvent
	for _, gameEvent := range ref.GameEventProposals {
		proposedGameEvents = append(proposedGameEvents, gameEvent.GameEvent...)
	}

	d.currentPhase.GameEventsProposed = processGameEvents(d.currentPhase.GameEventsProposed, proposedGameEvents, *ref.PacketTimestamp)
}

func processGameEvents(gameEvents []*GameEventTimed, newGameEvents []*referee.GameEvent, timestamp uint64) []*GameEventTimed {

	for _, presentGameEvent := range gameEvents {
		if !containsGameEvent(presentGameEvent.GameEvent, newGameEvents) {
			presentGameEvent.Withdrawn = true
		}
	}

	for _, presentGameEvent := range newGameEvents {
		if !containsGameEventTimed(presentGameEvent, gameEvents) {
			gameEventTimed := GameEventTimed{GameEvent: presentGameEvent, Timestamp: timestamp, Withdrawn: false}
			gameEvents = append(gameEvents, &gameEventTimed)
		}
	}

	return gameEvents
}

func containsGameEventTimed(gameEvent *referee.GameEvent, gameEvents []*GameEventTimed) bool {
	for _, existingEvent := range gameEvents {
		if reflect.DeepEqual(existingEvent.GameEvent, gameEvent) {
			return true
		}
	}
	return false
}

func containsGameEvent(gameEvent *referee.GameEvent, gameEvents []*referee.GameEvent) bool {
	for _, existingEvent := range gameEvents {
		if reflect.DeepEqual(existingEvent, gameEvent) {
			return true
		}
	}
	return false
}

func mapProtoCommandToCommand(command referee.Referee_Command) *Command {
	return &Command{
		Type:    mapProtoCommandToCommandType(command),
		ForTeam: mapProtoCommandToTeam(command),
	}
}

//goland:noinspection GoDeprecation
func mapProtoCommandToCommandType(command referee.Referee_Command) CommandType {
	switch command {
	case referee.Referee_HALT:
		return CommandType_COMMAND_HALT
	case referee.Referee_STOP:
		return CommandType_COMMAND_STOP
	case referee.Referee_NORMAL_START:
		return CommandType_COMMAND_NORMAL_START
	case referee.Referee_FORCE_START:
		return CommandType_COMMAND_FORCE_START
	case referee.Referee_PREPARE_KICKOFF_YELLOW,
		referee.Referee_PREPARE_KICKOFF_BLUE:
		return CommandType_COMMAND_PREPARE_KICKOFF
	case referee.Referee_PREPARE_PENALTY_YELLOW,
		referee.Referee_PREPARE_PENALTY_BLUE:
		return CommandType_COMMAND_PREPARE_PENALTY
	case referee.Referee_DIRECT_FREE_YELLOW,
		referee.Referee_DIRECT_FREE_BLUE:
		return CommandType_COMMAND_DIRECT_FREE
	case referee.Referee_INDIRECT_FREE_YELLOW,
		referee.Referee_INDIRECT_FREE_BLUE:
		return CommandType_COMMAND_INDIRECT_FREE
	case referee.Referee_TIMEOUT_YELLOW,
		referee.Referee_TIMEOUT_BLUE:
		return CommandType_COMMAND_TIMEOUT
	case referee.Referee_BALL_PLACEMENT_YELLOW,
		referee.Referee_BALL_PLACEMENT_BLUE:
		return CommandType_COMMAND_BALL_PLACEMENT
	case referee.Referee_GOAL_YELLOW, referee.Referee_GOAL_BLUE:
		return CommandType_COMMAND_GOAL
	}
	log.Printf("Command %v not mapped to any command type", command)
	return CommandType_COMMAND_UNKNOWN
}

//goland:noinspection GoDeprecation
func mapProtoCommandToGamePhaseType(command referee.Referee_Command) GamePhaseType {
	switch command {
	case referee.Referee_HALT:
		return GamePhaseType_PHASE_HALT
	case referee.Referee_STOP:
		return GamePhaseType_PHASE_STOP
	case referee.Referee_NORMAL_START:
		return GamePhaseType_PHASE_RUNNING
	case referee.Referee_FORCE_START:
		return GamePhaseType_PHASE_RUNNING
	case referee.Referee_PREPARE_KICKOFF_YELLOW, referee.Referee_PREPARE_KICKOFF_BLUE:
		return GamePhaseType_PHASE_PREPARE_KICKOFF
	case referee.Referee_PREPARE_PENALTY_YELLOW, referee.Referee_PREPARE_PENALTY_BLUE:
		return GamePhaseType_PHASE_PREPARE_PENALTY
	case referee.Referee_DIRECT_FREE_YELLOW, referee.Referee_DIRECT_FREE_BLUE:
		return GamePhaseType_PHASE_RUNNING
	case referee.Referee_INDIRECT_FREE_YELLOW, referee.Referee_INDIRECT_FREE_BLUE:
		return GamePhaseType_PHASE_RUNNING
	case referee.Referee_TIMEOUT_YELLOW, referee.Referee_TIMEOUT_BLUE:
		return GamePhaseType_PHASE_TIMEOUT
	case referee.Referee_BALL_PLACEMENT_YELLOW, referee.Referee_BALL_PLACEMENT_BLUE:
		return GamePhaseType_PHASE_BALL_PLACEMENT
	case referee.Referee_GOAL_YELLOW, referee.Referee_GOAL_BLUE:
		return GamePhaseType_PHASE_UNKNOWN
	}
	log.Printf("Command %v not mapped to any phase type", command)
	return GamePhaseType_PHASE_UNKNOWN
}

//goland:noinspection GoDeprecation
func mapProtoCommandToTeam(command referee.Referee_Command) TeamColor {
	switch command {
	case referee.Referee_HALT,
		referee.Referee_STOP,
		referee.Referee_NORMAL_START,
		referee.Referee_FORCE_START:
		return TeamColor_TEAM_NONE
	case referee.Referee_PREPARE_KICKOFF_YELLOW,
		referee.Referee_PREPARE_PENALTY_YELLOW,
		referee.Referee_DIRECT_FREE_YELLOW,
		referee.Referee_INDIRECT_FREE_YELLOW,
		referee.Referee_TIMEOUT_YELLOW,
		referee.Referee_BALL_PLACEMENT_YELLOW,
		referee.Referee_GOAL_YELLOW:
		return TeamColor_TEAM_YELLOW
	case referee.Referee_PREPARE_KICKOFF_BLUE,
		referee.Referee_PREPARE_PENALTY_BLUE,
		referee.Referee_DIRECT_FREE_BLUE,
		referee.Referee_INDIRECT_FREE_BLUE,
		referee.Referee_TIMEOUT_BLUE,
		referee.Referee_BALL_PLACEMENT_BLUE,
		referee.Referee_GOAL_BLUE:
		return TeamColor_TEAM_BLUE
	}
	log.Printf("Command %v not mapped to any team", command)
	return TeamColor_TEAM_UNKNOWN
}

func mapProtoStageToStageType(stage referee.Referee_Stage) StageType {
	switch stage {
	case referee.Referee_NORMAL_FIRST_HALF_PRE:
		return StageType_STAGE_NORMAL_FIRST_HALF_PRE
	case referee.Referee_NORMAL_FIRST_HALF:
		return StageType_STAGE_NORMAL_FIRST_HALF
	case referee.Referee_NORMAL_HALF_TIME:
		return StageType_STAGE_NORMAL_HALF_TIME
	case referee.Referee_NORMAL_SECOND_HALF_PRE:
		return StageType_STAGE_NORMAL_SECOND_HALF_PRE
	case referee.Referee_NORMAL_SECOND_HALF:
		return StageType_STAGE_NORMAL_SECOND_HALF
	case referee.Referee_EXTRA_TIME_BREAK:
		return StageType_STAGE_EXTRA_TIME_BREAK
	case referee.Referee_EXTRA_FIRST_HALF_PRE:
		return StageType_STAGE_EXTRA_FIRST_HALF_PRE
	case referee.Referee_EXTRA_FIRST_HALF:
		return StageType_STAGE_EXTRA_FIRST_HALF
	case referee.Referee_EXTRA_HALF_TIME:
		return StageType_STAGE_EXTRA_HALF_TIME
	case referee.Referee_EXTRA_SECOND_HALF_PRE:
		return StageType_STAGE_EXTRA_SECOND_HALF_PRE
	case referee.Referee_EXTRA_SECOND_HALF:
		return StageType_STAGE_EXTRA_SECOND_HALF
	case referee.Referee_PENALTY_SHOOTOUT_BREAK:
		return StageType_STAGE_PENALTY_SHOOTOUT_BREAK
	case referee.Referee_PENALTY_SHOOTOUT:
		return StageType_STAGE_PENALTY_SHOOTOUT
	case referee.Referee_POST_GAME:
		return StageType_STAGE_POST_GAME
	}
	return StageType_STAGE_UNKNOWN
}
