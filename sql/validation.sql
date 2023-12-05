-- game phase durations should not be negative
select * from game_phases gp
join matches m on m.id = gp.match_id_fk
where gp.duration < '0s'::interval;

-- there should be no very long game phases
select * from game_phases gp
                  join matches m on m.id = gp.match_id_fk
where gp.duration > '1h'::interval;

-- sum of game phase duration should the same as match duration
select matches.file_name, duration, duration_sum, duration - duration_sum as diff from matches
join (select match_id_fk, sum(game_phases.duration) duration_sum from game_phases
group by match_id_fk) as mifds on match_id_fk = matches.id
order by diff desc;