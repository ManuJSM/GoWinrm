package connection

import (
	"GoWinrm/internal/shell"
	"GoWinrm/internal/transport"
	"GoWinrm/internal/wsmv"
)

type ConnConf struct {
	*wsmv.SessionOptions
	*transport.NegotiateOpts
}

//TODO seguir con la configuracion

const DEFAULTLOCALE = "en-US"

func NewDefaultConf(endpoint string, domain string) *ConnConf {
	cc := &ConnConf{
		SessionOptions: &wsmv.SessionOptions{},
		NegotiateOpts:  &transport.NegotiateOpts{},
	}
	cc.Endpoint = endpoint
	cc.Domain = domain
	cc.Locale = DEFAULTLOCALE

	return cc
}

type Connection struct {
	shell     shell.Shell
	transport transport.Transport
}

func NewConnection(endpoint string, connConf *ConnConf) *Connection {
	t := transport.NewNtlmNego(endpoint, connConf.NegotiateOpts)

	return &Connection{
		transport: t,
		shell:     shell.NewCmdShell(t, connConf.SessionOptions),
	}

}
