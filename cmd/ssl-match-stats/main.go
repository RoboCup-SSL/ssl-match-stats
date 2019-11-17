package main

import (
	"flag"
	"fmt"
	"github.com/RoboCup-SSL/ssl-match-stats/pkg/matchstats"
	"log"
	"os"
)

func main() {
	flag.Usage = usage

	flag.Parse()

	args := flag.Args()

	if len(args) == 0 {
		usage()
		return
	}

	a := matchstats.NewCollector()
	for _, filename := range args {
		log.Println("Processing", filename)

		err := a.Process(filename)
		if err != nil {
			log.Printf("%v: %v\n", filename, err)
		} else {
			log.Printf("Processed %v\n", filename)
		}
	}

	if err := a.WriteBin("match-stats.bin"); err != nil {
		log.Println("Could not write binary file", err)
	}

	if err := a.WriteJson("match-stats.json"); err != nil {
		log.Println("Could not write JSON file", err)
	}
}

func usage() {
	_, err := fmt.Fprintln(os.Stderr, "Pass one or more log files that should be processed.")
	if err != nil {
		fmt.Println("Pass one or more log files that should be processed.")
	}
	flag.PrintDefaults()
}
