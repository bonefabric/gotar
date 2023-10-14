package main

import (
	"flag"
	"log"
	"slices"
)

var (
	concatFlag  bool
	createFlag  bool
	diffFlag    bool
	deleteFlag  bool
	addFlag     bool
	showFlag    bool
	updateFlag  bool
	extractFlag bool

	timePreserveFlag bool
	fileFlag         string
	verboseFlag      bool
)

func init() {

	//Main options

	//todo add --catenate, --concatenate
	flag.BoolVar(&concatFlag, "A", false, "Add file to archive")
	//todo add --create
	flag.BoolVar(&createFlag, "c", false, "Create archive")
	//todo --diff, --compare
	flag.BoolVar(&diffFlag, "d", false, "Show difference between archive and filesystem")
	flag.BoolVar(&deleteFlag, "delete", false, "Delete files from archive")
	//todo add --append
	flag.BoolVar(&addFlag, "r", false, "Add files to the end to archive")
	//todo add --list
	flag.BoolVar(&showFlag, "t", false, "Show files in archive")
	//todo add --update
	flag.BoolVar(&updateFlag, "u", false, "Show files in archive")
	//todo --extract, --get
	flag.BoolVar(&extractFlag, "x", false, "Extract archive")

	//Additional options

	flag.BoolVar(&timePreserveFlag, "atime-preserve", false, "Do not change file access time")

	// todo add more flags

	flag.StringVar(&fileFlag, "f", "archive.tar", "Archive file name")
	flag.BoolVar(&verboseFlag, "v", false, "Verbosely list files processed")

	flag.Parse()
}

func findAction() (res func()) {
	actions := map[int]func(){
		0: concat,
		1: create,
		2: diff,
		3: delete,
		4: add,
		5: show,
		6: update,
		7: extract,
	}

	flags := [...]bool{concatFlag, createFlag, diffFlag, deleteFlag, addFlag, showFlag, updateFlag, extractFlag}

	for i, v := range flags {
		if v {
			res = actions[i]
			if slices.Contains(flags[i+1:], true) {
				log.Fatal("only one option must be specified")
			}
			return
		}
	}
	log.Fatal("one option must be specified")
	return
}
