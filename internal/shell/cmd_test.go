package shell

import (
	"GoWinrm/internal/transport"
	"GoWinrm/internal/wsmv"
	"testing"
	"time"

	"github.com/ManuJSM/GoNtlm/client"
)

var (
	dummyDomain      = ""
	dummyWorkstation = "TESTWORKSTATION"
	dummyUsername    = ""
	dummyPassword    = ""
)

func TestCmd(t *testing.T) {
	endpoint := "http://192.168.1.6:5985/wsman"

	ntlmNego := transport.NewNtlmNego(endpoint, client.ClientOpts{
		Domain:      dummyDomain,
		Workstation: dummyWorkstation,
		Username:    dummyUsername,
		Password:    dummyPassword,
	})
	sessionOpts := wsmv.SessionOptions{
		Endpoint:         endpoint,
		MaxEnvelopeSize:  153600,
		SessionID:        "123e4567-e89b-12d3-a456-426614174000",
		Locale:           "en-US",
		OperationTimeout: 60 * time.Second,
	}

	cmd := NewCmdShell(ntlmNego, &sessionOpts)

	err := cmd.Open()
	if err != nil {
		t.Error(err)
	}

	idCommand, err := cmd.SendCommand("whoami")
	if err != nil {
		t.Error(err)
	}

	err = cmd.cleanCommand(idCommand)
	if err != nil {
		t.Error(err)
	}
	err = cmd.Close()
	if err != nil {
		t.Error(err)
	}

}
