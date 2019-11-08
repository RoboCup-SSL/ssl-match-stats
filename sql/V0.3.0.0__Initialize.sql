-- a tournament could be e.g. RoboCup 2019
create table tournaments
(
    id   uuid primary key,
    name varchar(255)
);

-- custom type for divisions
create type division as enum ('DivA', 'DivB', 'none');

-- a single match between two teams
create table matches
(
    id               uuid primary key,
    file_name        varchar(255) not null unique,
    tournament_id_fk uuid references tournaments (id),
    download_link    varchar(255),
    division         division,
    start_time       timestamp,
    duration         interval(3), -- millisecond precision
    extra_time       bool,
    shootout         bool
);

-- custom enum for team colors
create type team_color as enum ('UNKNOWN', 'YELLOW', 'BLUE', 'NONE');

-- statistics of a single match from the perspective of a single team -> two rows per match
create table team_match_stats
(
    id                      uuid primary key,
    match_id_fk             uuid references matches (id),
    team_color              team_color,
    team_name               varchar(255),
    opponent_name           varchar(255),
    goals                   int,
    conceded_goals          int,
    fouls                   int,
    yellow_cards            int,
    red_cards               int,
    timeout_time            interval(3), -- millisecond precision
    timeouts_taken          int,
    timeouts_left           int,
    ball_placement_time     interval(3), -- millisecond precision
    ball_placements         int,
    max_active_yellow_cards int,
    penalty_shots_total     int,
    constraint unique_match unique (match_id_fk, team_color)
);

-- custom type for game phases (as derived from a stream of referee messages)
create type game_phase as enum (
    'UNKNOWN',
    'RUNNING',
    'PREPARE_KICKOFF',
    'PREPARE_PENALTY',
    'STOP',
    'BALL_PLACEMENT',
    'TIMEOUT',
    'BREAK',
    'HALT'
    );

-- custom type for referee commands from referee messages
create type command as enum (
    'UNKNOWN',
    'HALT',
    'STOP',
    'BALL_PLACEMENT',
    'NORMAL_START',
    'FORCE_START',
    'DIRECT_FREE',
    'INDIRECT_FREE',
    'PREPARE_KICKOFF',
    'PREPARE_PENALTY',
    'TIMEOUT',
    'GOAL'
    );

-- custom type for stages from referee messages
create type stage as enum (
    'UNKNOWN',
    'NORMAL_FIRST_HALF_PRE',
    'NORMAL_FIRST_HALF',
    'NORMAL_HALF_TIME',
    'NORMAL_SECOND_HALF_PRE',
    'NORMAL_SECOND_HALF',
    'EXTRA_TIME_BREAK',
    'EXTRA_FIRST_HALF_PRE',
    'EXTRA_FIRST_HALF',
    'EXTRA_HALF_TIME',
    'EXTRA_SECOND_HALF_PRE',
    'EXTRA_SECOND_HALF',
    'PENALTY_SHOOTOUT_BREAK',
    'PENALTY_SHOOTOUT'
    );

-- derived game phases from a stream of referee messages
create table game_phases
(
    id                             uuid primary key,
    match_id_fk                    uuid references matches (id),
    start_time                     timestamp,
    end_time                       timestamp,
    duration                       interval(3), -- millisecond precision
    type                           game_phase,
    for_team                       team_color,
    entry_command                  command,
    entry_command_for_team         team_color,
    exit_command                   command,
    exit_command_for_team          team_color,
    proposed_next_command          command,
    proposed_next_command_for_team team_color,
    previous_command               command,
    previous_command_for_team      team_color,
    stage                          stage,
    stage_time_left_entry          interval(3), -- millisecond precision
    stage_time_left_exit           interval(3)  -- millisecond precision
);

-- grant read access to view user
create user ssl_match_stats_view with encrypted password 'ssl_match_stats_view';
grant usage on schema public to ssl_match_stats_view;

grant select on table tournaments to ssl_match_stats_view;
grant select on table matches to ssl_match_stats_view;
grant select on table team_match_stats to ssl_match_stats_view;
grant select on table game_phases to ssl_match_stats_view;
