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
<env:Envelope xmlns:p="http://schemas.microsoft.com/wbem/wsman/1/wsman.xsd" xmlns:rsp="http://schemas.microsoft.com/wbem/wsman/1/windows/shell" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:a="http://schemas.xmlsoap.org/ws/2004/08/addressing" xmlns:n="http://schemas.xmlsoap.org/ws/2004/09/enumeration" xmlns:w="http://schemas.dmtf.org/wbem/wsman/1/wsman.xsd" xmlns:cfg="http://schemas.microsoft.com/wbem/wsman/1/config" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:env="http://www.w3.org/2003/05/soap-envelope" xmlns:b="http://schemas.dmtf.org/wbem/wsman/1/cimbinding.xsd" xmlns:x="http://schemas.xmlsoap.org/ws/2004/09/transfer">
  <env:Header>
    <w:MaxEnvelopeSize mustUnderstand="true">153600</w:MaxEnvelopeSize>
    <p:DataLocale xml:lang="en-US" mustUnderstand="false"></p:DataLocale>
    <w:OperationTimeout>PT60S</w:OperationTimeout>
    <a:To>http://192.168.1.6:5985/wsman</a:To>
    <a:ReplyTo>
      <a:Address mustUnderstand="true">http://schemas.xmlsoap.org/ws/2004/08/addressing/role/anonymous</a:Address>
    </a:ReplyTo>
    <w:Locale mustUnderstand="false" xml:lang="en-US"></w:Locale>
    <w:ResourceURI mustUnderstand="true">http://schemas.microsoft.com/wbem/wsman/1/windows/shell/cmd</w:ResourceURI>
    <a:Action mustUnderstand="true">http://schemas.xmlsoap.org/ws/2004/09/transfer/Create</a:Action>
    <p:SessionId mustUnderstand="false">uuid:123e4567-e89b-12d3-a456-426614174000</p:SessionId>
    <a:MessageID>uuid:d2df3f48-1705-49fa-af18-ac88dee47bf7</a:MessageID>
    <w:OptionSet>
      <w:Option Name="WINRS_NOPROFILE">FALSE</w:Option>
      <w:Option Name="WINRS_CODEPAGE">65001</w:Option>
    </w:OptionSet>
  </env:Header>
  <env:Body>
    <rsp:Shell>
      <rsp:OutputStreams>stdout stderr</rsp:OutputStreams>
      <rsp:WorkingDirectory>C:\Users\Administrator</rsp:WorkingDirectory>
      <rsp:IdleTimeOut>PT5M0S</rsp:IdleTimeOut>
      <rsp:Environment>
        <rsp:Variable Name="PATH">C:\Windows\System32;C:\Windows</rsp:Variable>
        <rsp:Variable Name="TEMP">C:\Temp</rsp:Variable>
        <rsp:Variable Name="LANGUAGE">en_US.UTF-8</rsp:Variable>
      </rsp:Environment>
      <rsp:InputStreams>stdin</rsp:InputStreams>
    </rsp:Shell>
  </env:Body>
</env:Envelope>`

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
