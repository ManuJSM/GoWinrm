package shell

import (
	"testing"
	"time"

	"github.com/ManuJSM/GoWinrm/internal/transport"
	"github.com/ManuJSM/GoWinrm/internal/wsmv"
)

func TestPs(t *testing.T) {
	endpoint := "http://localhost:5985/wsman"

	ntlmNego := transport.NewNtlmNego(endpoint, &transport.NegotiateOpts{
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

	ps := NewPsShell(ntlmNego, &sessionOpts)
	cmd := `Get-Process | Sort-Object -Property WS -Descending | Select-Object -First 3 -Property Name, Id, @{Name='Memoria (%)'; Expression = { [math]::Round(($_.WS / (Get-CimInstance Win32_ComputerSystem).TotalPhysicalMemory) * 100, 2) }} | Format-Table -AutoSize`

	oc, err := ps.RunCommand(cmd)
	if err != nil {
		t.Error(err)
	} else {
		t.Log("STDOUT: ", oc.Stdout.String())
		t.Log("STDERR: ", oc.Stderr.String())
		t.Log("EXITCODE: ", oc.ExitCode)
	}

	err = ps.Close()
	if err != nil {
		t.Error(err)
	}

}
