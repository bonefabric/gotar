package main

import (
	"flag"
	"log"
)

var (
	fileFlag    string
	createFlag  bool
	verboseFlag bool
	extractFlag bool
)

func init() {
	flag.StringVar(&fileFlag, "f", "archive.tar", "Archive file name")
	flag.BoolVar(&createFlag, "c", false, "Create archive")
	flag.BoolVar(&extractFlag, "x", false, "Extract archive")
	flag.BoolVar(&verboseFlag, "v", false, "Verbosely list files processed")

	flag.Parse()
}

func validateFlags() {
	if (createFlag && extractFlag) || (!createFlag && !extractFlag) {
		log.Fatal("specify one flag: -c to create or -x to extract the archive")
	}

	if createFlag && len(flag.Args()) == 0 {
		log.Fatal("specify paths to create the archive")
	}

	if extractFlag && len(flag.Args()) != 1 {
		log.Fatal("specify path to extract the archive")
	}
}
