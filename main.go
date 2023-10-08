package main

import (
	"archive/tar"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func main() {
	configureLogger()

	if len(flag.Args()) == 0 {
		log.Fatal("No files to achieve")
	}

	makeArchive(flag.Args())
}

func makeArchive(sources []string) {
	abs, err := filepath.Abs(fileFlag)
	if err != nil {
		log.Fatalf("Failed to find absolute archive file path to %s, error: %s", fileFlag, err)
	}

	tarFile, err := os.Create(abs)
	if err != nil {
		log.Fatalf("Failed to create archive, error: %s", err)
	}

	defer func(tarFile *os.File) {
		if err := tarFile.Close(); err != nil {
			log.Printf("Failed to close arhive, error: %s\n", err)
		}
	}(tarFile)

	tw := tar.NewWriter(tarFile)
	defer func(tw *tar.Writer) {
		if err := tw.Close(); err != nil {
			log.Printf("Failed to close tar writer, error: %s", err)
		}
	}(tw)

	for _, srcName := range sources {
		absSrc, err := filepath.Abs(srcName)
		if err != nil {
			log.Printf("Failed to find absolute path to file %s, error: %s. Skiped\n", absSrc, err)
			continue
		}

		if err = writeFile(absSrc, tw); err != nil {
			log.Printf("Failed to write file %s, error: %s. Skiped\n", srcName, err)
			continue
		}
	}
}

func writeFile(sourcePath string, writer *tar.Writer) error {
	src, err := os.Open(sourcePath)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to open file, err: %s", err))
	}
	defer src.Close()

	stat, err := src.Stat()
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to read file info, err: %s", err))
	}

	hdr, err := tar.FileInfoHeader(stat, stat.Name())
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to create file header, err: %s", err))
	}

	if err = writer.WriteHeader(hdr); err != nil {
		return errors.New(fmt.Sprintf("Failed to write file header, err: %s", err))
	}

	_, err = io.Copy(writer, src)
	return err
}
