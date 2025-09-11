package shell

import (
	"strings"

	"github.com/ManuJSM/GoWinrm/internal/transport"
	"github.com/ManuJSM/GoWinrm/internal/utils"
	"github.com/ManuJSM/GoWinrm/internal/wsmv"
)

const PsPath = "C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe"

type Ps struct {
	cmd Shell
}

func NewPsShell(transport transport.Transport, opts *wsmv.SessionOptions) *Ps {
	return &Ps{
		cmd: NewCmdShell(transport, opts),
	}

}
func (ps *Ps) RunCommand(psCommand string, arguments ...string) (*utils.Output, error) {

	command := psCommand + " " + strings.Join(arguments, " ")

	escapedQuotes := strings.ReplaceAll(command, `"`, `\"`)

	fullCommand := PsPath + ` -Command "` + escapedQuotes + `"`

	return ps.cmd.RunCommand(fullCommand)
}

func (ps *Ps) Close() error {
	return ps.cmd.Close()
}
