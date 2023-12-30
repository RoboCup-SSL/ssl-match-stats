create view game_phase_durations as
select gp.match_id_fk                  as match_id_fk,
       gp.type                         as type,
       m.file_name                     as file_name,
       m.tournament_name               as tournament_name,
       m.division::text                as division,
       gp.for_team                     as for_team,
       extract(epoch FROM gp.duration) as duration
from game_phases gp
         join matches m on m.id = gp.match_id_fk;

create view match_durations_by_type as
select gpd.file_name                                              as file_name,
       gpd.type::text                                             as game_phase_type,
       gpd.tournament_name                                        as tournament_name,
       gpd.division::text                                         as division,
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
group by gpd.file_name, gpd.type, gpd.tournament_name, gpd.division;

create view match_durations_by_type_per_team as
select tms.team_name                                              as team_name,
       gpd.type::text                                             as game_phase_type,
       gpd.tournament_name                                        as tournament_name,
       gpd.division::text                                         as division,
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
         join team_match_stats tms on tms.match_id_fk = gpd.match_id_fk
where gpd.for_team = 'NONE'
   or gpd.for_team = tms.team_color
group by tms.team_name, gpd.type, gpd.tournament_name, gpd.division;

create view match_durations as
select m.file_name                         as file_name,
       m.tournament_name                   as tournament_name,
       m.division::text                    as division,
       extract(epoch FROM m.duration)      as duration,
       extract(epoch FROM m.duration) / 60 as duration_minutes
from matches m;
