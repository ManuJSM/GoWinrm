package msg

import (
	"GoWinrm/internal/utils"
	"GoWinrm/internal/wsmv"
	"fmt"
)

type OutputCommand struct {
	sessionOpts wsmv.SessionOptions
	shellURI    string
	shellID     string
	commandID   string
	outStreams  string
}

func NewOutputCommand(sessionOpts wsmv.SessionOptions, shellOpts utils.Arguments) *OutputCommand {

	return &OutputCommand{
		sessionOpts: sessionOpts,
		shellURI:    utils.OptOrDefault(shellOpts, "shell_uri", wsmv.RESOURCEURICMD).(string),
		shellID:     shellOpts["shell_id"].(string),
		commandID:   shellOpts["command_id"].(string),
		outStreams:  utils.OptOrDefault(shellOpts, "o_streams", "stdout stderr").(string),
	}
}
func (oc *OutputCommand) headerOpts() map[string]any {
	options := []map[string]any{
		{
			"_": "TRUE",
			":attributes!": map[string]any{
				"Name": "WSMAN_CMDSHELL_OPTION_KEEPALIVE",
			},
		},
	}
	return map[string]any{
		fmt.Sprintf("%s:OptionSet", wsmv.NS_WSMAN_DMTF): map[string]any{
			fmt.Sprintf("%s:Option", wsmv.NS_WSMAN_DMTF): options,
		},
	}
}
func (oc *OutputCommand) Headers() map[string]any {
	return wsmv.MergeHeaders(
		wsmv.SharedHeaders(oc.sessionOpts),
		wsmv.ResourceURIShell(oc.shellURI),
		wsmv.ActionReceive(),
		oc.headerOpts(),
		wsmv.SelectorShellID(oc.shellID),
	)
}

func (oc *OutputCommand) outputBody() map[string]any {

	return map[string]any{
		fmt.Sprintf("%s:DesiredStream", wsmv.NS_WIN_SHELL): map[string]any{
			"_": oc.outStreams,
			":attributes!": map[string]any{
				"CommandId": oc.commandID,
			},
		},
	}
}
func (oc *OutputCommand) Body() map[string]any {

	return map[string]any{
		fmt.Sprintf("%s:Receive", wsmv.NS_WIN_SHELL): oc.outputBody(),
	}

}
