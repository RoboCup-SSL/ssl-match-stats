package matchstats

import (
	"github.com/RoboCup-SSL/ssl-match-stats/internal/referee"
	"log"
	"time"
)

func packetTimestampToTime(packetTimestamp uint64) time.Time {
	seconds := int64(packetTimestamp / 1_000_000)
	nanoSeconds := int64(packetTimestamp-uint64(seconds*1_000_000)) * 1000
	return time.Unix(seconds, nanoSeconds)
}

func mapProtoCommandToCommand(command referee.Referee_Command, timestamp uint64) *Command {
	return &Command{
		Type:      mapProtoCommandToCommandType(command),
		ForTeam:   mapProtoCommandToTeam(command),
		Timestamp: timestamp,
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
