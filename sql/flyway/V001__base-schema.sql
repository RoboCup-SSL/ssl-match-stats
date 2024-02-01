-- custom type for divisions
create type division as enum ('DivA', 'DivB', 'none');

create type match_type as enum (
    'UNKNOWN',
    'GROUP_PHASE',
    'ELIMINATION_PHASE',
    'FRIENDLY'
    );

-- a single match between two teams
create table matches
(
    id                 uuid primary key,
    file_name          varchar(255) not null unique,
    tournament_name    varchar(255) not null,
    type               match_type,
    division           division,
    start_time         timestamp,
    start_time_planned timestamp,
    duration           interval(3), -- millisecond precision
    extra_time         bool,
    shootout           bool
);

-- custom enum for team colors
create type team_color as enum ('UNKNOWN', 'YELLOW', 'BLUE', 'NONE');

-- statistics of a single match from the perspective of a single team -> two rows per match
create table team_match_stats
(
    id                      uuid primary key,
    match_id_fk             uuid references matches (id) ON DELETE CASCADE,
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
    'HALT',
    'POST_GAME'
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
    'PENALTY_SHOOTOUT',
    'POST_GAME'
    );

-- derived game phases from a stream of referee messages
create table game_phases
(
    id                             uuid primary key,
    match_id_fk                    uuid references matches (id) ON DELETE CASCADE,
    start_time                     timestamp,
    end_time                       timestamp,
    duration                       interval(3), -- millisecond precision
    type                           game_phase,
    for_team                       team_color,
    entry_command                  command,
    entry_command_for_team         team_color,
    entry_command_timestamp        timestamp,
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

-- custom type for game events
create type game_event as enum (
    'UNKNOWN_GAME_EVENT_TYPE',
    'PREPARED',
    'NO_PROGRESS_IN_GAME',
    'PLACEMENT_FAILED',
    'PLACEMENT_SUCCEEDED',
    'BOT_SUBSTITUTION',
    'TOO_MANY_ROBOTS',
    'BALL_LEFT_FIELD_TOUCH_LINE',
    'BALL_LEFT_FIELD_GOAL_LINE',
    'POSSIBLE_GOAL',
    'GOAL',
    'INDIRECT_GOAL',
    'CHIPPED_GOAL',
    'AIMLESS_KICK',
    'KICK_TIMEOUT',
    'KEEPER_HELD_BALL',
    'ATTACKER_DOUBLE_TOUCHED_BALL',
    'ATTACKER_TOUCHED_BALL_IN_DEFENSE_AREA',
    'ATTACKER_TOUCHED_OPPONENT_IN_DEFENSE_AREA',
    'ATTACKER_TOUCHED_OPPONENT_IN_DEFENSE_AREA_SKIPPED',
    'BOT_DRIBBLED_BALL_TOO_FAR',
    'BOT_KICKED_BALL_TOO_FAST',
    'ATTACKER_TOO_CLOSE_TO_DEFENSE_AREA',
    'BOT_INTERFERED_PLACEMENT',
    'BOT_CRASH_DRAWN',
    'BOT_CRASH_UNIQUE',
    'BOT_CRASH_UNIQUE_SKIPPED',
    'BOT_PUSHED_BOT',
    'BOT_PUSHED_BOT_SKIPPED',
    'BOT_HELD_BALL_DELIBERATELY',
    'BOT_TIPPED_OVER',
    'BOT_TOO_FAST_IN_STOP',
    'DEFENDER_TOO_CLOSE_TO_KICK_POINT',
    'DEFENDER_IN_DEFENSE_AREA_PARTIALLY',
    'DEFENDER_IN_DEFENSE_AREA',
    'MULTIPLE_CARDS',
    'MULTIPLE_PLACEMENT_FAILURES',
    'MULTIPLE_FOULS',
    'UNSPORTING_BEHAVIOR_MINOR',
    'UNSPORTING_BEHAVIOR_MAJOR',
    'BOUNDARY_CROSSING',
    'INVALID_GOAL',
    'PENALTY_KICK_FAILED',
    'CHALLENGE_FLAG',
    'EMERGENCY_STOP',
    'CHALLENGE_FLAG_HANDLED'
    );

create type game_event_category as enum (
    'UNKNOWN',
    'BALL_OUT',
    'FOUL',
    'GOAL',
    'OTHER'
    );

-- game events with reference to game phases
create table game_events
(
    id                uuid primary key,
    game_phase_id_fk  uuid references game_phases (id) on delete cascade,
    type              game_event          not null,
    by_team           varchar(7)          not null,
    timestamp         timestamp           not null,
    -- timestamp when event was created by the GC. This data was not available in the past, so it might be null.
    created_timestamp timestamp,
    withdrawn         bool                not null,
    proposed          bool                not null,
    payload           jsonb               not null
);

-- mapping between game events and origins (note: there can be multiple origins per event)
create table game_event_origin_mapping
(
    game_event_id_fk  uuid references game_events (id) on delete cascade,
    game_event_origin varchar,
    primary key (game_event_id_fk, game_event_origin)
);

create table game_event_categories
(
    type game_event primary key,
    category game_event_category not null
);