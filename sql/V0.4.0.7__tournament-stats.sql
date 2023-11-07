create or replace view tournament_stats as
select matches.tournament_id_fk,
       matches.division,
       count(distinct team_name)  as num_teams,
       count(distinct matches.id) as num_matches
from team_match_stats
         join matches on matches.id = team_match_stats.match_id_fk
group by matches.tournament_id_fk, matches.division;
