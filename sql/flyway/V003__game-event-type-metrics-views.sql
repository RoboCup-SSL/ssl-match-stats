create view game_event_type_metrics_overall as
select type,
       count(*)                                               as total,
       count(*) FILTER (WHERE not withdrawn and not proposed) as accepted,
       count(*) FILTER (WHERE withdrawn and not proposed)     as accepted_withdrawn,
       count(*) FILTER (WHERE not withdrawn and proposed)     as proposed,
       count(*) FILTER (WHERE withdrawn and proposed)         as proposed_withdrawn
from game_events
group by type
order by type;

create view game_event_type_metrics_per_game_phase as
select gp.type                                                as game_phase_type,
       ge.type,
       count(*)                                               as total,
       count(*) FILTER (WHERE not withdrawn and not proposed) as accepted,
       count(*) FILTER (WHERE withdrawn and not proposed)     as accepted_withdrawn,
       count(*) FILTER (WHERE not withdrawn and proposed)     as proposed,
       count(*) FILTER (WHERE withdrawn and proposed)         as proposed_withdrawn
from game_events ge
         join game_phases gp on ge.game_phase_id_fk = gp.id
group by ge.type, gp.type
order by gp.type, ge.type;

create view game_event_type_metrics_per_match as
select m.id,
       ge.type,
       count(*)                                               as total,
       count(*) FILTER (WHERE not withdrawn and not proposed) as accepted,
       count(*) FILTER (WHERE withdrawn and not proposed)     as accepted_withdrawn,
       count(*) FILTER (WHERE not withdrawn and proposed)     as proposed,
       count(*) FILTER (WHERE withdrawn and proposed)         as proposed_withdrawn
from game_events ge
         join game_phases gp on ge.game_phase_id_fk = gp.id
         join matches m on gp.match_id_fk = m.id
group by ge.type, m.id
order by m.id, ge.type;

create view game_event_type_metrics_per_team as
select tms.team_name,
       ge.type,
       count(*)                                               as total,
       count(*) FILTER (WHERE not withdrawn and not proposed) as accepted,
       count(*) FILTER (WHERE withdrawn and not proposed)     as accepted_withdrawn,
       count(*) FILTER (WHERE not withdrawn and proposed)     as proposed,
       count(*) FILTER (WHERE withdrawn and proposed)         as proposed_withdrawn
from game_events ge
         join game_phases gp on ge.game_phase_id_fk = gp.id
         join matches m on gp.match_id_fk = m.id
         join team_match_stats tms on m.id = tms.match_id_fk
group by ge.type, tms.team_name
order by tms.team_name, ge.type;

create view game_event_type_metrics_per_tournament as
select m.tournament_name,
       m.division,
       ge.type,
       count(*)                                               as total,
       count(*) FILTER (WHERE not withdrawn and not proposed) as accepted,
       count(*) FILTER (WHERE withdrawn and not proposed)     as accepted_withdrawn,
       count(*) FILTER (WHERE not withdrawn and proposed)     as proposed,
       count(*) FILTER (WHERE withdrawn and proposed)         as proposed_withdrawn
from game_events ge
         join game_phases gp on ge.game_phase_id_fk = gp.id
         join matches m on gp.match_id_fk = m.id
group by ge.type, m.tournament_name, m.division
order by m.tournament_name, m.division, ge.type;

