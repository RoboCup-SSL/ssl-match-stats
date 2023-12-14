create view team_metrics_per_tournament_total AS
select mm.team_name,
       m.tournament_name,
       m.division,
       sum(mm.goals)                   as goals,
       sum(mm.conceded_goals)          as conceded_goals,
       sum(mm.fouls)                   as fouls,
       sum(mm.yellow_cards)            as yellow_cards,
       sum(mm.red_cards)               as red_cards,
       sum(mm.penalty_shots_total)     as penalty_shots,
       sum(mm.timeouts_taken)          as timeouts_taken,
       sum(mm.timeout_time)            as timeout_time,
       sum(mm.ball_placements)         as ball_placements,
       sum(mm.ball_placement_time)     as ball_placement_time,
       sum(mm.max_active_yellow_cards) as max_active_yellow_cards,
       sum(gecm.accepted)              as ball_out
from matches m
         join team_match_stats mm on mm.match_id_fk = m.id
         left join game_event_category_metrics gecm on gecm.team_name = mm.team_name
    and gecm.tournament_name = m.tournament_name
    and gecm.division = m.division
    and gecm.category = 'BALL_OUT'
group by mm.team_name, m.tournament_name, m.division
order by m.tournament_name, m.division;

create view team_metrics_per_tournament_avg AS
select mm.team_name,
       m.tournament_name,
       m.division,
       avg(mm.goals)                   as goals,
       avg(mm.conceded_goals)          as conceded_goals,
       avg(mm.fouls)                   as fouls,
       avg(mm.yellow_cards)            as yellow_cards,
       avg(mm.red_cards)               as red_cards,
       avg(mm.penalty_shots_total)     as penalty_shots,
       avg(mm.timeouts_taken)          as timeouts_taken,
       avg(mm.timeout_time)            as timeout_time,
       avg(mm.ball_placements)         as ball_placements,
       avg(mm.ball_placement_time)     as ball_placement_time,
       avg(mm.max_active_yellow_cards) as max_active_yellow_cards,
       avg(gecm.accepted)              as ball_out
from matches m
         join team_match_stats mm on mm.match_id_fk = m.id
         left join game_event_category_metrics gecm on gecm.team_name = mm.team_name
    and gecm.tournament_name = m.tournament_name
    and gecm.division = m.division
    and gecm.category = 'BALL_OUT'
group by mm.team_name, m.tournament_name, m.division
order by m.tournament_name, m.division;

create view team_metrics_overall_total AS
select mm.team_name,
       sum(mm.goals)               as goals,
       sum(mm.fouls)               as fouls,
       sum(mm.yellow_cards)        as yellow_cards,
       sum(mm.red_cards)           as red_cards,
       sum(mm.penalty_shots_total) as penalty_shots,
       sum(mm.timeouts_taken)      as timeouts_taken,
       sum(mm.timeout_time)        as timeout_time,
       sum(mm.ball_placements)     as ball_placements,
       sum(mm.ball_placement_time) as ball_placement_time,
       sum(gecm.accepted)          as ball_out
from matches m
         join team_match_stats mm on mm.match_id_fk = m.id
         left join game_event_category_metrics gecm on gecm.team_name = mm.team_name
    and gecm.category = 'BALL_OUT'
group by mm.team_name;

create view team_metrics_overall_avg AS
select mm.team_name,
       avg(mm.goals)               as goals,
       avg(mm.fouls)               as fouls,
       avg(mm.yellow_cards)        as yellow_cards,
       avg(mm.red_cards)           as red_cards,
       avg(mm.penalty_shots_total) as penalty_shots,
       avg(mm.timeouts_taken)      as timeouts_taken,
       avg(mm.timeout_time)        as timeout_time,
       avg(mm.ball_placements)     as ball_placements,
       avg(mm.ball_placement_time) as ball_placement_time,
       avg(gecm.accepted)          as ball_out
from matches m
         join team_match_stats mm on mm.match_id_fk = m.id
         left join game_event_category_metrics gecm on gecm.team_name = mm.team_name
    and gecm.category = 'BALL_OUT'
group by mm.team_name;
