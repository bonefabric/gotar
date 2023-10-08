package main

import (
	"archive/tar"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
)

type tarWalker struct {
	rootAbs string
	skipFor []string
	writer  *tar.Writer
	warns   chan<- error
	written chan<- *writeInfo
}

type writeInfo struct {
	filename string
	size     int64
}

func (walker *tarWalker) walkerFunc(path string, _ fs.DirEntry, err error) error {
	absSrcPath, err := filepath.Abs(path)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to find absolute path to file %s, error: %s", path, err))
	}

	if slices.Contains(walker.skipFor, absSrcPath) {
		return nil
	}

	src, err := os.Open(absSrcPath)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to open file, err: %s", err))
	}
	defer func(src *os.File) {
		if err = src.Close(); err != nil {
			walker.warns <- errors.New(fmt.Sprintf("failed to close file %s, error: %s", path, err))
		}
	}(src)

	link, err := filepath.Rel(filepath.Dir(walker.rootAbs), absSrcPath)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get related path to %s, error: %s", path, err))
	}

	written, err := walker.fileToTar(src, link, walker.writer)

	walker.written <- &writeInfo{
		filename: path,
		size:     written,
	}

	return err
}

// fileToTar writes file data to tar writer
func (walker *tarWalker) fileToTar(file *os.File, link string, writer *tar.Writer) (int64, error) {
	fstat, err := file.Stat()
	if err != nil {
		return 0, errors.New(fmt.Sprintf("failed to get file stat, err: %s", err))
	}

	hdr, err := tar.FileInfoHeader(fstat, link)
	if err != nil {
		return 0, errors.New(fmt.Sprintf("failed to create file info header, err: %s", err))
	}

	hdr.Name = link

	if err = writer.WriteHeader(hdr); err != nil {
		return 0, errors.New(fmt.Sprintf("failed to write file info header: err: %s", err))
	}

	if fstat.IsDir() {
		return 0, nil
	}

	return io.Copy(writer, file)
}

func makeTarArchive(tarFilePath string, sources []string, warns chan<- error, written chan<- *writeInfo) error {
	absTarPath, err := filepath.Abs(tarFilePath)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to find absolute path to %s, error: %s", tarFilePath, err))
	}

	tarFile, err := os.Create(absTarPath)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to create archive file, error: %s", err))
	}

	defer func(tarFile *os.File) {
		if err = tarFile.Close(); err != nil {
			warns <- errors.New(fmt.Sprintf("failed to close arhive, error: %s\n", err))
		}
	}(tarFile)

	tw := tar.NewWriter(tarFile)
	defer func(tw *tar.Writer) {
		if err = tw.Close(); err != nil {
			warns <- errors.New(fmt.Sprintf("failed to close tar writer, error: %s", err))
		}
	}(tw)

	skips := []string{
		absTarPath,
	}

	for _, srcName := range sources {
		if err = pathToTar(srcName, tw, skips, warns, written); err != nil {
			warns <- errors.New(fmt.Sprintf("Failed to write file %s, error: %s. Skiped\n", srcName, err))
			continue
		}
	}

	close(warns)
	close(written)

	return nil
}

// pathToTar writes path files to tar writer
func pathToTar(sourcePath string, writer *tar.Writer, skip []string, warns chan<- error, written chan<- *writeInfo) error {
	absRootPath, err := filepath.Abs(sourcePath)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to find absolute path to file %s, error: %s", sourcePath, err))
	}

	walker := &tarWalker{
		rootAbs: absRootPath,
		skipFor: skip,
		writer:  writer,
		warns:   warns,
		written: written,
	}

	if err = filepath.WalkDir(absRootPath, walker.walkerFunc); err != nil {
		return errors.New(fmt.Sprintf("failed to walk directories: %s", err))
	}
	return nil
}
