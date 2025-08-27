package transport

import (
	"testing"

	"github.com/ManuJSM/GoNtlm/client"
)

var (
	dummyDomain      = ""
	dummyWorkstation = "TESTWORKSTATION"
	dummyUsername    = ""
	dummyPassword    = ""
)

func generateSOAPMessage() []byte {

	soap := `<?xml version="1.0" encoding="UTF-8"?>
<env:Envelope xmlns:w="http://schemas.dmtf.org/wbem/wsman/1/wsman.xsd" xmlns:cfg="http://schemas.microsoft.com/wbem/wsman/1/config" xmlns:x="http://schemas.xmlsoap.org/ws/2004/09/transfer" xmlns:p="http://schemas.microsoft.com/wbem/wsman/1/wsman.xsd" xmlns:rsp="http://schemas.microsoft.com/wbem/wsman/1/windows/shell" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:env="http://www.w3.org/2003/05/soap-envelope" xmlns:a="http://schemas.xmlsoap.org/ws/2004/08/addressing" xmlns:b="http://schemas.dmtf.org/wbem/wsman/1/cimbinding.xsd" xmlns:n="http://schemas.xmlsoap.org/ws/2004/09/enumeration">
  <env:Header>
    <a:To>http://192.168.1.6:5985/wsman</a:To>
    <a:ReplyTo>
      <a:Address mustUnderstand="true">http://schemas.xmlsoap.org/ws/2004/08/addressing/role/anonymous</a:Address>
    </a:ReplyTo>
    <w:ResourceURI mustUnderstand="true">http://schemas.microsoft.com/wbem/wsman/1/windows/shell/cmd</w:ResourceURI>
    <w:SelectorSet>
      <w:Selector Name="ShellId">123e4567-e89b-12d3-a456-426614174000</w:Selector>
    </w:SelectorSet>
    <w:MaxEnvelopeSize mustUnderstand="true">153600</w:MaxEnvelopeSize>
    <p:SessionId mustUnderstand="false">uuid:123e4567-e89b-12d3-a456-426614174000</p:SessionId>
    <a:MessageID>uuid:8294d95c-b0ce-4bfa-9d08-0e56cc242d37</a:MessageID>
    <w:Locale mustUnderstand="false" xml:lang="en-US"></w:Locale>
    <p:DataLocale xml:lang="en-US" mustUnderstand="false"></p:DataLocale>
    <a:Action mustUnderstand="true">http://schemas.xmlsoap.org/ws/2004/09/transfer/Delete</a:Action>
    <w:OperationTimeout>PT60S</w:OperationTimeout>
  </env:Header>
  <env:Body></env:Body>
</env:Envelope>`

	return []byte(soap)
}

func TestAuth(t *testing.T) {
	endpoint := "http://localhost:5985/wsman"

	ntlmNego := NewNtlmNego(endpoint, &client.ClientOpts{
		Domain:      dummyDomain,
		Workstation: dummyWorkstation,
		Username:    dummyUsername,
		Password:    dummyPassword,
	})

	msg, err := ntlmNego.SendRequest(generateSOAPMessage())
	if err != nil {
		t.Error(err)
	} else {
		t.Log(string(msg))
	}

}
