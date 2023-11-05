package main

import (
	"flag"
	"fmt"
	"github.com/RoboCup-SSL/ssl-match-stats/pkg/matchstats"
	"log"
	"os"
	"path/filepath"
)

var targetDir = flag.String("targetDir", "", "directory where the match stats should be written to")

func main() {
	flag.Usage = usage

	parallel := flag.Int("parallel", 1, "number of parallel processes")

	flag.Parse()

	args := flag.Args()

	if len(args) == 0 {
		usage()
		return
	}

	guard := make(chan struct{}, *parallel)
	for _, filename := range args {
		guard <- struct{}{}
		go func(filename string) {
			process(filename)
			<-guard
		}(filename)
	}
}

func process(filename string) {
	a := matchstats.NewCollector()
	log.Println("Processing", filename)

	err := a.Process(filename)
	if err != nil {
		log.Printf("%v: %v\n", filename, err)
	} else {
		log.Printf("Processed %v\n", filename)
	}

	baseFilename := filepath.Base(filename)
	if err := a.WriteBin(filepath.Join(*targetDir, baseFilename+".bin")); err != nil {
		log.Println("Could not write binary file", err)
	}

	if err := a.WriteJson(filepath.Join(*targetDir, baseFilename+".json")); err != nil {
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
