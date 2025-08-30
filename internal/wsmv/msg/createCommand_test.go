package msg

import (
	"fmt"
	"testing"
	"time"

	"github.com/ManuJSM/GoWinrm/internal/wsmv"
)

func TestCommandFullMessage(t *testing.T) {
	sessionOpts := wsmv.SessionOptions{
		Endpoint:         "http://192.168.1.6:5985/wsman",
		MaxEnvelopeSize:  153600,
		SessionID:        "123e4567-e89b-12d3-a456-426614174000",
		Locale:           "en-US",
		OperationTimeout: 60 * time.Second,
	}
	cmdOpts := map[string]any{
		"shell_id":           "123e4567-e89b-12d3-a456-426614174000",
		"command":            "whoami",
		"arguments":          []string{"/user"},
		"console_mode_stdin": "TRUE",
		"skip_cmd_shell":     "FALSE",
	}

	cmd := NewCommand(sessionOpts, cmdOpts)

	headers := cmd.Headers()
	body := cmd.Body()

	// Verificaciones básicas de headers
	if headers["w:ResourceURI"] != wsmv.RESOURCEURICMD {
		t.Errorf("ResourceURI incorrecto: %v", headers["w:ResourceURI"])
	}

	if _, ok := headers["w:SelectorSet"]; !ok {
		t.Error("SelectorSet no encontrado en headers")
	}

	if _, ok := headers["w:OptionSet"]; !ok {
		t.Error("OptionSet no encontrado en headers (esperado para shell cmd)")
	}

	// Verificaciones básicas de body
	cmdLine, ok := body["rsp:CommandLine"].(map[string]any)
	if !ok {
		t.Fatal("CommandLine no presente o tipo incorrecto")
	}

	if cmdLine["rsp:Command"] != "whoami" {
		t.Errorf("Command esperado 'whoami', recibido: %v", cmdLine["rsp:Command"])
	}

	args, ok := cmdLine["rsp:Arguments"].([]string)
	if !ok || len(args) != 1 || args[0] != "/user" {
		t.Errorf("Arguments incorrectos: %v", cmdLine["rsp:Arguments"])
	}
	xml, err := wsmv.BuildXML(cmd.Headers(), cmd.Body())
	if err != nil {
		t.Error(err)
	} else {
		fmt.Println(string(xml))
	}
}
