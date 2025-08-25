package shell

import (
	"GoWinrm/internal/transport"
	"GoWinrm/internal/utils"
	"GoWinrm/internal/wsmv"
	"GoWinrm/internal/wsmv/msg"
	"time"
)

type Cmd struct {
	*Shell
}

func NewCmdShell(transport *transport.NtlmNego, opts *wsmv.SessionOptions) *Cmd {

	s := &Shell{
		transport:   transport,
		shellURI:    wsmv.RESOURCEURICMD,
		sessionOpts: opts,
	}

	c := &Cmd{
		Shell: s,
	}
	s.sendCommand = c.sendCommand
	s.open = c.open
	return c
}

func (c *Cmd) open() error {

	shellOpts := map[string]any{
		"shell_uri":         c.shellURI,
		"i_stream":          "stdin",
		"o_stream":          "stdout stderr",
		"codepage":          65001,
		"noprofile":         "FALSE",
		"working_directory": "C:\\Users\\",
		"idle_timeout":      300 * time.Second, // 5 minutos
		"env_vars": map[string]string{
			"PATH":     "C:\\Windows\\System32;C:\\Windows",
			"TEMP":     "C:\\Temp",
			"LANGUAGE": "en_US.UTF-8",
		},
	}
	msg := msg.NewCreateShell(*c.sessionOpts, shellOpts)

	xml, err := wsmv.BuildXML(msg.Headers(), msg.Body())
	if err != nil {
		return err
	}
	resp, err := c.transport.SendRequest(xml)
	if err != nil {
		return err
	}

	c.shellID = utils.GetShellId(string(resp))

	return nil

}

func (c *Cmd) sendCommand(command string, arguments ...string) (string, error) {

	cmdOpts := map[string]any{
		"shell_id":           c.shellID,
		"command":            command,
		"arguments":          arguments,
		"console_mode_stdin": "TRUE",
		"skip_cmd_shell":     "FALSE",
	}

	m := msg.NewCommand(*c.sessionOpts, cmdOpts)
	request, err := wsmv.BuildXML(m.Headers(), m.Body())
	if err != nil {
		return "", err
	}
	resp, err := c.transport.SendRequest(request)
	if err != nil {
		return "", err
	}
	commandId := utils.GetCommandId(string(resp))

	return commandId, nil

}
