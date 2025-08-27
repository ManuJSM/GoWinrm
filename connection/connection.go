package connection

import (
	"GoWinrm/internal/shell"
	"GoWinrm/internal/transport"
	"GoWinrm/internal/utils"
	"GoWinrm/internal/wsmv"
	"fmt"
	"time"
)

type ConnConf struct {
	*wsmv.SessionOptions
	*transport.NegotiateOpts
}

const (
	DEFAULTLOCALE             = "en-US"
	DEFAULT_OPERATION_TIMEOUT = 60 * time.Second
	DEFAULT_MAX_ENV_SIZE      = 153600
)

func NewConf(arg utils.Arguments) (*ConnConf, error) {
	cc := &ConnConf{
		SessionOptions: &wsmv.SessionOptions{},
		NegotiateOpts:  &transport.NegotiateOpts{},
	}
	var ok bool

	cc.Endpoint, ok = arg["endpoint"].(string)
	if !ok {
		return nil, fmt.Errorf("falta endpoint")
	}
	cc.Username, ok = arg["username"].(string)
	if !ok {
		return nil, fmt.Errorf("falta username")
	}
	cc.Password, ok = arg["password"].(string)
	if !ok {
		return nil, fmt.Errorf("falta password")
	}

	cc.Workstation = utils.OptOrDefault(arg, "workstation", "").(string)
	cc.Domain = utils.OptOrDefault(arg, "domain", "").(string)

	cc.Locale = utils.OptOrDefault(arg, "locale", DEFAULTLOCALE).(string)
	cc.SessionID = utils.NewUuid()
	cc.Flags = utils.OptOrDefault(arg, "flags", uint32(0)).(uint32)
	cc.MaxEnvelopeSize = utils.OptOrDefault(arg, "maxEnvelopeSize", DEFAULT_MAX_ENV_SIZE).(int)
	cc.OperationTimeout = utils.OptOrDefault(arg, "opTimeout", DEFAULT_OPERATION_TIMEOUT).(time.Duration)

	return cc, nil
}

type Connection struct {
	Shell     shell.Shell
	transport transport.Transport
}

func NewConnection(connConf *ConnConf) *Connection {

	t := transport.NewNtlmNego(connConf.Endpoint, connConf.NegotiateOpts)

	return &Connection{
		transport: t,
		Shell:     shell.NewCmdShell(t, connConf.SessionOptions),
	}

}
