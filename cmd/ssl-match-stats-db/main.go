package main

import (
	"flag"
	"fmt"
	"github.com/RoboCup-SSL/ssl-match-stats/pkg/matchstats"
	"github.com/RoboCup-SSL/ssl-match-stats/pkg/sqldbexport"
	"log"
	"os"
)

func main() {
	flag.Usage = usage

	sqlDriver := flag.String("sqlDriver", "postgres", "SQL driver")
	sqlDbSource := flag.String("sqlDbSource", "postgres://user:password@host:port/ssl_match_stats", "SQL connection string")

	//tournament := flag.String("tournament", "", "The tournament the log files are for")

	flag.Parse()

	exporter := sqldbexport.SqlDbExporter{}
	if err := exporter.Connect(*sqlDriver, *sqlDbSource); err != nil {
		log.Fatalf("Could not connect to database with driver '%v' at '%v'", *sqlDriver, *sqlDbSource)
	}

	a := matchstats.NewAggregator()

	if err := a.ReadBin("out.bin"); err != nil {
		log.Fatal(err)
	}

	if err := exporter.WriteLogFiles(&a.Collection); err != nil {
		log.Fatal("Could not write log files: ", err)
	}
	if err := exporter.WriteTeamStats(&a.Collection); err != nil {
		log.Fatal("Could not write team stats: ", err)
	}
}

func usage() {
	_, err := fmt.Fprintln(os.Stderr, "Pass one or more log files that should be processed.")
	if err != nil {
		fmt.Println("Pass one or more log files that should be processed.")
	}
	flag.PrintDefaults()
}
