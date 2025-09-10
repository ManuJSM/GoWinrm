package fm

import (
	"testing"
	"time"

	"github.com/ManuJSM/GoWinrm/internal/transport"
	"github.com/ManuJSM/GoWinrm/internal/wsmv"
)

var (
	dummyDomain      = ""
	dummyWorkstation = "TESTWORKSTATION"
	dummyUsername    = ""
	dummyPassword    = ""
)

func TestFM(t *testing.T) {
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

	fm := NewFileManager(ntlmNego, &sessionOpts)
	src := ""
	dst := ""

	err := fm.UploadFile(src, dst)
	if err != nil {
		t.Error(err)
	}
}
