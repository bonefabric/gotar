package main

import (
	"flag"
	"log"
)

func main() {
	setupLogger()

	if createFlag {
		makeArchive(flag.Args(), verboseFlag)
	}
}

func makeArchive(sources []string, verbose bool) {
	if err := makeTarArchive(fileFlag, sources, warnChan(), verboseChan(verbose)); err != nil {
		log.Fatalf("failed to create tar archive: %s\n", err)
	}
}

func warnChan() chan<- error {
	warns := make(chan error, 1)

	go func() {
		for warn := range warns {
			log.Printf("warning: %s\n", warn)
		}
	}()
	return warns
}

func verboseChan(show bool) chan<- *writeInfo {
	writes := make(chan *writeInfo, 1)

	go func() {
		for write := range writes {
			if show {
				log.Printf("file %s archived %d bytes\n", write.filename, write.size)
			}
		}
	}()
	return writes
}
