package main

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func makeZip() error {
	zipf, err := os.Create(fmt.Sprintf("%s/export.zip", outputDir))
	if err != nil {
		return err
	}
	defer func() { _ = zipf.Close() }()

	archive := zip.NewWriter(zipf)
	defer func() { _ = archive.Close() }()
	//todo make concurrently
	err = filepath.WalkDir(outputDir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.Name() == ".gitignore" || d.IsDir() || path.Ext(p) == ".zip" {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return fmt.Errorf("can't get fileinfo %s", err)
		}
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return fmt.Errorf("can't create fileinfo header %s", err)
		}

		header.Name = path.Join(strings.TrimPrefix(p, outputDir))
		header.Method = zip.Deflate
		zwrt, err := archive.CreateHeader(header)
		if err != nil {
			return fmt.Errorf("can't create zip header %s", err)
		}
		f, err := os.Open(p)
		if err != nil {
			return err
		}
		defer func() { _ = f.Close() }()

		_, err = io.Copy(zwrt, f)
		if err != nil {
			return err
		}
		return err
	})
	//todo remove files out of the archive
	return err
}
