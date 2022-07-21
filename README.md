[![CircleCI](https://circleci.com/gh/RoboCup-SSL/ssl-match-stats/tree/master.svg?style=svg)](https://circleci.com/gh/RoboCup-SSL/ssl-match-stats/tree/master)
[![Go Report Card](https://goreportcard.com/badge/github.com/RoboCup-SSL/ssl-match-stats?style=flat-square)](https://goreportcard.com/report/github.com/RoboCup-SSL/ssl-match-stats)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/RoboCup-SSL/ssl-match-stats)
[![Coverage](https://img.shields.io/badge/coverage-report-blue.svg)](https://circleci.com/api/v1.1/project/github/RoboCup-SSL/ssl-match-stats/latest/artifacts/0/coverage?branch=master)


# ssl-match-stats

A tool that generates statistics from [official SSL log files](https://ssl.robocup.org/game-logs/).

## Requirements
You need to install following dependencies first: 
 * Go >= 1.16
 
## Installation

Use go get to install all packages / executables:

```
go install github.com/RoboCup-SSL/ssl-match-stats/...@latest
```

## Run
The executables are installed to your $GOPATH/bin folder. If you have it on your $PATH, you can directly run them. Else, switch to this folder first.

## Usage

### Generate Statistics

The `ssl-match-stats` command will generate the statistics into an intermediate data structure `match-stats.{json|bin}`
from which they can be exported. This must be done per tournament and division. 

The command takes a list of log files as input:
```
ssl-match-stats *.log.gz
```

### Export Statistics to CSV files

The generated statistics can be exported into CSV files for further processing, 
for example with a spreadsheet software or Matlab. 
The `ssl-match-stats-csv` command will read the `match-stats.bin` protobuf file 
from the current folder and produces a set of CSV files: 

```
ssl-match-stats-csv
```

### Export Statistics to a Database

The generated statistics can be exported into a PostgreSQL database (other DBs not yet tested).
This is useful if you want to apply some BI (Business Intelligence) application on the data.

See [Setup for Match Stats DB](./setup/matchStatsDb/README.md) for instructions on setting up the database.

See [Setup for Metabase](./setup/metabase/README.md) for instructions on setting up Metabase, an open-source BI software.

The command requires some parameters:

* tournament: A unique name for the tournament of the log files, for example 'RoboCup2019'
* division: The division of the log files, one of: 'DivA', 'DivB', 'none'
* sqlDbSource: A connection string to the target database

```
ssl-match-stats-db -sqlDbSource postgres://user:password@host:port/db-name -tournament RoboCup2019 -division DivA
```

If you want to evaluate the start time offset between scheduled match start time and actual start time,
you need to add the scheduled time into `matches.start_time_planned` manually.

## Development

### Database

A docker-compose config is provided to startup a postgres DB.
You have to install `docker` and `docker-compose` for that. 

To start, run:
```shell script
docker-compose up
```
The database connection string is: `postgres://ssl_match_stats:ssl_match_stats@localhost:5432/ssl_match_stats?sslmode=disable`

The database schema will automatically be installed with flyway. 
Data will be stored in a volume, so it is not lost when stopping docker-compose.
To reset the database, run `docker-compose down`.

### Protobuf

To generate the sources from the `.proto` files, run the [generateProto.sh](./generateProto.sh) script.


## Implementation Details

This tool reads log files one by one with the `ssl-match-stats` command and creates a protobuf file based on [ssl_match_stats.proto](./proto/ssl_match_stats.proto).
This structure contains a `MatchStats` object per match (log file), which contains:

 * Some meta data like name, duration, etc.
 * A list of all game phases (which map roughly to referee commands)
 * `TeamStats` for each team

The incentive is that the `MatchStats` structure contains an improved representation of the referee messages
for further processing. The match is split into game phases, where a new game phase is started for each new command 
(except the goal commands). Each game phase has meta data, entry and exit states and game events attached to it.
`TeamStats` contain aggregated final counters and timers.

With this improved representation, the data can be exported for further analysis. 
There is currently a CSV exporter which does some additional aggregation and outputs simple CSV files.
As this type of output format is quite limited, there is also a database exporter which connects to a 
postgres DB. Data is added in a way that allows filtering and aggregation on database level.
