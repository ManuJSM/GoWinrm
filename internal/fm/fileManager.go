package fm

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"

	"github.com/ManuJSM/GoWinrm/internal/log"
	"github.com/ManuJSM/GoWinrm/internal/shell"
	"github.com/ManuJSM/GoWinrm/internal/transport"
	"github.com/ManuJSM/GoWinrm/internal/wsmv"
)

const chunkSize = 5900

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

	// Preparar archivos para transferencia
	// files, err := ft.prepareFiles(localPaths, remotePath)
	// if err != nil {
	// 	return nil, fmt.Errorf("error preparing files: %w", err)
	// }

	// Verificar archivos existentes en destino
	// err = ft.checkRemoteFiles(files)
	// if err != nil {
	// 	return nil, fmt.Errorf("error checking remote files: %w", err)
	// }

	// Transferir archivos
	err := fm.transferFile(File{SrcPath: localpath, DstPath: remotepath, Utd: false})
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

func (fm *FileManager) transferFile(file File) error {

	if !file.Utd {
		log.Debug(fmt.Sprintf("Skipping %s (already up to date)", file.SrcPath))
	}

	bytesTransferred, err := fm.streamUpload(file.SrcPath, file.DstPath)
	if err != nil {
		return fmt.Errorf("error uploading %s: %v", file.SrcPath, err)
	}

	log.Debug(fmt.Sprintf("Uploaded %s (%d bytes)", file.SrcPath, bytesTransferred))

	return nil
}

func (fm *FileManager) streamUpload(src string, dst string) (int64, error) {
	var sendBytes int64

	file, err := os.Open(src)
	if err != nil {
		return sendBytes, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	buffer := make([]byte, chunkSize)

	initScript := fmt.Sprintf(`$to = "%s";$parent = Split-Path $to;if(!(Test-Path $parent)) { New-Item -ItemType Directory -Path $parent | Out-Null }`, dst)

	output, err := fm.shell.RunCommand(initScript)
	if err != nil || output.ExitCode != 0 {
		return sendBytes, fmt.Errorf("error creando archivo remoto")
	}

	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return sendBytes, fmt.Errorf("error reading file: %w", err)
		}

		encoded := base64.StdEncoding.EncodeToString(buffer[:n])

		writeScript := fmt.Sprintf(`$bytes = [Convert]::FromBase64String("%s");$fs = [System.IO.File]::Open('%s', 'Append', 'Write', 'ReadWrite');$fs.Write($bytes, 0, $bytes.Length);$fs.Close()`, encoded, dst)

		output, err := fm.shell.RunCommand(writeScript)
		if err != nil || output.ExitCode != 0 {
			return sendBytes, fmt.Errorf("error escribiendo en el archivo, %v", output.Stderr.String())

		}

		sendBytes += int64(n)
	}

	return sendBytes, nil
}

func NewFileManager(transport transport.Transport, opt *wsmv.SessionOptions) *FileManager {
	return &FileManager{
		shell: shell.NewPsShell(transport, opt),
	}
}
