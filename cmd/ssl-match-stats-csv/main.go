package main

import (
	"flag"
	"fmt"
	"github.com/RoboCup-SSL/ssl-match-stats/pkg/csvexport"
	"github.com/RoboCup-SSL/ssl-match-stats/pkg/csvexport/aggregator"
	"github.com/RoboCup-SSL/ssl-match-stats/pkg/matchstats"
	"log"
	"os"
)

func main() {
	flag.Usage = usage

	flag.Parse()

	c := matchstats.NewCollector()

	if err := c.ReadBin("match-stats.bin"); err != nil {
		log.Fatal(err)
	}

	a := aggregator.NewAggregator(c.Collection)

	if err := a.Aggregate(); err != nil {
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
	if err := csvexport.WriteGamePhaseDurationStatsAggregated(a, "game-phase-duration-stats-aggregated.csv"); err != nil {
		log.Fatal(err)
	}
	if err := csvexport.WriteGameEventDurationStats(&a.Collection, "game-event-duration-stats.csv"); err != nil {
		log.Fatal(err)
	}
	if err := csvexport.WriteGameEventDurationStatsAggregated(a, "game-event-duration-stats-aggregated.csv"); err != nil {
		log.Fatal(err)
	}
}

func usage() {
	_, err := fmt.Fprintln(os.Stderr, "Pass one or more log files that should be processed.")
	if err != nil {
		fmt.Println("Pass one or more log files that should be processed.")
	}
	flag.PrintDefaults()
}
