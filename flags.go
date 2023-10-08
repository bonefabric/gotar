package main

import "flag"

var fileFlag string

func init() {
	flag.StringVar(&fileFlag, "f", "archive.tar", "Archive file name")

	flag.Parse()
}
