// SPDX-FileCopyrightText: 2025 Margo
//
// SPDX-License-Identifier: MIT
//
// SPDX-FileContributor: Michael Adler <michael.adler@siemens.com>

package main

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func unpackTgz(src io.Reader, destDir string, skipHidden bool) error {
	gzr, err := gzip.NewReader(src)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			return err
		}
		if skipHidden && strings.HasPrefix(header.Name, ".") {
			log.Println("WARN: Skipping hidden entry", header.Name)
			continue
		}

		target := filepath.Join(destDir, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, os.FileMode(header.Mode)); err != nil {
				return err
			}
		case tar.TypeReg:
			file, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			defer file.Close()

			if _, err := io.Copy(file, tr); err != nil {
				return err
			}
		default:
			log.Println("WARN: Skipping unsupported file type", header.Typeflag)
		}
	}
	return nil
}

func findAppFiles(dir string) ([]string, error) {
	var appFiles []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".app") {
			appFiles = append(appFiles, path)
		}
		return nil
	})
	return appFiles, err
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}
