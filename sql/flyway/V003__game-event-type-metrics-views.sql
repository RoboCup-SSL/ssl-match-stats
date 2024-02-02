drop view if exists game_events_per_match cascade;

create view game_events_per_match as
select gp.match_id_fk,
       ge.type,
       ge.by_team,
       ge.withdrawn,
       ge.proposed
from game_events ge
         join game_phases gp on ge.game_phase_id_fk = gp.id;

create view game_event_type_metrics as
select s.match_id_fk,
       m.tournament_name,
       m.division::text,
       m.file_name,
       gec.category::text as game_event_category,
       s.type::text       as game_event_type,
       s.total,
       s.accepted,
       s.accepted_withdrawn,
       s.proposed,
       s.proposed_withdrawn
from (select ge.match_id_fk,
             ge.type,
             count(*)                                               as total,
             count(*) FILTER (WHERE not withdrawn and not proposed) as accepted,
             count(*) FILTER (WHERE withdrawn and not proposed)     as accepted_withdrawn,
             count(*) FILTER (WHERE not withdrawn and proposed)     as proposed,
             count(*) FILTER (WHERE withdrawn and proposed)         as proposed_withdrawn
      from game_events_per_match ge
      group by ge.match_id_fk, ge.type) s
         join matches m on s.match_id_fk = m.id
         left join game_event_categories gec on s.type = gec.type
order by tournament_name, division, file_name, game_event_category, game_event_type;

create view game_event_type_metrics_per_team as
select team_name,
       s.match_id_fk,
       m.tournament_name,
       m.division::text,
       m.file_name,
       gec.category::text as game_event_category,
       s.type::text       as game_event_type,
       s.total,
       s.accepted,
       s.accepted_withdrawn,
       s.proposed,
       s.proposed_withdrawn
from (select ge.match_id_fk,
             ge.type,
             ge.by_team,
             count(*)                                               as total,
             count(*) FILTER (WHERE not withdrawn and not proposed) as accepted,
             count(*) FILTER (WHERE withdrawn and not proposed)     as accepted_withdrawn,
             count(*) FILTER (WHERE not withdrawn and proposed)     as proposed,
             count(*) FILTER (WHERE withdrawn and proposed)         as proposed_withdrawn
      from game_events_per_match ge
      group by ge.match_id_fk, ge.type, ge.by_team) s
         join matches m on s.match_id_fk = m.id
         join team_match_stats tms on s.match_id_fk = tms.match_id_fk
    and (s.by_team = tms.team_color or s.by_team = 'UNKNOWN')
         left join game_event_categories gec on s.type = gec.type
order by team_name, tournament_name, division, file_name, game_event_category, game_event_type;
