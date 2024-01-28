create view match_metrics as
select m.file_name             as file_name,
       m.tournament_name       as tournament_name,
       m.division              as division,
       tms.goals               as goals,
       tms.fouls               as fouls,
       tms.yellow_cards        as yellow_cards,
       tms.red_cards           as red_cards,
       tms.penalty_shots_total as penalty_shots,
       tms.timeouts_taken      as timeouts_taken,
       tms.timeout_time        as timeout_time,
       tms.ball_placements     as ball_placements,
       tms.ball_placement_time as ball_placement_time,
       gecm.accepted           as ball_outs
from matches m
         join team_match_stats tms on tms.match_id_fk = m.id
         left join game_event_category_metrics gecm on gecm.match_id_fk = m.id
    and gecm.game_event_category = 'BALL_OUT';

create view match_metrics_per_team as
select tms.team_name               as team_name,
       m.file_name                 as file_name,
       m.tournament_name           as tournament_name,
       m.division                  as division,
       tms.goals                   as goals,
       tms.conceded_goals          as conceded_goals,
       tms.fouls                   as fouls,
       tms.yellow_cards            as yellow_cards,
       tms.red_cards               as red_cards,
       tms.penalty_shots_total     as penalty_shots,
       tms.timeouts_taken          as timeouts_taken,
       tms.timeout_time            as timeout_time,
       tms.ball_placements         as ball_placements,
       tms.ball_placement_time     as ball_placement_time,
       tms.max_active_yellow_cards as max_active_yellow_cards,
       gecm.accepted               as ball_outs
from matches m
         join team_match_stats tms on tms.match_id_fk = m.id
         left join game_event_category_metrics_per_team gecm on gecm.team_name = tms.team_name
    and gecm.tournament_name = m.tournament_name
    and gecm.division = m.division::text
    and gecm.game_event_category = 'BALL_OUT';
