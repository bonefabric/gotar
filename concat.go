package main

import (
	"archive/tar"
	"flag"
	"io/fs"
	"log"
	"os"
)

func concat() {
	if len(flag.Args()) == 0 {
		log.Fatal("no one file specified")
	}

	archive, err := os.OpenFile(fileFlag, os.O_APPEND, fs.ModePerm)
	if err != nil {
		log.Fatalf("failed to open %s, error: %s\n", fileFlag, err)
	}

	defer archive.Close()

	tw := tar.NewWriter(archive)
	defer tw.Close()
	
}
