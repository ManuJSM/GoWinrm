package GoWinrm

import (
	"fmt"
	"time"

	"github.com/ManuJSM/GoWinrm/internal/fm"
	"github.com/ManuJSM/GoWinrm/internal/shell"
	"github.com/ManuJSM/GoWinrm/internal/transport"
	"github.com/ManuJSM/GoWinrm/internal/utils"
	"github.com/ManuJSM/GoWinrm/internal/wsmv"
)

// ConnConf encapsula opciones para sesión y negociación NTLM.
type ConnConf struct {
	*wsmv.SessionOptions     // Opciones específicas de la sesión WinRM (por ejemplo, locale, flags)
	*transport.NegotiateOpts // Opciones relacionadas con la negociación NTLM (endpoint, usuario, etc.)
}

// Valores por defecto para configuración
const (
	DEFAULTLOCALE             = "en-US"
	DEFAULT_OPERATION_TIMEOUT = 60 * time.Second
	DEFAULT_MAX_ENV_SIZE      = 153600
)

// NewConf crea una nueva configuración de conexión (ConnConf) a partir de un map de argumentos.
func NewConf(arg utils.Arguments) (*ConnConf, error) {
	cc := &ConnConf{
		SessionOptions: &wsmv.SessionOptions{},
		NegotiateOpts:  &transport.NegotiateOpts{},
	}

	var ok bool

	// Se requiere endpoint
	cc.Endpoint, ok = arg["endpoint"].(string)
	if !ok {
		return nil, fmt.Errorf("missing endpoint")
	}

	// Se requiere username
	cc.Username, ok = arg["username"].(string)
	if !ok {
		return nil, fmt.Errorf("missing username")
	}

	// Se requiere password
	cc.Password, ok = arg["password"].(string)
	if !ok {
		return nil, fmt.Errorf("missing password")
	}

	// Opcionales: workstation y domain
	cc.Workstation = utils.OptOrDefault(arg, "workstation", "").(string)
	cc.Domain = utils.OptOrDefault(arg, "domain", "").(string)

	// Opcionales: locale, flags, tamaño máximo de sobre y timeout
	cc.Locale = utils.OptOrDefault(arg, "locale", DEFAULTLOCALE).(string)
	cc.SessionID = utils.NewUuid() // Genera UUID único para la sesión
	cc.Flags = uint32(0)
	cc.MaxEnvelopeSize = utils.OptOrDefault(arg, "maxEnvelopeSize", DEFAULT_MAX_ENV_SIZE).(int)
	cc.OperationTimeout = utils.OptOrDefault(arg, "opTimeout", DEFAULT_OPERATION_TIMEOUT).(time.Duration)

	return cc, nil
}

// Connection representa una conexión activa que encapsula el transporte y el shell.
type Connection struct {
	Shell     shell.Shell         // Shell remoto (permite ejecutar comandos)
	transport transport.Transport // Transporte de red (ej. NTLM sobre HTTP)
	FileM     *fm.FileManager     // Descargar y subida de archivos
}

// NewConnection establece una nueva conexión usando configuración ConnConf.
func NewConnection(connConf *ConnConf) *Connection {
	// Crea transporte usando NTLM
	t := transport.NewNtlmNego(connConf.Endpoint, connConf.NegotiateOpts)

	// Retorna una conexión con shell sobre ese transporte
	return &Connection{
		transport: t,
		Shell:     shell.NewPsShell(t, connConf.SessionOptions),
		FileM:     fm.NewFileManager(t, connConf.SessionOptions),
	}
}
