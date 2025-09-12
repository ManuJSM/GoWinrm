package fm

import (
	"fmt"
	"os"
	"strconv"

	"github.com/ManuJSM/GoWinrm/internal/shell"
	"github.com/ManuJSM/GoWinrm/internal/transport"
	"github.com/ManuJSM/GoWinrm/internal/wsmv"
)

const chunkSize = 5800

const (
	tempPath       = `C:\Windows\Temp`
	zipTempPath    = tempPath + "\\temp.zip"
	base64TempPath = tempPath + "\\temp.b64"
	binTempPath    = tempPath + "\\temp.bin"
)

type FileManager struct {
	shell    shell.Shell
	uploader *Uploader
}
type File struct {
	SrcPath string
	DstPath string
	Utd     bool
	Size    int64
}

func (fm *FileManager) Close() error {
	return fm.shell.Close()
}

func NewFileManager(transport transport.Transport, opt *wsmv.SessionOptions) *FileManager {
	shell := shell.NewCmdShell(transport, opt)
	uploader := NewUploader(shell)

	return &FileManager{
		shell:    shell,
		uploader: uploader,
	}
}

func (fm *FileManager) UploadFile(localpath, dst string) error {

	return fm.uploader.UploadFile(localpath, dst)
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
