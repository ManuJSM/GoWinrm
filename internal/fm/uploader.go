package fm

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"os"

	"github.com/ManuJSM/GoWinrm/internal/shell"
	"github.com/ManuJSM/GoWinrm/internal/utils"
)

const chunkUploadSize = 5800

type Uploader struct {
	shell shell.Shell
}

func NewUploader(shell shell.Shell) *Uploader {
	return &Uploader{
		shell: shell,
	}
}

func (u *Uploader) UploadFile(localpath, dst string) error {

	//Preparar archivos para transferencia
	file, err := u.prepareFiles(localpath)
	if err != nil {
		return fmt.Errorf("error preparing files: %v", err)
	}

	// Transferir archivos
	err = u.streamUpload(file)
	if err != nil {
		return fmt.Errorf("error transferring files: %w", err)
	}

	// Extraer archivos comprimidos si es necesario
	command := fmt.Sprintf(`Expand-Archive -Path '%s' -DestinationPath '%s'`, remoteZipTempFile, dst)
	output, err := u.shell.RunCommand(shell.PsPath + " -NoProfile -Command " + command)
	if err != nil || output.ExitCode != 0 {
		return fmt.Errorf("error descomprimiendo, %v", output.Stderr.String())
	}

	return nil

}
func (u *Uploader) cleanup() {
	u.shell.RunCommand("del " + tempPath + "\\temp.*")
}

func (u *Uploader) prepareFiles(localpath string) ([]byte, error) {

	//Limpiar temps
	u.cleanup()

	output, err := u.shell.RunCommand(fmt.Sprintf("type nul>%s", remoteZipTempFile))
	if err != nil || output.ExitCode != 0 {
		return nil, fmt.Errorf("error preparando temps, %v", output.Stderr.String())
	}

	//Si ya es un zip no zipearlo
	if utils.IsCompressedFile(localpath) {

		data, err := os.ReadFile(localpath)
		if err != nil {
			return nil, err
		}
		return data, nil

	}
	return utils.ZipFileToMemory(localpath)

}

func (u *Uploader) streamUpload(buf []byte) error {
	var sendBytes int64
	totalSize := len(buf)

	buffer := make([]byte, chunkUploadSize)
	reader := bytes.NewReader(buf)

	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading file: %w", err)
		}

		encoded := base64.StdEncoding.EncodeToString(buffer[:n])

		writeScript := fmt.Sprintf("echo %s > %s &amp; certutil -decode %s %s &amp;copy /b %s+%s %s &amp;del %s", encoded, base64TempPath, base64TempPath, binTempPath, remoteZipTempFile, binTempPath, remoteZipTempFile, binTempPath)

		output, err := u.shell.RunCommand(writeScript)
		if err != nil || output.ExitCode != 0 {
			return fmt.Errorf("error escribiendo en el archivo, %v", output.Stderr.String())
		}

		sendBytes += int64(n)
		showProgress(int64(totalSize), sendBytes)
	}

	return nil
}
