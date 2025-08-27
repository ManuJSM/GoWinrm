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
	endpoint := "http://localhost:5985/wsman"

	ntlmNego := transport.NewNtlmNego(endpoint, &client.ClientOpts{
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

	oc, err := cmd.RunCommand("whoami")
	if err != nil {
		t.Error(err)
	} else {
		t.Log("STDOUT: ", oc.Stdout.String())
		t.Log("STDERR: ", oc.Stderr.String())
		t.Log("EXITCODE: ", oc.ExitCode)
	}

	err = cmd.Close()
	if err != nil {
		t.Error(err)
	}

}
