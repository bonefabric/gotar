package main

import "flag"

var (
	fileFlag    string
	createFlag  bool
	verboseFlag bool
)

func init() {
	flag.StringVar(&fileFlag, "f", "archive.tar", "Archive file name")
	flag.BoolVar(&createFlag, "c", true, "Create archive")
	flag.BoolVar(&verboseFlag, "v", false, "Verbosely list files processed")

	flag.Parse()
}
