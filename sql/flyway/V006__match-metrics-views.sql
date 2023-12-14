create view match_metrics AS
select m.id                         as match_id_fk,
       sum(tms.goals)               as goals,
       sum(tms.fouls)               as fouls,
       sum(tms.yellow_cards)        as yellow_cards,
       sum(tms.red_cards)           as red_cards,
       sum(tms.penalty_shots_total) as penalty_shots,
       sum(tms.timeouts_taken)      as timeouts_taken,
       sum(tms.timeout_time)        as timeout_time,
       sum(tms.ball_placements)     as ball_placements,
       sum(tms.ball_placement_time) as ball_placement_time,
       sum(gecm.accepted)           as ball_out
from matches m
         join team_match_stats tms on tms.match_id_fk = m.id
         left join game_event_category_metrics gecm on gecm.match_id_fk = m.id
    and gecm.category = 'BALL_OUT'
group by m.id;

create view match_metrics_per_tournament_total AS
select m.tournament_name,
       m.division,
       sum(mm.goals)               as goals,
       sum(mm.fouls)               as fouls,
       sum(mm.yellow_cards)        as yellow_cards,
       sum(mm.red_cards)           as red_cards,
       sum(mm.penalty_shots)       as penalty_shots,
       sum(mm.timeouts_taken)      as timeouts_taken,
       sum(mm.timeout_time)        as timeout_time,
       sum(mm.ball_placements)     as ball_placements,
       sum(mm.ball_placement_time) as ball_placement_time,
       sum(gecm.accepted)          as ball_out
from matches m
         join match_metrics mm on mm.match_id_fk = m.id
         left join game_event_category_metrics gecm on gecm.tournament_name = m.tournament_name
    and gecm.division = m.division
    and gecm.category = 'BALL_OUT'
group by m.tournament_name, m.division
order by m.tournament_name, m.division;

create view match_metrics_per_tournament_avg AS
select m.tournament_name,
       m.division,
       avg(mm.goals)               as goals,
       avg(mm.fouls)               as fouls,
       avg(mm.yellow_cards)        as yellow_cards,
       avg(mm.red_cards)           as red_cards,
       avg(mm.penalty_shots)       as penalty_shots,
       avg(mm.timeouts_taken)      as timeouts_taken,
       avg(mm.timeout_time)        as timeout_time,
       avg(mm.ball_placements)     as ball_placements,
       avg(mm.ball_placement_time) as ball_placement_time,
       avg(gecm.accepted)          as ball_out
from matches m
         join match_metrics mm on mm.match_id_fk = m.id
         left join game_event_category_metrics gecm on gecm.tournament_name = m.tournament_name
    and gecm.division = m.division
    and gecm.category = 'BALL_OUT'
group by m.tournament_name, m.division
order by m.tournament_name, m.division;

create view match_metrics_overall_total AS
select sum(mm.goals)               as goals,
       sum(mm.fouls)               as fouls,
       sum(mm.yellow_cards)        as yellow_cards,
       sum(mm.red_cards)           as red_cards,
       sum(mm.penalty_shots)       as penalty_shots,
       sum(mm.timeouts_taken)      as timeouts_taken,
       sum(mm.timeout_time)        as timeout_time,
       sum(mm.ball_placements)     as ball_placements,
       sum(mm.ball_placement_time) as ball_placement_time,
       sum(gecm.accepted)          as ball_out
from matches m
         join match_metrics mm on mm.match_id_fk = m.id
         left join game_event_category_metrics gecm on gecm.category = 'BALL_OUT';

create view match_metrics_overall_avg AS
select avg(mm.goals)               as goals,
       avg(mm.fouls)               as fouls,
       avg(mm.yellow_cards)        as yellow_cards,
       avg(mm.red_cards)           as red_cards,
       avg(mm.penalty_shots)       as penalty_shots,
       avg(mm.timeouts_taken)      as timeouts_taken,
       avg(mm.timeout_time)        as timeout_time,
       avg(mm.ball_placements)     as ball_placements,
       avg(mm.ball_placement_time) as ball_placement_time,
       sum(gecm.accepted)          as ball_out
from matches m
         join match_metrics mm on mm.match_id_fk = m.id
         left join game_event_category_metrics gecm on gecm.category = 'BALL_OUT';
