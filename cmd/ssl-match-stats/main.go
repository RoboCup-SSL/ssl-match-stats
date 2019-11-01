package main

import (
	"flag"
	"fmt"
	"github.com/RoboCup-SSL/ssl-match-stats/pkg/csvexport"
	"github.com/RoboCup-SSL/ssl-match-stats/pkg/matchstats"
	"github.com/RoboCup-SSL/ssl-match-stats/pkg/sqldbexport"
	"log"
	"os"
)

func main() {
	flag.Usage = usage

	fGenerate := flag.Bool("generate", false, "Generate statistics based on passed in log files")
	fExportCsv := flag.Bool("exportCsv", false, "Export data from a generated out.bin file to CSV")
	fExportSqlDb := flag.Bool("exportSqlDb", false, "Export data from a generated out.bin file to SQL DB")
	sqlDriver := flag.String("sqlDriver", "postgres", "SQL driver")
	sqlDbSource := flag.String("sqlDbSource", "postgres://user:password@172.16.1.1:5432/ssl_match_stats", "SQL connection string")

	flag.Parse()

	if *fGenerate {
		generate()
	}
	if *fExportCsv {
		exportCsv()
	}
	if *fExportSqlDb {
		exportSqlDb(*sqlDriver, *sqlDbSource)
	}
}

func generate() {
	args := flag.Args()

	if len(args) == 0 {
		usage()
		return
	}

	a := matchstats.NewAggregator()
	for _, filename := range args {
		log.Println("Processing", filename)

		err := a.Process(filename)
		if err != nil {
			log.Printf("%v: %v\n", filename, err)
		} else {
			log.Printf("Processed %v\n", filename)
		}
	}

	if err := a.Aggregate(); err != nil {
		log.Println("Failed to aggregate match stats", err)
	}

	if err := a.WriteBin("out.bin"); err != nil {
		log.Println("Could not write binary file", err)
	}

	if err := a.WriteJson("out.json"); err != nil {
		log.Println("Could not write JSON file", err)
	}
}

func exportCsv() {

	a := matchstats.NewAggregator()

	if err := a.ReadBin("out.bin"); err != nil {
		log.Fatal(err)
	}

	if err := csvexport.WriteGamePhases(&a.Collection, "game-phases.csv"); err != nil {
		log.Fatal(err)
	}
	if err := csvexport.WriteGamePhaseDurations(&a.Collection, "game-phase-durations.csv"); err != nil {
		log.Fatal(err)
	}

	if err := csvexport.WriteTeamMetricsPerGame(&a.Collection, "team-metrics-per-game.csv"); err != nil {
		log.Fatal(err)
	}
	if err := csvexport.WriteTeamMetricsSum(&a.Collection, "team-metrics-sum.csv"); err != nil {
		log.Fatal(err)
	}

	if err := csvexport.WriteGamePhaseDurationStats(&a.Collection, "game-phase-duration-stats.csv"); err != nil {
		log.Fatal(err)
	}
	if err := csvexport.WriteGamePhaseDurationStatsAggregated(&a.Collection, "game-phase-duration-stats-aggregated.csv"); err != nil {
		log.Fatal(err)
	}
	if err := csvexport.WriteGameEventDurationStats(&a.Collection, "game-event-duration-stats.csv"); err != nil {
		log.Fatal(err)
	}
	if err := csvexport.WriteGameEventDurationStatsAggregated(&a.Collection, "game-event-duration-stats-aggregated.csv"); err != nil {
		log.Fatal(err)
	}
}

func exportSqlDb(sqlDriver string, sqlDbSource string) {
	exporter := sqldbexport.SqlDbExporter{}
	if err := exporter.Connect(sqlDriver, sqlDbSource); err != nil {
		log.Fatalf("Could not connect to database with driver '%v' at '%v'", sqlDriver, sqlDbSource)
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
