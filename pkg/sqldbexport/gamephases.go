package sqldbexport

import (
	"github.com/RoboCup-SSL/ssl-match-stats/pkg/matchstats"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (p *SqlDbExporter) removeOldMatchData(matchId *uuid.UUID) error {
	// delete all game phases and all foreign key references to it by cascading
	if _, err := p.db.Exec("DELETE FROM game_phases WHERE match_id_fk=$1", matchId); err != nil {
		return errors.Wrap(err, "Could not delete previous game phases for id: "+matchId.String())
	}
	return nil
}

func (p *SqlDbExporter) WriteGamePhases(matchStats *matchstats.MatchStats, matchId *uuid.UUID) error {
	if err := p.removeOldMatchData(matchId); err != nil {
		return err
	}

	for _, gamePhase := range matchStats.GamePhases {
		if err := p.insertGamePhase(gamePhase, matchId); err != nil {
			return err
		}
	}
	return nil
}

func (p *SqlDbExporter) insertGamePhase(gamePhase *matchstats.GamePhase, matchId *uuid.UUID) error {
	var nextCommandType *string
	var nextCommandForTeam *string
	if gamePhase.NextCommandProposed != nil {
		nextCommandType = new(string)
		*nextCommandType = gamePhase.NextCommandProposed.Type.String()[8:]
		nextCommandForTeam = new(string)
		*nextCommandForTeam = gamePhase.NextCommandProposed.ForTeam.String()[5:]
	}

	id := uuid.New()

	_, err := p.db.Exec(
		`INSERT INTO game_phases (
                     id, 
                     match_id_fk, 
                     start_time,
					 end_time,
					 duration,
					 type,
					 for_team,
					 entry_command,
					 entry_command_for_team,
					 exit_command,
					 exit_command_for_team,
					 proposed_next_command,
					 proposed_next_command_for_team,
					 previous_command,
					 previous_command_for_team,
					 stage,
					 stage_time_left_entry,
					 stage_time_left_exit
                     ) 
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)`,
		id,
		matchId,
		convertTime(gamePhase.StartTime),
		convertTime(gamePhase.EndTime),
		convertDuration(gamePhase.Duration),
		gamePhase.Type.String()[6:],
		gamePhase.ForTeam.String()[5:],
		gamePhase.CommandEntry.Type.String()[8:],
		gamePhase.CommandEntry.ForTeam.String()[5:],
		gamePhase.CommandExit.Type.String()[8:],
		gamePhase.CommandExit.ForTeam.String()[5:],
		nextCommandType,
		nextCommandForTeam,
		gamePhase.CommandPrev.Type.String()[8:],
		gamePhase.CommandPrev.ForTeam.String()[5:],
		gamePhase.Stage.String()[6:],
		convertDuration(gamePhase.StageTimeLeftEntry),
		convertDuration(gamePhase.StageTimeLeftExit),
	)

	if err != nil {
		return err
	}

	if err := p.WriteGameEvents(gamePhase.GameEventsApplied, id, GameEventKindApplied); err != nil {
		return err
	}

	return p.WriteGameEvents(gamePhase.GameEventsProposed, id, GameEventKindProposed)
}
