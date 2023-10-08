package main

import (
	"archive/tar"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"slices"
)

func main() {
	configureLogger()

	if len(flag.Args()) == 0 {
		log.Fatal("No files to achieve")
	}

	makeArchive(flag.Args())
}

func makeArchive(sources []string) {
	tarFile := makeTarFile(fileFlag)

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

	skips := skipFor(fileFlag)

	for _, srcName := range sources {
		if err := writeFile(srcName, tw, skips); err != nil {
			log.Printf("Failed to write file %s, error: %s. Skiped\n", srcName, err)
			continue
		}
	}
}

func makeTarFile(path string) *os.File {
	absTar, err := filepath.Abs(path)
	if err != nil {
		log.Fatalf("Failed to find absolute archive file path to %s, error: %s", path, err)
	}

	tarFile, err := os.Create(absTar)
	if err != nil {
		log.Fatalf("Failed to create archive, error: %s", err)
	}
	return tarFile
}

func writeFile(sourcePath string, writer *tar.Writer, skip []string) error {
	absRootPath, err := filepath.Abs(sourcePath)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to find absolute path to file %s, error: %s", sourcePath, err))
	}

	walker := func(path string, d fs.DirEntry, err error) error {
		absSrcPath, err := filepath.Abs(path)
		if err != nil {
			return errors.New(fmt.Sprintf("Failed to find absolute path to file %s, error: %s", path, err))
		}

		if slices.Contains(skip, absSrcPath) {
			return nil
		}

		src, err := os.Open(absSrcPath)
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

		if hdr.Name, err = filepath.Rel(filepath.Dir(absRootPath), absSrcPath); err != nil {
			return errors.New(fmt.Sprintf("Failed to find relative path for file %s, err: %s", absSrcPath, err))
		}

		if err = writer.WriteHeader(hdr); err != nil {
			return errors.New(fmt.Sprintf("Failed to write file header, err: %s", err))
		}

		if stat.IsDir() {
			return nil
		}

		_, err = io.Copy(writer, src)
		return err
	}

	if err = filepath.WalkDir(absRootPath, walker); err != nil {
		return errors.New(fmt.Sprintf("Failed to walk directories: %s", err))
	}
	return nil
}

func skipFor(archivePath string) []string {
	absArchive, err := filepath.Abs(archivePath)
	if err != nil {
		log.Fatalf("Failed to find absolute file path to archive %s, error: %s", archivePath, err)
	}

	return []string{
		absArchive,
	}
}
