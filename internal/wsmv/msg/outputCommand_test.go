package msg

import (
	"testing"
	"time"

	"github.com/ManuJSM/GoWinrm/internal/wsmv"
)

func TestOutputCommandSimple(t *testing.T) {
	// Configurar opciones de sesión
	sessionOpts := wsmv.SessionOptions{
		Endpoint:         "http://192.168.1.6:5985/wsman",
		MaxEnvelopeSize:  153600,
		SessionID:        "123e4567-e89b-12d3-a456-426614174000",
		Locale:           "en-US",
		OperationTimeout: 60 * time.Second,
	}

	// Configurar opciones del shell
	shellOpts := map[string]any{
		"shell_uri":  wsmv.RESOURCEURICMD,
		"shell_id":   "1234",
		"command_id": "1234",
	}

	// Crear instancia de CreateShell
	outputCommand := NewOutputCommand(sessionOpts, shellOpts)
	xml, err := wsmv.BuildXML(outputCommand.Headers(), outputCommand.Body())
	if err != nil {
		t.Error(err)
	} else {
		t.Log(string(xml))
	}

}
