package shell

import (
	"GoWinrm/internal/log"
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

	return &Cmd{
		Shell: s,
	}

}

func (c *Cmd) Open() error {

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

func (c *Cmd) SendCommand(command string, arguments ...string) (string, error) {

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

func (c *Cmd) cleanCommand(commandId string) error {

	log.Debug("Cleaning up command " + commandId)

	cmdOpts := map[string]any{
		"shell_id":   c.shellID,
		"command_id": commandId,
	}

	cleanMsg := msg.NewCleanCommand(*c.sessionOpts, cmdOpts)
	xml, err := wsmv.BuildXML(cleanMsg.Headers(), cleanMsg.Body())
	if err != nil {
		return err
	}

	_, err = c.transport.SendRequest(xml)
	if err != nil {
		return err
	}

	return nil

}
