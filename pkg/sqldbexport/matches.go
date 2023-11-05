package sqldbexport

import (
	"database/sql"
	"github.com/RoboCup-SSL/ssl-match-stats/pkg/matchstats"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"log"
	"strings"
)

func (p *SqlDbExporter) FindMatchId(logFileName string) *uuid.UUID {
	id := new(uuid.UUID)
	err := p.db.QueryRow(
		"SELECT id FROM matches WHERE file_name=$1",
		logFileName).Scan(id)
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		log.Print("Could not query matches:", err)
	}
	return id
}

func (p *SqlDbExporter) WriteMatches(matchStatsCollection *matchstats.MatchStatsCollection, tournamentId *uuid.UUID, division string) error {
	for _, matchStats := range matchStatsCollection.MatchStats {
		log.Println("Writing ", matchStats.Name)
		logFileName := strings.ReplaceAll(matchStats.Name, ".gz", "")
		matchId := p.FindMatchId(logFileName)
		if matchId == nil {
			matchId = new(uuid.UUID)
			*matchId = uuid.New()
		}

		if _, err := p.db.Exec(
			`INSERT INTO matches 
						(
							id, 
							file_name, 
							tournament_id_fk, 
							division,
							start_time,
							duration,
							extra_time,
							shootout,
							type
						) 
						VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
						ON CONFLICT ON CONSTRAINT matches_pkey DO UPDATE SET
							division=excluded.division,
							start_time=excluded.start_time,
							duration=excluded.duration,
							extra_time=excluded.extra_time,
							shootout=excluded.shootout,
							type=excluded.type`,
			matchId,
			logFileName,
			tournamentId,
			division,
			convertTime(matchStats.StartTime),
			convertDuration(matchStats.MatchDuration),
			matchStats.ExtraTime,
			matchStats.Shootout,
			matchStats.Type,
		); err != nil {
			return errors.Wrap(err, "Could not insert match")
		}

		if err := p.WriteTeamStats(matchStats, matchId); err != nil {
			return errors.Wrap(err, "Could not insert team stats")
		}

		if err := p.WriteGamePhases(matchStats, matchId); err != nil {
			return errors.Wrap(err, "Could not insert game phases")
		}
	}
	return nil
}
