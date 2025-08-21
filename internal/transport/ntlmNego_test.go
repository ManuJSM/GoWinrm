package transport

import (
	"crypto/rand"
	"fmt"
	"testing"

	"github.com/ManuJSM/GoNtlm/client"
)

//TODO usar .ENV

var (
	dummyDomain      = ""
	dummyWorkstation = "TESTWORKSTATION"
	dummyUsername    = ""
	dummyPassword    = ""
)

func generateSOAPMessage() []byte {
	uuid := make([]byte, 16)
	_, err := rand.Read(uuid)
	if err != nil {
		panic("no se pudo generar UUID")
	}
	// Formatear según RFC 4122
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // versión 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // variante 10

	uuidStr := fmt.Sprintf("uuid:%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:16])

	soap := fmt.Sprintf(`<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope"
            xmlns:a="http://schemas.xmlsoap.org/ws/2004/08/addressing"
            xmlns:w="http://schemas.dmtf.org/wbem/wsman/1/wsman.xsd"
            xmlns:rsp="http://schemas.microsoft.com/wbem/wsman/1/windows/shell">
  <s:Header>
    <a:To>http://192.168.1.6:5985/wsman</a:To>
    <w:ResourceURI>http://schemas.microsoft.com/wbem/wsman/1/windows/shell/cmd</w:ResourceURI>
    <a:ReplyTo>
      <a:Address>http://schemas.xmlsoap.org/ws/2004/08/addressing/role/anonymous</a:Address>
    </a:ReplyTo>
    <a:MessageID>%s</a:MessageID>
    <a:Action>http://schemas.xmlsoap.org/ws/2004/09/transfer/Create</a:Action>
  </s:Header>
  <s:Body>
    <rsp:Shell>
      <rsp:InputStreams>stdin</rsp:InputStreams>
      <rsp:OutputStreams>stdout stderr</rsp:OutputStreams>
    </rsp:Shell>
  </s:Body>
</s:Envelope>
`, uuidStr)

	return []byte(soap)
}

func TestAuth(t *testing.T) {
	endpoint := "http://192.168.1.6:5985/wsman"

	ntlmNego := NewNtlmNego(endpoint, client.ClientOpts{
		Domain:      dummyDomain,
		Workstation: dummyWorkstation,
		Username:    dummyUsername,
		Password:    dummyPassword,
	})
	// err := ntlmNego.InitAuth()
	// if err != nil {
	// 	t.Error(err)
	// }
	msg, err := ntlmNego.SendRequest(generateSOAPMessage())
	if err != nil {
		t.Error(err)
	} else {
		t.Log(string(msg))
	}

}
