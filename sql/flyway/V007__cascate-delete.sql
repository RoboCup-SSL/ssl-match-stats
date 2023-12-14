alter table team_match_stats
    DROP CONSTRAINT team_match_stats_match_id_fk_fkey,
    ADD CONSTRAINT team_match_stats_match_id_fk_fkey
        FOREIGN KEY (match_id_fk)
            REFERENCES matches (id)
            ON DELETE CASCADE;

alter table game_phases
    DROP CONSTRAINT game_phases_match_id_fk_fkey,
    ADD CONSTRAINT game_phases_match_id_fk_fkey
        FOREIGN KEY (match_id_fk)
            REFERENCES matches (id)
            ON DELETE CASCADE;
