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
	sqlDbSource := flag.String("sqlDbSource", "", "SQL connection string, for example: postgres://user:password@host:port/ssl_match_stats")
	tournament := flag.String("tournament", "", "The tournament the log files are from")
	division := flag.String("division", "", "The division of the log files. Should be one of: DivA, DivB, none")

	flag.Parse()

	if len(*sqlDbSource) == 0 {
		log.Fatal("You have to specify a db source")
	}
	if len(*tournament) == 0 {
		log.Fatal("You have to specify the tournament name")
	}
	if len(*division) == 0 {
		log.Fatal("You have to specify the division")
	} else if !validDivision(*division) {
		log.Fatal("The division must be one of: DivA, DivB or none")
	}

	exporter := sqldbexport.SqlDbExporter{}
	if err := exporter.Connect(*sqlDriver, *sqlDbSource); err != nil {
		log.Fatalf("Could not connect to database with driver '%v' at '%v'", *sqlDriver, *sqlDbSource)
	}

	a := matchstats.NewCollector()

	if err := a.ReadBin("match-stats.bin"); err != nil {
		log.Fatal(err)
	}

	tournamentId, err := exporter.AddTournamentIfNotPresent(*tournament)
	if err != nil {
		log.Fatal(err)
	}

	if err := exporter.WriteMatches(a.Collection, tournamentId, *division); err != nil {
		log.Fatal("Could not write matches: ", err)
	}
}

func validDivision(s string) bool {
	validValues := []string{"DivA", "DivB", "none"}
	for _, validValue := range validValues {
		if validValue == s {
			return true
		}
	}
	return false
}

func usage() {
	_, err := fmt.Fprintln(os.Stderr, "Pass one or more log files that should be processed.")
	if err != nil {
		fmt.Println("Pass one or more log files that should be processed.")
	}
	flag.PrintDefaults()
}
