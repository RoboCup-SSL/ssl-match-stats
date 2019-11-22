alter table matches alter column division TYPE varchar(4);
drop type division;

alter table team_match_stats alter column team_color type varchar(7);
alter table game_phases alter column for_team type varchar(7);
alter table game_phases alter column entry_command_for_team type varchar(7);
alter table game_phases alter column exit_command_for_team type varchar(7);
alter table game_phases alter column proposed_next_command_for_team type varchar(7);
alter table game_phases alter column previous_command_for_team type varchar(7);
drop type team_color;

drop view game_phase_duration_per_type;
alter table game_phases alter column type type varchar(255);
drop type game_phase;

alter table game_phases alter column entry_command type varchar(255);
alter table game_phases alter column exit_command type varchar(255);
alter table game_phases alter column proposed_next_command type varchar(255);
alter table game_phases alter column previous_command type varchar(255);
drop type command;

alter table game_phases alter column stage type varchar(255);
drop type stage;

alter table game_events alter column type type varchar(255);
drop type game_event;

create view game_phase_duration_per_type AS
select match_id_fk,
       type,
       sum(duration) as duration
from game_phases
group by match_id_fk, type;
