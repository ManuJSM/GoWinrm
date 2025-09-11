package utils

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func ZipFileToMemory(filePath string) ([]byte, error) {
	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)

	// Abrir archivo fuente
	srcFile, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("no se pudo abrir el archivo: %w", err)
	}
	defer srcFile.Close()

	// Obtener info del archivo para el header
	info, err := srcFile.Stat()
	if err != nil {
		return nil, fmt.Errorf("no se pudo obtener información del archivo: %w", err)
	}

	// Crear encabezado zip
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return nil, fmt.Errorf("no se pudo crear encabezado zip: %w", err)
	}
	header.Name = filepath.Base(filePath)
	header.Method = zip.Deflate // Comprimir

	// Crear archivo dentro del zip
	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return nil, fmt.Errorf("no se pudo crear entrada zip: %w", err)
	}

	// Copiar contenido del archivo original al zip
	if _, err := io.Copy(writer, srcFile); err != nil {
		return nil, fmt.Errorf("no se pudo copiar contenido al zip: %w", err)
	}

	// Finalizar el zip
	if err := zipWriter.Close(); err != nil {
		return nil, fmt.Errorf("error al cerrar zip: %w", err)
	}

	// Devolver el zip como bytes
	return buf.Bytes(), nil
}
func ReadFileToBuffer(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var buffer []byte
	temp := make([]byte, 1024)

	for {
		n, err := file.Read(temp)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if n == 0 {
			break
		}
		buffer = append(buffer, temp[:n]...)
	}

	return buffer, nil
}
