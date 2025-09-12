package utils

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func ZipFileToMemory(sourceDir string) ([]byte, error) {
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	absSource, err := filepath.Abs(sourceDir)
	if err != nil {
		return nil, fmt.Errorf("no se pudo obtener ruta absoluta: %w", err)
	}
	baseName := filepath.Base(absSource)

	err = filepath.Walk(absSource, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(absSource, path)
		if err != nil {
			return err
		}

		zipPath := filepath.Join(baseName, relPath)
		zipPath = filepath.ToSlash(zipPath)

		if info.IsDir() {
			if relPath == "." {
				return nil
			}
			_, err := zipWriter.Create(zipPath + "/")
			return err
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		header := &zip.FileHeader{
			Name:   zipPath,
			Method: zip.Deflate,
		}

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		_, err = io.Copy(writer, file)
		return err
	})

	if err != nil {
		zipWriter.Close()
		return nil, fmt.Errorf("error al crear zip: %w", err)
	}

	if err := zipWriter.Close(); err != nil {
		return nil, fmt.Errorf("error al cerrar zip: %w", err)
	}
	return buf.Bytes(), nil
}
