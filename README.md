[![CircleCI](https://circleci.com/gh/RoboCup-SSL/ssl-match-stats/tree/master.svg?style=svg)](https://circleci.com/gh/RoboCup-SSL/ssl-match-stats/tree/master)

# ssl-match-stats

A tool that generates statistics from [official SSL log files](https://ssl.robocup.org/game-logs/).

## Requirements

You need to install following dependencies:

* Go

## Usage

### Build

```shell
# Build the binaries and install them in $GOPATH/bin
make install
```

### Generate Statistics

The `ssl-match-stats` command will generate the statistics into an intermediate data structure `*.{json|bin}`
from which they can be exported. This must be done per tournament and division.

The command takes a list of log files as input:

```
mkdir -p stats
ssl-match-stats -parallel 16 -targetDir stats *.log.gz
```

### Import Statistics into a Database

The generated statistics can be exported into a PostgreSQL database.
This makes it easier to query and analyze the data.

The command requires some parameters:

* tournament: A unique name for the tournament of the log files, for example 'RoboCup2019'
* division: The division of the log files, one of: 'DivA', 'DivB', 'none'
* sqlDbSource: A connection string to the target database

```shell
ssl-match-stats-db \
  -parallel=16 \
  -sqlDbSource="postgres://ssl_match_stats:ssl_match_stats@localhost:5432/ssl_match_stats?sslmode=disable" \
  -tournament=Test \
  -division=DivA \
  stats/*.bin
```

### Local Setup

Use the provided compose setup to run the database and metabase locally.

```shell
docker compose up -d
```

The database connection string
is: `postgres://ssl_match_stats:ssl_match_stats@localhost:5432/ssl_match_stats?sslmode=disable`

### Remote Setup

The host the statistics one a server, following the below instructions.

See [Setup for Match Stats DB](./setup/matchStatsDb/README.md) for instructions on setting up the PostgreSQL database.

See [Setup for Metabase](./setup/metabase/README.md) for instructions on setting up Metabase, an open-source BI
software.

### Manual data correction

If you want to evaluate the start time offset between scheduled match start time and actual start time,
you need to add the scheduled time into `matches.start_time_planned` manually.

### Protobuf

To generate the sources from the `.proto` files, run `make proto`.

## Implementation Details

This tool reads log files one by one with the `ssl-match-stats` command and creates a protobuf file based
on [ssl_match_stats.proto](./proto/ssl_match_stats.proto).
This structure contains a `MatchStats` object per match (log file), which contains:

* Some meta data like name, duration, etc.
* A list of all game phases (which map roughly to referee commands)
* `TeamStats` for each team

The incentive is that the `MatchStats` structure contains an improved representation of the referee messages
for further processing. The match is split into game phases, where a new game phase is started for each new command
(except the goal commands). Each game phase has metadata, entry and exit states and game events attached to it.
`TeamStats` contain aggregated final counters and timers.

With this improved representation, the data can be exported for further analysis.
