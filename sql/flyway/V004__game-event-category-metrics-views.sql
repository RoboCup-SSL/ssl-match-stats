drop view if exists game_event_category_metrics cascade;
drop view if exists game_event_category_metrics_per_team cascade;

create view game_event_category_metrics as
select s.match_id_fk,
       m.tournament_name,
       m.division::text,
       m.file_name,
       s.category::text as game_event_category,
       s.total,
       s.accepted,
       s.accepted_withdrawn,
       s.proposed,
       s.proposed_withdrawn
from (select ge.match_id_fk,
             gec.category,
             count(*)                                               as total,
             count(*) FILTER (WHERE not withdrawn and not proposed) as accepted,
             count(*) FILTER (WHERE withdrawn and not proposed)     as accepted_withdrawn,
             count(*) FILTER (WHERE not withdrawn and proposed)     as proposed,
             count(*) FILTER (WHERE withdrawn and proposed)         as proposed_withdrawn
      from game_events_per_match ge
               left join game_event_categories gec on ge.type = gec.type
      group by ge.match_id_fk, gec.category) s
         join matches m on s.match_id_fk = m.id
order by tournament_name, division, file_name, game_event_category;

create view game_event_category_metrics_per_team as
select team_name,
       s.match_id_fk,
       m.tournament_name,
       m.division::text,
       m.file_name,
       s.category::text as game_event_category,
       s.total,
       s.accepted,
       s.accepted_withdrawn,
       s.proposed,
       s.proposed_withdrawn
from (select ge.match_id_fk,
             gec.category,
             ge.by_team,
             count(*)                                               as total,
             count(*) FILTER (WHERE not withdrawn and not proposed) as accepted,
             count(*) FILTER (WHERE withdrawn and not proposed)     as accepted_withdrawn,
             count(*) FILTER (WHERE not withdrawn and proposed)     as proposed,
             count(*) FILTER (WHERE withdrawn and proposed)         as proposed_withdrawn
      from game_events_per_match ge
               left join game_event_categories gec on ge.type = gec.type
      group by ge.match_id_fk, gec.category, ge.by_team) s
         join matches m on s.match_id_fk = m.id
         join team_match_stats tms on s.match_id_fk = tms.match_id_fk
    and (s.by_team = tms.team_color or s.by_team = 'UNKNOWN')
order by team_name, tournament_name, division, file_name, game_event_category;
