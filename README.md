[![CircleCI](https://circleci.com/gh/RoboCup-SSL/ssl-match-stats/tree/master.svg?style=svg)](https://circleci.com/gh/RoboCup-SSL/ssl-match-stats/tree/master)
[![Go Report Card](https://goreportcard.com/badge/github.com/RoboCup-SSL/ssl-match-stats?style=flat-square)](https://goreportcard.com/report/github.com/RoboCup-SSL/ssl-match-stats)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/RoboCup-SSL/ssl-match-stats)
[![Coverage](https://img.shields.io/badge/coverage-report-blue.svg)](https://circleci.com/api/v1.1/project/github/RoboCup-SSL/ssl-match-stats/latest/artifacts/0/coverage?branch=master)


# ssl-match-stats

A tool that generates statistics from [official SSL log files](https://ssl.robocup.org/game-logs/).

## Requirements
You need to install following dependencies first: 
 * Go >= 1.11
 
## Installation

Use go get to install all packages / executables:

```
go get -u github.com/RoboCup-SSL/ssl-match-stats/...
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
