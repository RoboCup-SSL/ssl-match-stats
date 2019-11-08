package sqldbexport

import (
	"encoding/json"
	"github.com/RoboCup-SSL/ssl-match-stats/pkg/matchstats"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (p *SqlDbExporter) WriteGameEvents(gameEvents []*matchstats.GameEventTimed, gamePhaseId uuid.UUID) error {
	if _, err := p.db.Exec("DELETE FROM game_events WHERE game_phase_id_fk=$1", gamePhaseId); err != nil {
		return errors.Wrap(err, "Could not delete previous game phases for id:"+gamePhaseId.String())
	}

	for _, gameEvent := range gameEvents {
		if err := p.insertGameEvent(gameEvent, gamePhaseId); err != nil {
			return err
		}
	}
	return nil
}

func (p *SqlDbExporter) insertGameEvent(gameEvent *matchstats.GameEventTimed, gamePhaseId uuid.UUID) error {

	payload, err := json.Marshal(gameEvent.GameEvent)
	if err != nil {
		return err
	}

	_, err = p.db.Exec(
		`INSERT INTO game_events (
                     id, 
                     game_phase_id_fk, 
                     type,
					 timestamp,
					 withdrawn,
					 payload
                     ) 
				VALUES ($1, $2, $3, $4, $5, $6)`,
		uuid.New(),
		gamePhaseId,
		gameEvent.GameEvent.Type.String(),
		convertTime(gameEvent.Timestamp),
		gameEvent.Withdrawn,
		payload,
	)
	return err
}
