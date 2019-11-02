package sqldbexport

import (
	"database/sql"
	"github.com/RoboCup-SSL/ssl-match-stats/pkg/matchstats"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"log"
	"time"
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
		logFileName := matchStats.Name
		matchId := p.FindMatchId(logFileName)
		if matchId == nil {
			matchId = new(uuid.UUID)
			*matchId = uuid.New()
		}
		startTimeSec := int64(matchStats.StartTime / 1000000)
		startTimeNsec := (int64(matchStats.StartTime) - (startTimeSec * 1000000)) * 1000
		startTime := time.Unix(startTimeSec, startTimeNsec)

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
							shootout
						) 
						VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
						ON CONFLICT ON CONSTRAINT matches_pkey DO UPDATE SET
							division=excluded.division,
							start_time=excluded.start_time,
							duration=excluded.duration,
							extra_time=excluded.extra_time,
							shootout=excluded.shootout`,
			matchId,
			logFileName,
			tournamentId,
			division,
			startTime,
			time.Duration(matchStats.MatchDuration*1000).Seconds()*1000,
			matchStats.ExtraTime,
			matchStats.Shootout,
		); err != nil {
			return errors.Wrap(err, "Could not insert match")
		}

		if err := p.WriteTeamStats(matchStats, matchId); err != nil {
			return errors.Wrap(err, "Could not insert team stats")
		}
	}
	return nil
}

func (p *SqlDbExporter) WriteTeamStats(matchStats *matchstats.MatchStats, matchId *uuid.UUID) error {
	if err := p.insertTeamStats(matchId, matchStats.TeamStatsYellow, "yellow", matchStats.TeamStatsBlue.Name); err != nil {
		return err
	}
	if err := p.insertTeamStats(matchId, matchStats.TeamStatsBlue, "blue", matchStats.TeamStatsYellow.Name); err != nil {
		return err
	}
	return nil
}

func (p *SqlDbExporter) insertTeamStats(matchId *uuid.UUID, teamStats *matchstats.TeamStats, teamColor string, opponent string) error {

	_, err := p.db.Exec(
		`INSERT INTO team_match_stats (
                     id, 
                     match_id_fk, 
                     team_color, 
                     team_name,
                     opponent_name,
                     goals,
                     conceded_goals,
                     fouls,
                     yellow_cards,
                     red_cards,
                     timeout_time,
                     timeouts_taken,
                     timeouts_left,
                     ball_placement_time,
                     ball_placements,
                     max_active_yellow_cards,
                     penalty_shots_total
                     ) 
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
				ON CONFLICT ON CONSTRAINT unique_match DO UPDATE SET
					  team_name=excluded.team_name,
					  opponent_name=excluded.opponent_name,
					  goals=excluded.goals,
					  conceded_goals=excluded.conceded_goals,
					  fouls=excluded.fouls,
					  yellow_cards=excluded.yellow_cards,
					  red_cards=excluded.red_cards,
					  timeout_time=excluded.timeout_time,
					  timeouts_taken=excluded.timeouts_taken,
					  timeouts_left=excluded.timeouts_left,
					  ball_placement_time=excluded.ball_placement_time,
					  ball_placements=excluded.ball_placements,
					  max_active_yellow_cards=excluded.max_active_yellow_cards,
					  penalty_shots_total=excluded.penalty_shots_total`,
		uuid.New(),
		matchId,
		teamColor,
		teamStats.Name,
		opponent,
		teamStats.Goals,
		teamStats.ConcededGoals,
		teamStats.Fouls,
		teamStats.YellowCards,
		teamStats.RedCards,
		time.Duration(teamStats.TimeoutTime*1000).Seconds()*1000,
		teamStats.TimeoutsTaken,
		teamStats.TimeoutsLeft,
		time.Duration(teamStats.BallPlacementTime*1000).Seconds()*1000,
		teamStats.BallPlacements,
		teamStats.MaxActiveYellowCards,
		teamStats.PenaltyShotsTotal,
	)
	return err
}
