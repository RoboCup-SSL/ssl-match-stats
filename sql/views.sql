drop view if exists game_phase_duration_per_match;
create view game_phase_duration_per_match AS
select gp.match_id_fk,
       gp.type,
       count(*) as count,
       sum(gp.duration) as total,
       min(gp.duration) as min,
       avg(gp.duration) as avg,
       max(gp.duration) as max,
       percentile_disc(0.25) within group (order by gp.duration) as p25,
       percentile_disc(0.5) within group (order by gp.duration) as median,
       percentile_disc(0.75) within group (order by gp.duration) as p75
from game_phases gp
group by gp.type, gp.match_id_fk;

drop view if exists game_phase_duration_per_tournament;
create view game_phase_duration_per_tournament AS
select m.tournament_name,
       m.division,
       gp.type,
       count(*) as count,
       sum(gp.duration) as total,
       min(gp.duration) as min,
       avg(gp.duration) as avg,
       max(gp.duration) as max,
       percentile_disc(0.25) within group (order by gp.duration) as p25,
       percentile_disc(0.5) within group (order by gp.duration) as median,
       percentile_disc(0.75) within group (order by gp.duration) as p75
from game_phases gp
join matches m on m.id = gp.match_id_fk
group by m.tournament_name, m.division, gp.type;

drop view if exists game_phase_duration_per_team;
create view game_phase_duration_per_team AS
select m.team_name,
       gp.type,
       count(*) as count,
       sum(gp.duration) as total,
       min(gp.duration) as min,
       avg(gp.duration) as avg,
       max(gp.duration) as max,
       percentile_disc(0.25) within group (order by gp.duration) as p25,
       percentile_disc(0.5) within group (order by gp.duration) as median,
       percentile_disc(0.75) within group (order by gp.duration) as p75
from game_phases gp
         join team_match_stats m on m.match_id_fk = gp.match_id_fk
group by m.team_name, gp.type;

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
