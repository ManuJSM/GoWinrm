package connection

import (
	"GoWinrm/internal/utils"
	"testing"
)

var args = utils.Arguments{
	"endpoint":    "http://localhost:5985/wsman",
	"workstation": "TESTWORKSTATION",
	"username":    "",
	"password":    "",
}

func TestCmd(t *testing.T) {

	conf, err := NewConf(args)

	if err != nil {
		t.Error(err)
		return
	}

	conn := NewConnection(conf)

	oc, err := conn.Shell.RunCommand("C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe", "-Command", "whoami /priv")
	if err != nil {
		t.Error(err)
	} else {
		t.Log("STDOUT: ", oc.Stdout.String())
		t.Log("STDERR: ", oc.Stderr.String())
		t.Log("EXITCODE: ", oc.ExitCode)
	}

	err = conn.Shell.Close()
	if err != nil {
		t.Error(err)
	}

}
