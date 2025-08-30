package connection_test

import (
	"testing"

	"github.com/ManuJSM/GoWinrm"
	"github.com/ManuJSM/GoWinrm/internal/utils"
)

var args = utils.Arguments{
	"endpoint":    "http://localhost:5985/wsman",
	"workstation": "TESTWORKSTATION",
	"username":    "",
	"password":    "",
}

func TestCmd(t *testing.T) {

	conf, err := GoWinrm.NewConf(args)

	if err != nil {
		t.Error(err)
		return
	}

	conn := GoWinrm.NewConnection(conf)

	oc, err := conn.Shell.RunCommand("whoami")
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
