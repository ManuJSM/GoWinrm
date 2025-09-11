package fm

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/ManuJSM/GoWinrm/internal/shell"
	"github.com/ManuJSM/GoWinrm/internal/transport"
	"github.com/ManuJSM/GoWinrm/internal/utils"
	"github.com/ManuJSM/GoWinrm/internal/wsmv"
)

const chunkSize = 5800

var (
	tempPath       = os.TempDir()
	zipTempPath    = tempPath + "\\temp.zip"
	base64TempPath = tempPath + "\\temp.b64"
	binTempPath    = tempPath + "\\temp.bin"
)

type FileManager struct {
	shell shell.Shell
}
type File struct {
	SrcPath string
	DstPath string
	Utd     bool
	Size    int64
}

func (fm *FileManager) UploadFile(localpath, remotepath string) error {

	//Preparar archivos para transferencia
	file, err := fm.prepareFiles(localpath)
	if err != nil {
		return fmt.Errorf("error preparing files: %v", err)
	}

	// Verificar archivos existentes en destino
	// err = ft.checkRemoteFiles(files)
	// if err != nil {
	// 	return nil, fmt.Errorf("error checking remote files: %w", err)
	// }

	// Transferir archivos
	err = fm.streamUpload(file, remotepath)
	if err != nil {
		return fmt.Errorf("error transferring files: %w", err)
	}

	// Extraer archivos comprimidos si es necesario
	// err = ft.extractFiles(files)
	// if err != nil {
	// 	return nil, fmt.Errorf("error extracting files: %w", err)
	// }

	// Limpiar archivos temporales
	// ft.cleanup(files)

	// ft.logger.Debugf("Uploaded %d files (%d bytes) in %v",
	// 	len(files), totalBytes, duration)

	// return &UploadResult{
	// 	TotalBytes: totalBytes,
	// 	TotalFiles: len(files),
	// 	Duration:   duration,
	// }, nil
	return nil

}

func (fm *FileManager) prepareFiles(localpath string) ([]byte, error) {
	//TODO soporte para carpetas?

	return utils.ZipFileToMemory(localpath)
}

func (fm *FileManager) streamUpload(buf []byte, dst string) error {
	var sendBytes int64
	totalSize := len(buf)

	buffer := make([]byte, chunkSize)
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

		writeScript := fmt.Sprintf("echo %s > %s &amp; certutil -decode %s %s &amp; type %s >> %s &amp;del %s", encoded, base64TempPath, base64TempPath, binTempPath, binTempPath, zipTempPath, binTempPath)

		output, err := fm.shell.RunCommand(writeScript)
		if err != nil || output.ExitCode != 0 {
			return fmt.Errorf("error escribiendo en el archivo, %v", output.Stderr.String())
		}

		sendBytes += int64(n)
		showProgress(int64(totalSize), sendBytes)
	}

	command := fmt.Sprintf(`Expand-Archive -Path '%s' -DestinationPath '%s'`, zipTempPath, dst)
	output, err := fm.shell.RunCommand(shell.PsPath + " -NoProfile -Command " + command)
	if err != nil || output.ExitCode != 0 {
		return fmt.Errorf("error descomprimiendo, %v", output.Stderr.String())
	}
	output, err = fm.shell.RunCommand("del " + tempPath + "\\temp.*")
	if err != nil || output.ExitCode != 0 {
		return fmt.Errorf("error descomprimiendo, %v", output.Stderr.String())
	}

	return nil
}

func NewFileManager(transport transport.Transport, opt *wsmv.SessionOptions) *FileManager {
	return &FileManager{
		shell: shell.NewCmdShell(transport, opt),
	}
}

func showProgress(totalSize, sendBytes int64) {
	progress := (sendBytes * 100) / totalSize

	os.Stdout.WriteString("\r")
	os.Stdout.WriteString("[*] ")
	os.Stdout.WriteString(strconv.Itoa(int(progress)))
	os.Stdout.WriteString("%")
	if progress >= 100 {
		fmt.Println()
	}

}
