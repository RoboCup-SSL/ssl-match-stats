select * from game_phases gp
join matches m on m.id = gp.match_id_fk
where gp.duration < '0s'::interval;

select * from game_phases gp
                  join matches m on m.id = gp.match_id_fk
where gp.duration > '1h'::interval;