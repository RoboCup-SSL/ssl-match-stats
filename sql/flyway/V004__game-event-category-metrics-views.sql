create view game_event_category_metrics as
select m.id                                                   as match_id_fk,
       m.file_name                                            as file_name,
       ge.category::text                                      as game_event_category,
       m.tournament_name                                      as tournament_name,
       m.division::text                                       as division,
       count(*)                                               as total,
       count(*) FILTER (WHERE not withdrawn and not proposed) as accepted,
       count(*) FILTER (WHERE withdrawn and not proposed)     as accepted_withdrawn,
       count(*) FILTER (WHERE not withdrawn and proposed)     as proposed,
       count(*) FILTER (WHERE withdrawn and proposed)         as proposed_withdrawn
from game_events ge
         join game_phases gp on ge.game_phase_id_fk = gp.id
         join matches m on gp.match_id_fk = m.id
group by m.id, ge.category, m.tournament_name, m.division;

create view game_event_category_metrics_per_team as
select tms.team_name                                          as team_name,
       m.id                                                   as match_id_fk,
       m.file_name                                            as file_name,
       ge.category::text                                      as game_event_category,
       m.tournament_name                                      as tournament_name,
       m.division::text                                       as division,
       count(*)                                               as total,
       count(*) FILTER (WHERE not withdrawn and not proposed) as accepted,
       count(*) FILTER (WHERE withdrawn and not proposed)     as accepted_withdrawn,
       count(*) FILTER (WHERE not withdrawn and proposed)     as proposed,
       count(*) FILTER (WHERE withdrawn and proposed)         as proposed_withdrawn
from game_events ge
         join game_phases gp on ge.game_phase_id_fk = gp.id
         join matches m on gp.match_id_fk = m.id
         join team_match_stats tms on m.id = tms.match_id_fk
where gp.for_team = 'NONE'
   or gp.for_team = tms.team_color
group by m.id, tms.team_name, ge.category, m.tournament_name, m.division;
