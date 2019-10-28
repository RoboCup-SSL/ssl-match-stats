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

The binary is called `ssl-match-stats`.

### Generate statistics from log files

Pass in a list of log files to be processed: `ssl-match-stats -generate *.log.gz`

### Export statistics to CSV files

First, generate the statistics with the command above. This will produce a `out.json` and `out.bin` file.

Run: `ssl-match-stats -exportCsv`

This will generate `*.csv` files that you can import in your favorite tool, like a spreadsheet tool.
