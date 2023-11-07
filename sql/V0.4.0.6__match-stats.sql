
create or replace view game_phase_duration_per_type AS
select match_id_fk,
       type,
       sum(duration) as duration,
       count(*) as count
from game_phases
group by match_id_fk, type;

drop materialized view if exists match_stats;
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
