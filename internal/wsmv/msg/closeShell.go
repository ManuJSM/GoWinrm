package msg

import (
	"GoWinrm/internal/utils"
	"GoWinrm/internal/wsmv"
)

type CloseShell struct {
	sessionOpts wsmv.SessionOptions
	shellID     string
	shellURI    string
}

func NewCloseShell(sessionOpts wsmv.SessionOptions, shellOpts utils.Arguments) *CloseShell {

	return &CloseShell{
		sessionOpts: sessionOpts,
		shellID:     shellOpts["shell_id"].(string),
		shellURI:    utils.OptOrDefault(shellOpts, "shell_uri", wsmv.RESOURCEURICMD).(string),
	}
}

func (c *CloseShell) Headers() map[string]any {
	return wsmv.MergeHeaders(
		wsmv.SharedHeaders(c.sessionOpts),
		wsmv.ActionDelete(),
		wsmv.ResourceURIShell(c.shellURI),
		wsmv.SelectorShellID(c.shellID),
	)
}

func (c *CloseShell) Body() map[string]any {
	return nil
}
