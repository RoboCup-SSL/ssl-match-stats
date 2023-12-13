-- the view can be used with Metabase to great a histogram of game phase durations
create view game_phase_duration as
select id, type, extract(epoch FROM duration) duration_seconds
from game_phases;

create view game_phase_duration_per_match AS
select gp.match_id_fk,
       gp.type,
       count(*)                                                  as count,
       sum(gp.duration)                                          as total,
       min(gp.duration)                                          as min,
       avg(gp.duration)                                          as avg,
       max(gp.duration)                                          as max,
       percentile_disc(0.25) within group (order by gp.duration) as p25,
       percentile_disc(0.5) within group (order by gp.duration)  as median,
       percentile_disc(0.75) within group (order by gp.duration) as p75
from game_phases gp
group by gp.type, gp.match_id_fk;

create view game_phase_duration_per_team AS
select m.team_name,
       gp.type,
       count(*)                                                  as count,
       sum(gp.duration)                                          as total,
       min(gp.duration)                                          as min,
       avg(gp.duration)                                          as avg,
       max(gp.duration)                                          as max,
       percentile_disc(0.25) within group (order by gp.duration) as p25,
       percentile_disc(0.5) within group (order by gp.duration)  as median,
       percentile_disc(0.75) within group (order by gp.duration) as p75
from game_phases gp
         join team_match_stats m on m.match_id_fk = gp.match_id_fk
group by m.team_name, gp.type;

create view game_phase_duration_per_tournament AS
select m.tournament_name,
       m.division,
       gp.type,
       count(*)                                                  as count,
       sum(gp.duration)                                          as total,
       min(gp.duration)                                          as min,
       avg(gp.duration)                                          as avg,
       max(gp.duration)                                          as max,
       percentile_disc(0.25) within group (order by gp.duration) as p25,
       percentile_disc(0.5) within group (order by gp.duration)  as median,
       percentile_disc(0.75) within group (order by gp.duration) as p75
from game_phases gp
         join matches m on m.id = gp.match_id_fk
group by m.tournament_name, m.division, gp.type;

create view game_phase_duration_overall AS
select gp.type,
       count(*)                                                  as count,
       sum(gp.duration)                                          as total,
       min(gp.duration)                                          as min,
       avg(gp.duration)                                          as avg,
       max(gp.duration)                                          as max,
       percentile_disc(0.25) within group (order by gp.duration) as p25,
       percentile_disc(0.5) within group (order by gp.duration)  as median,
       percentile_disc(0.75) within group (order by gp.duration) as p75
from game_phases gp
         join team_match_stats m on m.match_id_fk = gp.match_id_fk
group by gp.type;
