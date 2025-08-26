package shell

import (
	"GoWinrm/internal/log"
	"GoWinrm/internal/transport"
	"GoWinrm/internal/utils"
	"GoWinrm/internal/wsmv"
	"GoWinrm/internal/wsmv/msg"
)

const (
	TOO_MANY_COMMANDS       = 2150859174
	ERROR_OPERATION_ABORTED = 995
	SHELL_NOT_FOUND         = 2150858843
)

// Lista de errores que requieren reinicio
var FAULTS_FOR_RESET = []uint32{
	SHELL_NOT_FOUND,   // Shell has been closed
	2147943418,        // Error reading registry key
	TOO_MANY_COMMANDS, // Maximum commands per user exceeded
}

type Shell struct {
	shellID     string
	shellURI    string
	transport   *transport.NtlmNego
	sessionOpts *wsmv.SessionOptions
	sendCommand func(string, ...string) (string, error)
	open        func() error
}

func (s *Shell) Close() error {
	if s.shellID == "" {
		return nil
	}

	cmdOpts := map[string]any{
		"shell_id": s.shellID,
	}

	msg := msg.NewCloseShell(*s.sessionOpts, cmdOpts)

	xml, err := wsmv.BuildXML(msg.Headers(), msg.Body())
	if err != nil {
		return err
	}

	_, err = s.transport.SendRequest(xml)
	if err != nil {
		return err
	}
	s.shellID = ""

	return nil
}

func (s *Shell) readOutput(commandId string) ([]byte, error) {
	cmdOpts := map[string]any{
		"shell_id":   s.shellID,
		"command_id": commandId,
	}

	outputMsg := msg.NewOutputCommand(*s.sessionOpts, cmdOpts)

	xml, err := wsmv.BuildXML(outputMsg.Headers(), outputMsg.Body())
	if err != nil {
		return nil, err
	}
	resp, err := s.transport.SendRequest(xml)
	if err != nil {
		return nil, err
	}
	return resp, nil

}

func (s *Shell) RunCommand(command string, arguments ...string) (*utils.OutputCommand, error) {
	//TODO implementar una logica de reintentos (2)

	if s.shellID == "" {
		s.open()
	}
	commandId, err := s.sendCommand(command, arguments...)
	if err != nil {
		return nil, err
	}
	defer s.cleanCommand(commandId)

	log.Debug("creating command_id:" + commandId + " on shell_id " + s.shellID)

	return s.handleOutput(commandId)
}

func (s *Shell) handleOutput(commandId string) (*utils.OutputCommand, error) {

	output := &utils.OutputCommand{}
	finished := false

	for !finished {

		resp, err := s.readOutput(commandId)
		if err != nil {
			return nil, err
		}
		finished = utils.ParseOutput(resp, output)
	}

	return output, nil
}

func (c *Shell) cleanCommand(commandId string) error {

	log.Debug("cleaning up command " + commandId)

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
