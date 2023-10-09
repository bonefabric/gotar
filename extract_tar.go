package main

import (
	"archive/tar"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func extractTar(tarFilepath, extractFilepath string, warns chan<- error, readed chan<- *archiveInfo) error {
	defer close(warns)
	defer close(readed)

	absTar, err := filepath.Abs(tarFilepath)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to find absolute file path for %s, error: %s", tarFilepath, err))
	}

	absExtract, err := filepath.Abs(extractFilepath)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to find absolute file path for %s, error: %s", extractFilepath, err))
	}

	srcArchive, err := os.Open(absTar)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to open archive, error: %s", err))
	}

	defer func(absTar *os.File) {
		if err = srcArchive.Close(); err != nil {
			warns <- errors.New(fmt.Sprintf("failed to close archive %s, error: %s", tarFilepath, err))
		}
	}(srcArchive)

	tr := tar.NewReader(srcArchive)

	for {
		hdr, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			warns <- errors.New(fmt.Sprintf("failed to extract file header, error: %s", err))
			continue
		}

		extractPath := filepath.Join(absExtract, hdr.Name)
		if hdr.FileInfo().IsDir() {
			if err = os.MkdirAll(extractPath, hdr.FileInfo().Mode()); err != nil {
				warns <- errors.New(fmt.Sprintf("failed to create directory %s, skiped, error: %s", extractPath, err))
				continue
			}
		} else {
			dst, err := os.OpenFile(extractPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL|os.O_TRUNC, hdr.FileInfo().Mode())
			if err != nil {
				warns <- errors.New(fmt.Sprintf("failed to create file %s, skiped, error: %s", extractPath, err))
				continue
			}
			readSize, err := io.Copy(dst, tr)
			if err != nil {
				warns <- errors.New(fmt.Sprintf("failed to read file from archive %s, skiped, error: %s", extractPath, err))
				continue
			}

			readed <- &archiveInfo{
				filename: extractPath,
				size:     readSize,
			}
		}
	}
	return nil
}
