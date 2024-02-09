create view robot_count_total as
select rc.game_phase_id_fk,
       rc.start_time,
       rc.duration,
       sum(rc.count) as count
from robot_count rc
group by rc.game_phase_id_fk, rc.start_time, rc.duration;

drop view if exists robot_count_per_match cascade;
create view robot_count_per_match as
select s.match_id_fk,
       m.tournament_name,
       m.division,
       s.count,
       extract(epoch FROM s.duration) / extract(epoch FROM t.total_duration) * 100 as duration_percentage
from (select gp.match_id_fk,
             rc.count,
             sum(rc.duration) as duration
      from robot_count_total rc
               join game_phases gp on rc.game_phase_id_fk = gp.id
      where gp.type in ('RUNNING', 'PREPARE_KICKOFF', 'PREPARE_PENALTY', 'STOP', 'BALL_PLACEMENT')
        and gp.stage in ('NORMAL_FIRST_HALF', 'NORMAL_SECOND_HALF', 'EXTRA_FIRST_HALF', 'EXTRA_SECOND_HALF')
      group by match_id_fk, count) as s
         join (select gp.match_id_fk,
                      sum(rc.duration) as total_duration
               from robot_count_total rc
                        join game_phases gp on rc.game_phase_id_fk = gp.id
               where gp.type in ('RUNNING', 'PREPARE_KICKOFF', 'PREPARE_PENALTY', 'STOP', 'BALL_PLACEMENT')
                 and gp.stage in ('NORMAL_FIRST_HALF', 'NORMAL_SECOND_HALF', 'EXTRA_FIRST_HALF', 'EXTRA_SECOND_HALF')
               group by gp.match_id_fk) as t on s.match_id_fk = t.match_id_fk
         join matches m on s.match_id_fk = m.id
order by tournament_name, division, count;

drop view if exists robot_count_per_tournament cascade;
create view robot_count_per_tournament as
select s.tournament_name,
       s.division,
       s.count,
       extract(epoch FROM s.duration) / extract(epoch FROM t.total_duration) * 100 as duration_percentage
from (select m.tournament_name,
             m.division,
             rc.count,
             sum(rc.duration) as duration
      from robot_count_total rc
               join game_phases gp on rc.game_phase_id_fk = gp.id
               join matches m on gp.match_id_fk = m.id
      where gp.type in ('RUNNING', 'PREPARE_KICKOFF', 'PREPARE_PENALTY', 'STOP', 'BALL_PLACEMENT')
        and gp.stage in ('NORMAL_FIRST_HALF', 'NORMAL_SECOND_HALF', 'EXTRA_FIRST_HALF', 'EXTRA_SECOND_HALF')
      group by tournament_name, division, count) as s
         join (select m.tournament_name,
                      m.division,
                      sum(rc.duration) as total_duration
               from robot_count_total rc
                        join game_phases gp on rc.game_phase_id_fk = gp.id
                        join matches m on gp.match_id_fk = m.id
               where gp.type in ('RUNNING', 'PREPARE_KICKOFF', 'PREPARE_PENALTY', 'STOP', 'BALL_PLACEMENT')
                 and gp.stage in ('NORMAL_FIRST_HALF', 'NORMAL_SECOND_HALF', 'EXTRA_FIRST_HALF', 'EXTRA_SECOND_HALF')
               group by m.tournament_name, division) as t
              on s.tournament_name = t.tournament_name and s.division = t.division
order by s.tournament_name, s.division, s.count;

drop view if exists robot_count_per_team cascade;
create view robot_count_per_team as
select s.team_name,
       s.tournament_name,
       s.division,
       s.count,
       extract(epoch FROM s.duration) / extract(epoch FROM t.total_duration) * 100 as duration_percentage
from (select tms.team_name,
             m.tournament_name,
             m.division,
             rc.count,
             sum(rc.duration) as duration
      from robot_count rc
               join game_phases gp on rc.game_phase_id_fk = gp.id
               join matches m on gp.match_id_fk = m.id
               join team_match_stats tms on m.id = tms.match_id_fk and tms.team_color = rc.team_color
      where gp.type in ('RUNNING', 'PREPARE_KICKOFF', 'PREPARE_PENALTY', 'STOP', 'BALL_PLACEMENT')
        and gp.stage in ('NORMAL_FIRST_HALF', 'NORMAL_SECOND_HALF', 'EXTRA_FIRST_HALF', 'EXTRA_SECOND_HALF')
      group by tms.team_name, tournament_name, division, count) as s
         join (select tms.team_name,
                      m.tournament_name,
                      m.division,
                      sum(rc.duration) as total_duration
               from robot_count rc
                        join game_phases gp on rc.game_phase_id_fk = gp.id
                        join matches m on gp.match_id_fk = m.id
                        join team_match_stats tms on m.id = tms.match_id_fk and tms.team_color = rc.team_color
               where gp.type in ('RUNNING', 'PREPARE_KICKOFF', 'PREPARE_PENALTY', 'STOP', 'BALL_PLACEMENT')
                 and gp.stage in ('NORMAL_FIRST_HALF', 'NORMAL_SECOND_HALF', 'EXTRA_FIRST_HALF', 'EXTRA_SECOND_HALF')
               group by tms.team_name, m.tournament_name, division) as t
              on s.tournament_name = t.tournament_name and s.division = t.division and s.team_name = t.team_name
order by s.team_name, s.tournament_name, s.division, s.count;

drop view if exists robot_count_avg_per_match cascade;
create view robot_count_avg_per_match as
select rc.match_id_fk,
       rc.tournament_name,
       rc.division,
       sum(rc.count * rc.duration_percentage) / 100 as avg_count
from robot_count_per_match rc
group by rc.match_id_fk, tournament_name, division
order by tournament_name, division;

drop view if exists robot_count_avg_per_tournament cascade;
create view robot_count_avg_per_tournament as
select rc.tournament_name,
       rc.division,
       sum(rc.count * rc.duration_percentage) / 100 as avg_count
from robot_count_per_tournament rc
group by tournament_name, division
order by tournament_name, division;

drop view if exists robot_count_avg_per_team cascade;
create view robot_count_avg_per_team as
select rc.team_name,
       rc.tournament_name,
       rc.division,
       sum(rc.count * rc.duration_percentage) / 100 as avg_count
from robot_count_per_team rc
group by rc.team_name, tournament_name, division
order by rc.team_name, tournament_name, division;
