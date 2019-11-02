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
create type team_color as enum ('yellow', 'blue', 'neutral');

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
    'running',
    'prepare_kickoff',
    'prepare_penalty',
    'stop',
    'ball_placement',
    'timeout',
    'break',
    'halt'
    );

-- custom type for referee commands from referee messages
create type command as enum (
    'halt',
    'stop',
    'ball_placement',
    'normal_start',
    'force_start',
    'direct_free',
    'indirect_free',
    'prepare_kickoff',
    'prepare_penalty',
    'timeout',
    'goal'
    );

-- custom type for stages from referee messages
create type stage as enum (
    'normal_first_half_pre',
    'normal_first_half',
    'normal_half_time',
    'normal_second_half_pre',
    'normal_second_half',
    'extra_time_break',
    'extra_first_half_pre',
    'extra_first_half',
    'extra_half_time',
    'extra_second_half_pre',
    'extra_second_half',
    'penalty_shootout_break',
    'penalty_shootout'
    );

-- derived game phases from a stream of referee messages
create table game_phases
(
    id                    uuid primary key,
    match_id_fk           uuid references matches (id),
    start_time            timestamp,
    end_time              timestamp,
    duration              interval(3), -- millisecond precision
    type                  game_phase,
    for_team              team_color,
    entry_command         command,
    exit_command          command,
    proposed_next_command command,
    previous_command      command,
    stage                 stage,
    stage_time_left_entry interval(3), -- millisecond precision
    stage_time_left_exit  interval(3)  -- millisecond precision
);

-- grant read access to view user
grant select on table tournaments to ssl_match_stats_view;
grant select on table matches to ssl_match_stats_view;
grant select on table team_match_stats to ssl_match_stats_view;
grant select on table game_phases to ssl_match_stats_view;
