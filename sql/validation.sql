-- game phase durations should not be negative
select *
from game_phases gp
         join matches m on m.id = gp.match_id_fk
where gp.duration < '0s'::interval;

-- there should be no very long game phases
select *
from game_phases gp
         join matches m on m.id = gp.match_id_fk
where gp.duration > '1h'::interval;

-- sum of game phase duration should the same as match duration
select matches.file_name, duration, duration_sum, duration - duration_sum as diff
from matches
         join (select match_id_fk, sum(game_phases.duration) duration_sum
               from game_phases
               group by match_id_fk) as mifds on match_id_fk = matches.id
where abs(EXTRACT(epoch FROM duration - duration_sum)) > 1
order by diff desc;

-- there should not be duplicate key names due to different writing (check it manually)
select distinct team_name
from team_match_stats
order by team_name;

-- check sum of accepted/proposed/withdrawn events
select game_event_type,
       total,
       accepted + accepted_withdrawn + proposed + proposed_withdrawn as sum
from game_event_type_metrics
where total <> accepted + accepted_withdrawn + proposed + proposed_withdrawn;
select game_event_type,
       total,
       accepted + accepted_withdrawn + proposed + proposed_withdrawn as sum
from game_event_type_metrics_per_team
where total <> accepted + accepted_withdrawn + proposed + proposed_withdrawn;
select game_event_category,
       total,
       accepted + accepted_withdrawn + proposed + proposed_withdrawn as sum
from game_event_category_metrics
where total <> accepted + accepted_withdrawn + proposed + proposed_withdrawn;
select game_event_category,
       total,
       accepted + accepted_withdrawn + proposed + proposed_withdrawn as sum
from game_event_category_metrics_per_team
where total <> accepted + accepted_withdrawn + proposed + proposed_withdrawn;

-- check Game Event Type vs. Category stats
select c.tournament_name,
       c.game_event_category,
       t.total,
       c.total,
       t.accepted,
       c.accepted
from (select game_event_category,
             tournament_name,
             sum(total)    total,
             sum(accepted) accepted
      from game_event_type_metrics c
      group by game_event_category, tournament_name) as t
         join (select game_event_category,
                      tournament_name,
                      sum(total)    total,
                      sum(accepted) accepted
               from game_event_category_metrics
               group by game_event_category, tournament_name) c
              on c.game_event_category = t.game_event_category and c.tournament_name = t.tournament_name
where t.total <> c.total
   or t.accepted <> c.accepted;
select c.tournament_name,
       c.team_name,
       c.game_event_category,
       t.total,
       c.total,
       t.accepted,
       c.accepted
from (select game_event_category,
             tournament_name,
             team_name,
             sum(total)    total,
             sum(accepted) accepted
      from game_event_type_metrics_per_team c
      group by game_event_category, tournament_name, team_name) as t
         join (select game_event_category,
                      tournament_name,
                      team_name,
                      sum(total)    total,
                      sum(accepted) accepted
               from game_event_category_metrics_per_team
               group by game_event_category, tournament_name, team_name) c
              on c.game_event_category = t.game_event_category and c.tournament_name = t.tournament_name and
                 c.team_name = t.team_name
where t.total <> c.total
   or t.accepted <> c.accepted;

-- Number of goals and number of accepted goal game events should be similar.
select s.match_id_fk,
       m.tournament_name,
       m.division,
       m.file_name,
       accepted,
       goals
from (select ge.match_id_fk,
             ge.by_team,
             count(*) as accepted
      from game_events_per_match ge
      where type = 'GOAL'
        and withdrawn = false
        and proposed = false
      group by ge.match_id_fk, by_team) s
         join team_match_stats tms on tms.match_id_fk = s.match_id_fk and tms.team_color = s.by_team
         join matches m on m.id = s.match_id_fk
where abs(accepted - goals) > 2;
