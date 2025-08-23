package wsmv

import (
	"GoWinrm/internal/utils"
	"fmt"
	"maps"
	"time"
)

// Constantes equivalentes a los URI de WSMan y SOAP
const (
	RESOURCEURICMD        = "http://schemas.microsoft.com/wbem/wsman/1/windows/shell/cmd"
	RESOURCEURIPOWERSHELL = "http://schemas.microsoft.com/powershell/Microsoft.PowerShell"
)

type SessionOptions struct {
	Endpoint         string
	MaxEnvelopeSize  int
	SessionID        string
	Locale           string
	OperationTimeout time.Duration
}

// mergeHeaders simula la fusión de mapas, con merge especial para key ":attributes!"
func MergeHeaders(headers ...map[string]any) map[string]any {
	result := make(map[string]any)
	for _, h := range headers {
		for k, v := range h {
			if k == ":attributes!" {
				// Se espera que el valor sea un map[string] any
				if existingAttr, ok := result[k].(map[string]any); ok {
					if newAttr, ok2 := v.(map[string]any); ok2 {
						maps.Copy(existingAttr, newAttr)
						result[k] = existingAttr
						continue
					}
				}
			}
			result[k] = v
		}
	}
	return result
}

func SharedHeaders(sessionOpts SessionOptions) map[string]any {
	return map[string]any{
		fmt.Sprintf("%s:To", NS_ADDRESSING): sessionOpts.Endpoint,

		fmt.Sprintf("%s:ReplyTo", NS_ADDRESSING): map[string]any{
			fmt.Sprintf("%s:Address", NS_ADDRESSING): map[string]any{
				"_": "http://schemas.xmlsoap.org/ws/2004/08/addressing/role/anonymous",
				":attributes!": map[string]any{
					"mustUnderstand": true,
				},
			},
		},

		fmt.Sprintf("%s:MaxEnvelopeSize", NS_WSMAN_DMTF): sessionOpts.MaxEnvelopeSize,
		fmt.Sprintf("%s:MessageID", NS_ADDRESSING):       "uuid:" + utils.NewUuid(),
		fmt.Sprintf("%s:SessionId", NS_WSMAN_MSFT):       "uuid:" + sessionOpts.SessionID,

		fmt.Sprintf("%s:Locale", NS_WSMAN_DMTF): map[string]any{
			":attributes!": map[string]any{
				"xml:lang":       sessionOpts.Locale,
				"mustUnderstand": false,
			},
		},

		fmt.Sprintf("%s:DataLocale", NS_WSMAN_MSFT): map[string]any{
			":attributes!": map[string]any{
				"xml:lang":       sessionOpts.Locale,
				"mustUnderstand": false,
			},
		},

		fmt.Sprintf("%s:OperationTimeout", NS_WSMAN_DMTF): utils.Iso8601Duration(int(sessionOpts.OperationTimeout.Seconds())),

		":attributes!": map[string]any{
			fmt.Sprintf("%s:MaxEnvelopeSize", NS_WSMAN_DMTF): map[string]any{"mustUnderstand": true},
			fmt.Sprintf("%s:SessionId", NS_WSMAN_MSFT): map[string]any{
				"mustUnderstand": false,
			},
		},
	}
}

func ResourceURIShell(shellURI string) map[string]any {
	return map[string]any{
		fmt.Sprintf("%s:ResourceURI", NS_WSMAN_DMTF): shellURI,
		":attributes!": map[string]any{
			fmt.Sprintf("%s:ResourceURI", NS_WSMAN_DMTF): map[string]any{
				"mustUnderstand": true,
			},
		},
	}
}

func ResourceURIWMI(namespace string) map[string]any {
	if namespace == "" {
		namespace = "root/cimv2/*"
	}
	return map[string]any{
		fmt.Sprintf("%s:ResourceURI", NS_WSMAN_DMTF): "http://schemas.microsoft.com/wbem/wsman/1/wmi/" + namespace,
		":attributes!": map[string]any{
			fmt.Sprintf("%s:ResourceURI", NS_WSMAN_DMTF): map[string]any{
				"mustUnderstand": true,
			},
		},
	}
}

func ActionGet() map[string]any {
	return actionWithURL("http://schemas.xmlsoap.org/ws/2004/09/transfer/Get")
}

func ActionCreate() map[string]any {
	return actionWithURL("http://schemas.xmlsoap.org/ws/2004/09/transfer/Create")
}

func ActionDelete() map[string]any {
	return actionWithURL("http://schemas.xmlsoap.org/ws/2004/09/transfer/Delete")
}

func ActionCommand() map[string]any {
	return actionWithURL("http://schemas.microsoft.com/wbem/wsman/1/windows/shell/Command")
}

func ActionReceive() map[string]any {
	return actionWithURL("http://schemas.microsoft.com/wbem/wsman/1/windows/shell/Receive")
}

func ActionSend() map[string]any {
	return actionWithURL("http://schemas.microsoft.com/wbem/wsman/1/windows/shell/Send")
}

func ActionSignal() map[string]any {
	return actionWithURL("http://schemas.microsoft.com/wbem/wsman/1/windows/shell/Signal")
}

func ActionEnumerate() map[string]any {
	return actionWithURL("http://schemas.xmlsoap.org/ws/2004/09/enumeration/Enumerate")
}

func ActionEnumeratePull() map[string]any {
	return actionWithURL("http://schemas.xmlsoap.org/ws/2004/09/enumeration/Pull")
}

func actionWithURL(url string) map[string]any {
	return map[string]any{
		fmt.Sprintf("%s:Action", NS_ADDRESSING): url,
		":attributes!": map[string]any{
			fmt.Sprintf("%s:Action", NS_ADDRESSING): map[string]any{
				"mustUnderstand": true,
			},
		},
	}
}

func SelectorShellID(shellID string) map[string]any {
	return map[string]any{
		fmt.Sprintf("%s:SelectorSet", NS_WSMAN_DMTF): map[string]any{
			fmt.Sprintf("%s:Selector", NS_WSMAN_DMTF): map[string]any{
				"_": shellID,
				":attributes!": map[string]any{
					"Name": "ShellId",
				},
			},
		},
	}
}
