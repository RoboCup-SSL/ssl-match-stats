-- the view can be used with Metabase to great a histogram of game phase durations
create view game_phase_duration as
select id, type, extract(epoch FROM duration) duration_seconds from game_phases;
