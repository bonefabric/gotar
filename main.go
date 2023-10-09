package main

import (
	"flag"
	"log"
)

func main() {
	setupLogger()
	validateFlags()

	if createFlag {
		makeArchive(fileFlag, flag.Args(), verboseFlag)
		return
	}

	if extractFlag {
		extractArchive(fileFlag, flag.Args()[0], verboseFlag)
	}
}

func makeArchive(filepath string, sources []string, verbose bool) {
	if err := makeTarArchive(filepath, sources, warnChan(), verboseChan(verbose)); err != nil {
		log.Fatalf("failed to create archive: %s\n", err)
	}
}

func extractArchive(archivePath, extractPath string, verbose bool) {
	if err := extractTar(archivePath, extractPath, warnChan(), verboseChan(verbose)); err != nil {
		log.Fatalf("failed to extract archive: %s\n", err)
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

func verboseChan(show bool) chan<- *archiveInfo {
	writes := make(chan *archiveInfo, 1)

	go func() {
		for write := range writes {
			if show {
				log.Printf("file %s archived %d bytes\n", write.filename, write.size)
			}
		}
	}()
	return writes
}
