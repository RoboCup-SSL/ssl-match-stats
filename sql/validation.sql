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
select type,
       total,
       accepted + accepted_withdrawn + proposed + proposed_withdrawn as sum
from game_event_type_metrics_overall
where total <> accepted + accepted_withdrawn + proposed + proposed_withdrawn;
