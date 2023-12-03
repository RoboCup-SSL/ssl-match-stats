package main

import (
	"flag"
	"fmt"
	"github.com/RoboCup-SSL/ssl-match-stats/pkg/matchstats"
	"log"
	"os"
	"path/filepath"
	"sync"
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

	var ch = make(chan string, *parallel)
	var wg sync.WaitGroup
	wg.Add(*parallel)

	for i := 0; i < *parallel; i++ {
		go func() {
			for {
				filename, ok := <-ch
				if !ok {
					wg.Done()
					return
				}
				process(filename)
				log.Println("Done with ", filename)
			}
		}()
	}

	log.Println("Starting")
	for _, filename := range args {
		log.Printf("Adding %v to queue", filename)
		ch <- filename
	}

	close(ch)
	wg.Wait()
	log.Println("Done")
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

	if err := os.MkdirAll(*targetDir, os.ModePerm); err != nil {
		log.Fatalln("Could not create directory for binary output file", err)
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
