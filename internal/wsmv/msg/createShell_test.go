package msg

import (
	. "GoWinrm/internal/wsmv"
	"fmt"
	"testing"
	"time"
)

func TestCreateShellSimple(t *testing.T) {
	// Configurar opciones de sesión
	sessionOpts := SessionOptions{
		Endpoint:         "http://192.168.1.6:5985/wsman",
		MaxEnvelopeSize:  153600,
		SessionID:        "123e4567-e89b-12d3-a456-426614174000",
		Locale:           "en-US",
		OperationTimeout: 60 * time.Second,
	}

	// Configurar opciones del shell
	shellOpts := map[string]any{
		"shell_uri":         RESOURCEURICMD,
		"i_stream":          "stdin",
		"o_stream":          "stdout stderr",
		"codepage":          65001,
		"noprofile":         "FALSE",
		"working_directory": "C:\\Users\\Administrator",
		"idle_timeout":      300 * time.Second, // 5 minutos
		"env_vars": map[string]string{
			"PATH":     "C:\\Windows\\System32;C:\\Windows",
			"TEMP":     "C:\\Temp",
			"LANGUAGE": "en_US.UTF-8",
		},
	}

	// Crear instancia de CreateShell
	createShell := NewCreateShell(sessionOpts, shellOpts)

	// Test básico - que no sea nil
	if createShell == nil {
		t.Fatal("CreateShell no debería ser nil")
	}

	// Test de valores
	if createShell.shellURI != RESOURCEURICMD {
		t.Errorf("Shell URI esperado: %s, obtenido: %s", RESOURCEURICMD, createShell.shellURI)
	}

	if createShell.iStream != "stdin" {
		t.Errorf("IStream esperado: stdin, obtenido: %s", createShell.iStream)
	}

	if createShell.workingDirectory != "C:\\Users\\Administrator" {
		t.Errorf("Working directory esperado: C:\\Users\\Administrator, obtenido: %s", createShell.workingDirectory)
	}

	// Test de headers y body
	headers := createShell.Headers()
	body := createShell.Body()

	if headers == nil {
		t.Error("Headers no debería ser nil")
	}

	if body == nil {
		t.Error("Body no debería ser nil")
	}

	// Verificar que tenemos los campos básicos en el body
	shellBody := body["rsp:Shell"].(map[string]any)
	if shellBody["rsp:InputStreams"] != "stdin" {
		t.Error("InputStreams no está en el body")
	}

	if shellBody["rsp:OutputStreams"] != "stdout stderr" {
		t.Error("OutputStreams no está en el body")
	}

	xml, err := BuildXML(headers, body)
	if err != nil {
		t.Error(err)
	} else {
		fmt.Println(xml)
	}

}
