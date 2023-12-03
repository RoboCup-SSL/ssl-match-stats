drop view tournament_game_phase_duration;
create view tournament_game_phase_duration as
select tournament_id_fk,
       division,
       s.type,
       min(s.duration)                                                                                          as min,
       avg(s.duration)                                                                                          as avg,
       max(s.duration)                                                                                          as max,
       make_interval(secs := EXTRACT(epoch FROM percentile_disc(0.90) within group (order by s.duration desc))) as p90
from game_phase_duration_per_type s
         left join matches on matches.id = match_id_fk
group by tournament_id_fk, division, s.type;

create or replace view tournament_stats as
select matches.tournament_id_fk,
       matches.division,
       count(distinct team_name)  as num_teams,
       count(distinct matches.id) as num_matches,
       sum(s.goals)               as goals,
       sum(s.fouls)               as fouls,
       sum(s.yellow_cards)        as yellow_cards,
       sum(s.red_cards)           as red_cards,
       sum(s.timeouts_taken)      as timeouts_taken,
       sum(s.timeout_time)        as timeout_time,
       sum(s.ball_placement_time) as ball_placement_time,
       sum(s.ball_placements)     as ball_placements,
       sum(s.penalty_shots_total) as penalty_shots_total
from team_match_stats s
         join matches on matches.id = s.match_id_fk
group by matches.tournament_id_fk, matches.division;

drop view if exists tournament_stats_avg_match;
create view tournament_stats_avg_match as
select tournament_id_fk,
       division,
       num_teams,
       num_matches,
       goals::decimal / num_matches               as goals,
       fouls::decimal / num_matches               as fouls,
       yellow_cards::decimal / num_matches        as yellow_cards,
       red_cards::decimal / num_matches           as red_cards,
       timeouts_taken::decimal / num_matches      as timeouts_taken,
       timeout_time / num_matches                 as timeout_time,
       ball_placement_time / num_matches          as ball_placement_time,
       ball_placements::decimal / num_matches     as ball_placements,
       penalty_shots_total::decimal / num_matches as penalty_shots_total
from tournament_stats s;
