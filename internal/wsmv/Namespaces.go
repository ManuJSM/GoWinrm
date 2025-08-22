package wsmv

const (
	NS_SOAP_ENV    = "env" // http://www.w3.org/2003/05/soap-envelope
	NS_ADDRESSING  = "a"   // http://schemas.xmlsoap.org/ws/2004/08/addressing
	NS_CIMBINDING  = "b"   // http://schemas.dmtf.org/wbem/wsman/1/cimbinding.xsd
	NS_ENUM        = "n"   // http://schemas.xmlsoap.org/ws/2004/09/enumeration
	NS_TRANSFER    = "x"   // http://schemas.xmlsoap.org/ws/2004/09/transfer
	NS_WSMAN_DMTF  = "w"   // http://schemas.dmtf.org/wbem/wsman/1/wsman.xsd
	NS_WSMAN_MSFT  = "p"   // http://schemas.microsoft.com/wbem/wsman/1/wsman.xsd
	NS_SCHEMA_INST = "xsi" // http://www.w3.org/2001/XMLSchema-instance
	NS_WIN_SHELL   = "rsp" // http://schemas.microsoft.com/wbem/wsman/1/windows/shell
	NS_WSMAN_FAULT = "f"   // http://schemas.microsoft.com/wbem/wsman/1/wsmanfault
	NS_WSMAN_CONF  = "cfg" // http://schemas.microsoft.com/wbem/wsman/1/config
)

var Namespaces = map[string]string{
	"xmlns:xsd":              "http://www.w3.org/2001/XMLSchema",
	"xmlns:xsi":              "http://www.w3.org/2001/XMLSchema-instance",
	"xmlns:" + NS_SOAP_ENV:   "http://www.w3.org/2003/05/soap-envelope",
	"xmlns:" + NS_ADDRESSING: "http://schemas.xmlsoap.org/ws/2004/08/addressing",
	"xmlns:" + NS_CIMBINDING: "http://schemas.dmtf.org/wbem/wsman/1/cimbinding.xsd",
	"xmlns:" + NS_ENUM:       "http://schemas.xmlsoap.org/ws/2004/09/enumeration",
	"xmlns:" + NS_TRANSFER:   "http://schemas.xmlsoap.org/ws/2004/09/transfer",
	"xmlns:" + NS_WSMAN_DMTF: "http://schemas.dmtf.org/wbem/wsman/1/wsman.xsd",
	"xmlns:" + NS_WSMAN_MSFT: "http://schemas.microsoft.com/wbem/wsman/1/wsman.xsd",
	"xmlns:" + NS_WIN_SHELL:  "http://schemas.microsoft.com/wbem/wsman/1/windows/shell",
	"xmlns:" + NS_WSMAN_CONF: "http://schemas.microsoft.com/wbem/wsman/1/config",
}
