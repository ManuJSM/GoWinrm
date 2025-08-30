package msg

import (
	"fmt"

	"github.com/ManuJSM/GoWinrm/internal/utils"

	"github.com/ManuJSM/GoWinrm/internal/wsmv"
)

const terminateSignal = "http://schemas.microsoft.com/wbem/wsman/1/windows/shell/signal/terminate"

type CleanCommand struct {
	sessionOpts wsmv.SessionOptions
	shellID     string
	commandID   string
	shellURI    string
}

func NewCleanCommand(sessionOpts wsmv.SessionOptions, cmdOpts utils.Arguments) *CleanCommand {

	cCmd := &CleanCommand{
		sessionOpts: sessionOpts,
		shellID:     cmdOpts["shell_id"].(string),
		commandID:   cmdOpts["command_id"].(string),
		shellURI:    utils.OptOrDefault(cmdOpts, "shell_uri", wsmv.RESOURCEURICMD).(string),
	}

	return cCmd
}

func (c *CleanCommand) Headers() map[string]any {
	headers := wsmv.MergeHeaders(
		wsmv.SharedHeaders(c.sessionOpts),
		wsmv.ResourceURIShell(c.shellURI),
		wsmv.ActionSignal(),
		wsmv.SelectorShellID(c.shellID),
	)

	return headers
}

func (c *CleanCommand) Body() map[string]any {

	body := map[string]any{
		fmt.Sprintf("%s:Signal", wsmv.NS_WIN_SHELL): map[string]any{
			fmt.Sprintf("%s:Code", wsmv.NS_WIN_SHELL): terminateSignal,
			":attributes!": map[string]any{
				"CommandId": c.commandID,
			},
		},
	}

	return body
}
