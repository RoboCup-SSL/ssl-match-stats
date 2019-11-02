package sqldbexport

import (
	"github.com/RoboCup-SSL/ssl-match-stats/pkg/matchstats"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (p *SqlDbExporter) WriteTeamStats(matchStats *matchstats.MatchStats, matchId *uuid.UUID) error {
	if _, err := p.db.Exec("DELETE FROM team_match_stats WHERE match_id_fk=$1", matchId); err != nil {
		return errors.Wrap(err, "Could not delete previous match stats for id:"+matchId.String())
	}

	if err := p.insertTeamStats(matchId, matchStats.TeamStatsYellow, "YELLOW", matchStats.TeamStatsBlue.Name); err != nil {
		return err
	}
	if err := p.insertTeamStats(matchId, matchStats.TeamStatsBlue, "BLUE", matchStats.TeamStatsYellow.Name); err != nil {
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
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)`,
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
		convertDuration(teamStats.TimeoutTime),
		teamStats.TimeoutsTaken,
		teamStats.TimeoutsLeft,
		convertDuration(teamStats.BallPlacementTime),
		teamStats.BallPlacements,
		teamStats.MaxActiveYellowCards,
		teamStats.PenaltyShotsTotal,
	)
	return err
}
