create view game_phase_durations as
select gp.match_id_fk,
       gp.type,
       gp.for_team,
       extract(epoch FROM gp.duration) as duration
from game_phases gp;

create view match_durations_by_type AS
select m.tournament_name,
       m.division::text,
       m.file_name,
       gpd.type::text,
       count(*)                                                   as count,
       sum(gpd.duration)                                          as total,
       sum(gpd.duration) / 60                                     as total_minutes,
       min(gpd.duration)                                          as min,
       avg(gpd.duration)                                          as avg,
       max(gpd.duration)                                          as max,
       percentile_disc(0.25) within group (order by gpd.duration) as p25,
       percentile_disc(0.5) within group (order by gpd.duration)  as median,
       percentile_disc(0.75) within group (order by gpd.duration) as p75
from game_phase_durations gpd
         join matches m on m.id = gpd.match_id_fk
group by m.tournament_name, m.division, m.file_name, gpd.type;

create view match_durations_by_type_per_team AS
select m.tournament_name,
       m.division::text,
       tms.team_name,
       gpd.type::text,
       count(*)                                                   as count,
       sum(gpd.duration)                                          as total,
       sum(gpd.duration) / 60                                     as total_minutes,
       min(gpd.duration)                                          as min,
       avg(gpd.duration)                                          as avg,
       max(gpd.duration)                                          as max,
       percentile_disc(0.25) within group (order by gpd.duration) as p25,
       percentile_disc(0.5) within group (order by gpd.duration)  as median,
       percentile_disc(0.75) within group (order by gpd.duration) as p75
from game_phase_durations gpd
         join matches m on m.id = gpd.match_id_fk
         join team_match_stats tms on tms.match_id_fk = gpd.match_id_fk
where gpd.for_team = 'NONE'
   or gpd.for_team = tms.team_color
group by m.tournament_name, m.division, tms.team_name, gpd.type;

create view match_durations as
select m.tournament_name,
       m.division::text,
       m.file_name,
       extract(epoch FROM m.duration)      as duration,
       extract(epoch FROM m.duration) / 60 as duration_minutes
from matches m;
