package sqldbexport

import (
	"encoding/json"
	"github.com/RoboCup-SSL/ssl-match-stats/internal/referee"
	"github.com/RoboCup-SSL/ssl-match-stats/internal/sslcommon"
	"github.com/RoboCup-SSL/ssl-match-stats/pkg/matchstats"
	"github.com/google/uuid"
	"log"
	"reflect"
	"time"
)

type GameEventKind int

const (
	GameEventKindApplied  = iota
	GameEventKindProposed = iota
)

func (p *SqlDbExporter) WriteGameEvents(gameEvents []*matchstats.GameEventTimed, gamePhaseId uuid.UUID, kind GameEventKind) error {
	for _, gameEvent := range gameEvents {
		if err := p.insertGameEvent(gameEvent, gamePhaseId, kind); err != nil {
			return err
		}
	}
	return nil
}

func (p *SqlDbExporter) insertGameEvent(gameEvent *matchstats.GameEventTimed, gamePhaseId uuid.UUID, kind GameEventKind) error {

	payload, err := json.Marshal(gameEvent.GameEvent)
	if err != nil {
		return err
	}

	byTeam := ByTeam(gameEvent.GameEvent)

	var createdTimestamp *time.Time
	if gameEvent.GameEvent.CreatedTimestamp != nil {
		createdTimestamp = new(time.Time)
		*createdTimestamp = convertTime(*gameEvent.GameEvent.CreatedTimestamp)
	}

	if gameEvent.GameEvent.Id == nil {
		log.Fatal("Game event id is nil: ", gameEvent)
	}

	gameEventId := *gameEvent.GameEvent.Id
	result, err := p.db.Exec(
		`INSERT INTO game_events (
                     id, 
                     game_phase_id_fk, 
                     type,
					 category,
					 by_team,
					 timestamp,
					 created_timestamp,
					 withdrawn,
					 proposed,
					 payload
                     ) 
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
				ON CONFLICT DO NOTHING`,
		gameEventId,
		gamePhaseId,
		gameEvent.GameEvent.Type.String(),
		gameEvent.Category.String()[9:],
		byTeam.String()[5:],
		convertTime(gameEvent.Timestamp),
		createdTimestamp,
		gameEvent.Withdrawn,
		kind == GameEventKindProposed,
		payload,
	)

	if err != nil {
		return err
	}

	if rows, err := result.RowsAffected(); err != nil && rows == 1 {
		for _, origin := range gameEvent.GameEvent.Origin {
			if err := p.insertOriginToGameEventMapping(gameEventId, origin); err != nil {
				return err
			}
		}
	}

	return err
}

func (p *SqlDbExporter) insertOriginToGameEventMapping(gameEventId string, origin string) error {
	_, err := p.db.Exec(
		`INSERT INTO game_event_origin_mapping (
			     game_event_id_fk,
				 game_event_origin
            	 ) 
			   VALUES ($1, $2)`,
		gameEventId,
		origin,
	)

	return err
}

// ByTeam extracts the `ByTeam` attribute from the game event details, if present
func ByTeam(gameEvent *referee.GameEvent) matchstats.TeamColor {
	if gameEvent.GetEvent() == nil {
		return matchstats.TeamColor_TEAM_UNKNOWN
	}
	event := reflect.ValueOf(gameEvent.GetEvent())
	if event.Elem().NumField() == 0 {
		return matchstats.TeamColor_TEAM_UNKNOWN
	}
	// all structs have a single field that we need to access
	v := event.Elem().Field(0)
	if !v.IsNil() {
		byTeamValue := v.Elem().FieldByName("ByTeam")
		if byTeamValue.IsValid() && !byTeamValue.IsNil() {
			byTeam := sslcommon.Team(byTeamValue.Elem().Int())
			switch byTeam {
			case sslcommon.Team_YELLOW:
				return matchstats.TeamColor_TEAM_YELLOW
			case sslcommon.Team_BLUE:
				return matchstats.TeamColor_TEAM_BLUE
			}
		}
	}
	return matchstats.TeamColor_TEAM_UNKNOWN
}
