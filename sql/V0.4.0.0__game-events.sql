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
    'UNSPORTING_BEHAVIOR_MAJOR'
    );

-- all autoRefs that are referenced in game events
create table auto_refs
(
    id   uuid primary key,
    name varchar(255) unique
);

-- game events with reference to game phases
create table game_events
(
    id               uuid primary key,
    game_phase_id_fk uuid references game_phases (id),
    type             game_event,
    timestamp        timestamp,
    withdrawn        bool,
    payload          jsonb
);

-- mapping between game events and autoRefs
create table game_event_auto_ref_mapping
(
    auto_ref_id_fk   uuid references auto_refs (id),
    game_event_id_fk uuid references game_events (id),
    primary key (auto_ref_id_fk, game_event_id_fk)
);

-- grant read access to view user
grant select on table auto_refs to ssl_match_stats_view;
grant select on table game_events to ssl_match_stats_view;
grant select on table game_event_auto_ref_mapping to ssl_match_stats_view;
