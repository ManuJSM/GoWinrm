package msg

import (
	"fmt"

	"github.com/ManuJSM/GoWinrm/internal/utils"
	"github.com/ManuJSM/GoWinrm/internal/wsmv"
)

// Command representa un mensaje WSMV para ejecutar un comando dentro de una shell remota
type Command struct {
	sessionOpts  wsmv.SessionOptions
	shellID      string
	command      string
	arguments    []string
	shellURI     string
	consoleMode  string
	skipCmdShell string
}

// NewCommand crea una nueva instancia de Command
func NewCommand(sessionOpts wsmv.SessionOptions, cmdOpts utils.Arguments) *Command {

	cmd := &Command{
		sessionOpts:  sessionOpts,
		shellID:      cmdOpts["shell_id"].(string),
		command:      cmdOpts["command"].(string),
		arguments:    utils.OptStringSlice(cmdOpts, "arguments", []string{}),
		shellURI:     utils.OptOrDefault(cmdOpts, "shell_uri", wsmv.RESOURCEURICMD).(string),
		consoleMode:  utils.OptOrDefault(cmdOpts, "console_mode_stdin", "TRUE").(string),
		skipCmdShell: utils.OptOrDefault(cmdOpts, "skip_cmd_shell", "FALSE").(string),
	}

	return cmd
}

// Headers genera los encabezados para el mensaje Command
func (c *Command) Headers() map[string]any {
	headers := wsmv.MergeHeaders(
		wsmv.SharedHeaders(c.sessionOpts),
		wsmv.ResourceURIShell(c.shellURI),
		wsmv.ActionCommand(),
		wsmv.SelectorShellID(c.shellID),
	)

	// Solo se agrega OptionSet si se trata de la shell por defecto
	if c.shellURI == wsmv.RESOURCEURICMD {
		headers = wsmv.MergeHeaders(headers, c.commandHeaderOpts())
	}

	return headers
}

// Body genera el cuerpo del mensaje Command
func (c *Command) Body() map[string]any {
	body := map[string]any{
		fmt.Sprintf("%s:CommandLine", wsmv.NS_WIN_SHELL): c.commandBody(),
	}
	return body
}

func (c *Command) commandBody() map[string]any {

	body := map[string]any{
		fmt.Sprintf("%s:Command", wsmv.NS_WIN_SHELL): c.command,
	}
	if len(c.arguments) > 0 {
		args := make([]map[string]any, 0)

		for _, argument := range c.arguments {
			arg := map[string]any{"_": argument}
			args = append(args, arg)
		}
		body[fmt.Sprintf("%s:Arguments", wsmv.NS_WIN_SHELL)] = args
	}

	return body
}

func (c *Command) commandHeaderOpts() map[string]any {
	options := []map[string]any{
		{
			"_": c.consoleMode,
			":attributes!": map[string]any{
				"Name": "WINRS_CONSOLEMODE_STDIN",
			},
		},
		{
			"_": c.skipCmdShell,
			":attributes!": map[string]any{
				"Name": "WINRS_SKIP_CMD_SHELL",
			},
		},
	}

	return map[string]any{
		fmt.Sprintf("%s:OptionSet", wsmv.NS_WSMAN_DMTF): map[string]any{
			fmt.Sprintf("%s:Option", wsmv.NS_WSMAN_DMTF): options,
		},
	}
}
