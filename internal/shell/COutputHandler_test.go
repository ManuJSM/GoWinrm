package shell_test

import (
	"GoWinrm/internal/shell"
	"fmt"
	"testing"
)

// Simula una respuesta XML WinRM válida con stdout/stderr
var sampleXML = []byte(`
<s:Envelope xml:lang="en-US" xmlns:s="http://www.w3.org/2003/05/soap-envelope" xmlns:a="http://schemas.xmlsoap.org/ws/2004/08/addressing" xmlns:w="http://schemas.dmtf.org/wbem/wsman/1/wsman.xsd" xmlns:rsp="http://schemas.microsoft.com/wbem/wsman/1/windows/shell" xmlns:p="http://schemas.microsoft.com/wbem/wsman/1/wsman.xsd"><s:Header><a:Action>http://schemas.microsoft.com/wbem/wsman/1/windows/shell/ReceiveResponse</a:Action><a:MessageID>uuid:28C482F3-AE0B-41B4-9877-BC280ACF3522</a:MessageID><a:To>http://schemas.xmlsoap.org/ws/2004/08/addressing/role/anonymous</a:To><a:RelatesTo>uuid:9dcb9250-162b-4992-a48e-3d0cd1f40717</a:RelatesTo></s:Header><s:Body><rsp:ReceiveResponse><rsp:Stream Name="stderr" CommandId="4D00F60C-6C79-4ED2-B90C-E02DD4C30323">RVJST1I6IA==</rsp:Stream><rsp:Stream Name="stderr" CommandId="4D00F60C-6C79-4ED2-B90C-E02DD4C30323">SW52YWxpZCBhcmd1bWVudC9vcHRpb24gLSAnWy91c2VyXScuDQpUeXBlICJXSE9BTUkgLz8iIGZvciB1c2FnZS4NCg==</rsp:Stream><rsp:Stream Name="stdout" CommandId="4D00F60C-6C79-4ED2-B90C-E02D
D4C30323" End="true"></rsp:Stream><rsp:Stream Name="stderr" CommandId="4D00F60C-6C79-4ED2-B90C-E02DD4C30323" End="true"></rsp:Stream><rsp:CommandState CommandId="4D00F60C-6C79-4ED2-B90C-E02DD4C30323" State="http://schemas.microsoft.com/wbem/wsman/1/windows/shell/CommandState/Done"><rsp:ExitCode>1</rsp:ExitCode></rsp:CommandState></rsp:ReceiveResponse></s:Body></s:Envelope>
`)

func TestCOHandler_HandleOutput(t *testing.T) {
	handler := shell.NewCOHandler()
	finished := handler.HandleOutput(sampleXML)

	if !finished {
		t.Error("not finished?")
	} else {
		fmt.Println("STDOUT: ", handler.Stdout.String())
		fmt.Println("STDERR: ", handler.Stderr.String())
		fmt.Println("EXITCODE: ", handler.ExitCode)
	}

}
