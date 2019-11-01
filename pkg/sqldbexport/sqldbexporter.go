package sqldbexport

import (
	"context"
	"database/sql"
	"errors"
	"github.com/RoboCup-SSL/ssl-match-stats/pkg/matchstats"
	"github.com/google/uuid"
	"log"

	_ "github.com/lib/pq"
)

type SqlDbExporter struct {
	db  *sql.DB
	ctx context.Context
}

func (p *SqlDbExporter) Connect(driver string, dataSourceName string) error {
	db, err := sql.Open(driver, dataSourceName)
	if err != nil {
		return err
	}
	p.db = db
	p.ctx = context.Background()
	return nil
}

func (p *SqlDbExporter) FindLogFileId(logFileName string) *uuid.UUID {
	id := new(uuid.UUID)
	err := p.db.QueryRow(
		"SELECT id FROM log_files WHERE file_name=$1",
		logFileName).Scan(id)
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		log.Print("Could not query log_files:", err)
	}
	return id
}

func (p *SqlDbExporter) WriteLogFiles(matchStatsCollection *matchstats.MatchStatsCollection) error {
	for _, matchStats := range matchStatsCollection.MatchStats {
		logFileName := matchStats.Name
		id := p.FindLogFileId(logFileName)
		if id == nil {
			id = new(uuid.UUID)
			*id = uuid.New()
			if _, err := p.db.Exec(
				"INSERT INTO log_files (id, file_name) VALUES ($1, $2)",
				id,
				logFileName,
			); err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *SqlDbExporter) WriteTeamStats(matchStatsCollection *matchstats.MatchStatsCollection) error {
	for _, matchStats := range matchStatsCollection.MatchStats {
		logFileName := matchStats.Name
		logFileId := p.FindLogFileId(logFileName)
		if logFileId == nil {
			return errors.New("Could not find log file in DB: " + logFileName)
		}
		if err := p.insertTeamStats(logFileId, matchStats.TeamStatsYellow, "yellow", matchStats.TeamStatsBlue.Name); err != nil {
			return err
		}
		if err := p.insertTeamStats(logFileId, matchStats.TeamStatsBlue, "blue", matchStats.TeamStatsYellow.Name); err != nil {
			return err
		}
	}
	return nil
}

func (p *SqlDbExporter) insertTeamStats(logFileId *uuid.UUID, teamStats *matchstats.TeamStats, teamColor string, opponent string) error {

	_, err := p.db.Exec(
		`INSERT INTO matches (
                     id, 
                     log_file_id_fk, 
                     team_color, 
                     team_name,
                     opponent_name,
                     goals,
                     goals_conceded,
                     fouls,
                     cards_yellow,
                     cards_red,
                     timeout_time,
                     timeouts_taken,
                     timeouts_left,
                     ball_placement_time,
                     ball_placements,
                     max_active_yellow_cards,
                     penalty_shots_total
                     ) 
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
				ON CONFLICT ON CONSTRAINT unique_log_file DO UPDATE SET
					  team_name=excluded.team_name,
					  opponent_name=excluded.opponent_name,
					  goals=excluded.goals,
					  goals_conceded=excluded.goals_conceded,
					  fouls=excluded.fouls,
					  cards_yellow=excluded.cards_yellow,
					  cards_red=excluded.cards_red,
					  timeout_time=excluded.timeout_time,
					  timeouts_taken=excluded.timeouts_taken,
					  timeouts_left=excluded.timeouts_left,
					  ball_placement_time=excluded.ball_placement_time,
					  ball_placements=excluded.ball_placements,
					  max_active_yellow_cards=excluded.max_active_yellow_cards,
					  penalty_shots_total=excluded.penalty_shots_total`,
		uuid.New(),
		logFileId,
		teamColor,
		teamStats.Name,
		opponent,
		teamStats.Goals,
		teamStats.ConcededGoals,
		teamStats.Fouls,
		teamStats.YellowCards,
		teamStats.RedCards,
		teamStats.TimeoutTime,
		teamStats.Timeouts,
		nil, // TODO
		teamStats.BallPlacementTime,
		teamStats.BallPlacements,
		teamStats.MaxActiveYellowCards,
		teamStats.PenaltyShotsTotal,
	)
	return err
}
