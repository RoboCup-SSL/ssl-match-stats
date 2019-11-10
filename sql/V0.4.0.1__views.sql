create view game_phase_duration_per_type AS
select match_id_fk,
       type,
       sum(duration) as duration
from game_phases
group by match_id_fk, type;
