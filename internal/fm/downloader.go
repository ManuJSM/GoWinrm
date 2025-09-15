package fm

import (
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ManuJSM/GoWinrm/internal/shell"
	"github.com/ManuJSM/GoWinrm/internal/utils"
)

const chunkDownloadSize = 48024
const chunkPath = tempPath + "\\chunks_base64"

var localZipTempFile = os.TempDir() + "\\temp.zip"

type Downloader struct {
	shell shell.Shell
}

func NewDownloader(shell shell.Shell) *Downloader {
	return &Downloader{
		shell: shell,
	}
}

// Devuelve el numero de chunks a leer
func (d *Downloader) prepareDownload(srcPath string) (int, error) {
	//Create Remote zip
	createRemoteZip := fmt.Sprintf(`%s -Command "rm %s;Compress-Archive -Path %s -DestinationPath %s"`,
		shell.PsPath, remoteZipTempFile, srcPath, remoteZipTempFile)

	output, err := d.shell.RunCommand(createRemoteZip)
	if err != nil || output.ExitCode == 1 {
		return 0, fmt.Errorf("error creating Zip Remote File%v", output.Stderr.String())
	}

	//Create Chunks of the zip
	chunkDivider := fmt.Sprintf(`%s -Command "$chunkSize=%d; $in='%s'; $out='%s'; if (-not (Test-Path $out)) {New-Item -ItemType Directory $out | Out-Null}; $stream=[IO.File]::OpenRead($in); $i=0; $buffer=New-Object byte[] $chunkSize; while (($read=$stream.Read($buffer,0,$chunkSize)) -gt 0) { $chunk = if ($read -eq $chunkSize) { $buffer } else { $buffer[0..($read-1)] }; [IO.File]::WriteAllText(\"$out/chunk$i.b64\", [Convert]::ToBase64String($chunk)); $i++ }; $stream.Close();echo $i"`,
		shell.PsPath, chunkDownloadSize, remoteZipTempFile, chunkPath)

	output, err = d.shell.RunCommand(chunkDivider)
	if err != nil || output.ExitCode == 1 {
		return 0, fmt.Errorf("error creating Remote Chunks%v", output.Stderr.String())
	}

	clean := strings.TrimSpace(output.Stdout.String())
	return strconv.Atoi(clean)

}

func (d *Downloader) readChunk(n int) ([]byte, error) {

	command := fmt.Sprintf("type %s\\chunk%d.b64", chunkPath, n)

	output, err := d.shell.RunCommand(command)
	if err != nil || output.ExitCode == 1 {
		return nil, fmt.Errorf("error downloading %v", output.Stderr.String())
	}

	return base64.StdEncoding.DecodeString(output.Stdout.String())
}

func (d *Downloader) writeZipTemp(chunks int) error {
	os.Remove(localZipTempFile)
	file, err := os.OpenFile(localZipTempFile, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	total := chunks - 1
	for i := range chunks {

		decodedBytes, err := d.readChunk(i)
		if err != nil {
			return err
		}
		_, err = file.Write(decodedBytes)
		if err != nil {
			return err
		}
		showProgress(int64(total), int64(i))

	}

	return nil
}

func (d *Downloader) cleanup() {

	d.shell.RunCommand(fmt.Sprintf(`rmdir /s /q "%s" &amp; del %s`, chunkPath, remoteZipTempFile))
	os.Remove(localZipTempFile)

}

func (d *Downloader) DownloadFile(remotePath, dst string) error {

	chunks, err := d.prepareDownload(remotePath)
	if err != nil {
		return err
	}

	err = d.writeZipTemp(chunks)
	if err != nil {
		return err
	}

	//unzip file
	err = utils.UnzipFile(localZipTempFile, dst)
	if err != nil {
		return err
	}

	return nil
}
