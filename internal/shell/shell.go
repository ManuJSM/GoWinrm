package shell

import (
	"GoWinrm/internal/transport"
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
