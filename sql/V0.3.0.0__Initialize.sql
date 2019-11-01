create table log_files
(
    id        uuid primary key,
    file_name varchar(255) not null unique
);

CREATE TYPE team_color AS ENUM ('yellow', 'blue', 'neutral');

create table matches
(
    id                      uuid primary key,
    log_file_id_fk          uuid references log_files (id),
    team_color              team_color,
    team_name               varchar(255),
    opponent_name           varchar(255),
    goals                   int,
    goals_conceded          int,
    fouls                   int,
    cards_yellow            int,
    cards_red               int,
    timeout_time            int,
    timeouts_taken          int,
    timeouts_left           int,
    ball_placement_time     int,
    ball_placements         int,
    max_active_yellow_cards int,
    penalty_shots_total     int,
    constraint unique_log_file unique (log_file_id_fk, team_color)
);
