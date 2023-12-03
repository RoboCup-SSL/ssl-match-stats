create view game_phase_duration_per_type AS
select match_id_fk,
       type,
       sum(duration) as duration,
       count(*) as count
from game_phases
group by match_id_fk, type;

create materialized view match_stats as
select m.id as match_id_fk,
       m.duration,
       (select duration
        from game_phase_duration_per_type pt
        where pt.match_id_fk = m.id and pt.type = 'RUNNING') as duration_running,
       (select duration
        from game_phase_duration_per_type pt
        where pt.match_id_fk = m.id and pt.type = 'STOP') as duration_stop,
       (select duration
        from game_phase_duration_per_type pt
        where pt.match_id_fk = m.id and pt.type = 'HALT') as duration_halt,
       (select duration
        from game_phase_duration_per_type pt
        where pt.match_id_fk = m.id and pt.type = 'BALL_PLACEMENT') as duration_ball_placement,
       (select duration
        from game_phase_duration_per_type pt
        where pt.match_id_fk = m.id and pt.type = 'TIMEOUT') as duration_timeout,
       (select sum(duration)
        from game_phase_duration_per_type pt
        where pt.match_id_fk = m.id and pt.type not in ('RUNNING', 'STOP', 'HALT', 'BALL_PLACEMENT', 'TIMEOUT') group by id) as duration_other,
       (select count
        from game_phase_duration_per_type pt
        where pt.match_id_fk = m.id and pt.type = 'RUNNING') as count_running,
       (select count
        from game_phase_duration_per_type pt
        where pt.match_id_fk = m.id and pt.type = 'STOP') as count_stop,
       (select count
        from game_phase_duration_per_type pt
        where pt.match_id_fk = m.id and pt.type = 'HALT') as count_halt,
       (select count
        from game_phase_duration_per_type pt
        where pt.match_id_fk = m.id and pt.type = 'BALL_PLACEMENT') as count_ball_placement,
       (select count
        from game_phase_duration_per_type pt
        where pt.match_id_fk = m.id and pt.type = 'TIMEOUT') as count_timeout,
       (select sum(goals)
        from team_match_stats pt
        where pt.match_id_fk = m.id) as goals,
       (select sum(fouls)
        from team_match_stats pt
        where pt.match_id_fk = m.id) as fouls,
       (select sum(yellow_cards)
        from team_match_stats pt
        where pt.match_id_fk = m.id) as yellow_cards,
       (select sum(red_cards)
        from team_match_stats pt
        where pt.match_id_fk = m.id) as red_cards,
       (select sum(timeouts_taken)
        from team_match_stats pt
        where pt.match_id_fk = m.id) as timeouts_taken,
       (select sum(timeout_time)
        from team_match_stats pt
        where pt.match_id_fk = m.id) as timeout_time,
       (select sum(ball_placement_time)
        from team_match_stats pt
        where pt.match_id_fk = m.id) as ball_placement_time,
       (select sum(ball_placements)
        from team_match_stats pt
        where pt.match_id_fk = m.id) as ball_placements,
       (select sum(penalty_shots_total)
        from team_match_stats pt
        where pt.match_id_fk = m.id) as penalty_shots_total
from matches m;

drop view tournament_game_phase_duration;
create view tournament_game_phase_duration as
select tournament_name,
       division,
       s.type,
       min(s.duration)                                                                                          as min,
       avg(s.duration)                                                                                          as avg,
       max(s.duration)                                                                                          as max,
       make_interval(secs := EXTRACT(epoch FROM percentile_disc(0.90) within group (order by s.duration desc))) as p90
from game_phase_duration_per_type s
         left join matches on matches.id = match_id_fk
group by tournament_name, division, s.type;

create or replace view tournament_stats as
select matches.tournament_name,
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
group by matches.tournament_name, matches.division;

create view tournament_stats_avg_match as
select tournament_name,
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

-- the view can be used with Metabase to great a histogram of game phase durations
create view game_phase_duration as
select id, type, extract(epoch FROM duration) duration_seconds from game_phases;
