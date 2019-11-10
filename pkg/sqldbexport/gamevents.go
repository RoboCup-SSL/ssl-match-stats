package sqldbexport

import (
	"database/sql"
	"encoding/json"
	"github.com/RoboCup-SSL/ssl-match-stats/pkg/matchstats"
	"github.com/google/uuid"
	"log"
)

func (p *SqlDbExporter) WriteGameEvents(gameEvents []*matchstats.GameEventTimed, gamePhaseId uuid.UUID) error {
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

	// TODO proposed game events

	gameEventId := uuid.New()
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
		gameEventId,
		gamePhaseId,
		gameEvent.GameEvent.Type.String(),
		convertTime(gameEvent.Timestamp),
		gameEvent.Withdrawn,
		payload,
	)

	for _, origin := range gameEvent.GameEvent.Origin {
		var autoRefId uuid.UUID
		if id, ok := p.autoRefs[origin]; ok {
			autoRefId = id
		} else if presentAutoRefId := p.FindAutoRefId(origin); presentAutoRefId != nil {
			autoRefId = *presentAutoRefId
		} else {
			id, err := p.insertAutoRefId(origin)
			if err != nil {
				return err
			}
			autoRefId = id
		}
		p.autoRefs[origin] = autoRefId
		if err := p.insertAutoRefToGameEventMapping(autoRefId, gameEventId); err != nil {
			return err
		}
	}

	return err
}

func (p *SqlDbExporter) FindAutoRefId(autoRefName string) *uuid.UUID {
	id := new(uuid.UUID)
	err := p.db.QueryRow(
		"SELECT id FROM auto_refs WHERE name=$1",
		autoRefName).Scan(id)
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		log.Print("Could not query autoRefs:", err)
	}
	return id
}

func (p *SqlDbExporter) insertAutoRefId(autoRefName string) (uuid.UUID, error) {
	autoRefId := uuid.New()
	_, err := p.db.Exec(
		`INSERT INTO auto_refs (
                     id, 
                     name 
                     ) 
				VALUES ($1, $2)`,
		autoRefId,
		autoRefName,
	)

	return autoRefId, err
}

func (p *SqlDbExporter) insertAutoRefToGameEventMapping(autoRefId, gameEventId uuid.UUID) error {
	_, err := p.db.Exec(
		`INSERT INTO game_event_auto_ref_mapping (
                     auto_ref_id_fk, 
                     game_event_id_fk 
                     ) 
				VALUES ($1, $2)`,
		autoRefId,
		gameEventId,
	)

	return err
}
