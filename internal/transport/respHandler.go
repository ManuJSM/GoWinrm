package transport

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/antchfx/xmlquery"
)

func RespHandler(xml []byte, statusCode int) ([]byte, error) {

	err := raiseIfError(string(xml), statusCode)
	if err != nil {
		return nil, err
	}

	return xml, nil
}
func raiseIfError(xml string, statusCode int) error {

	if statusCode == http.StatusOK {
		return nil
	}

	if statusCode == http.StatusUnauthorized {
		return errors.New("authorization error: received HTTP 401")
	}
	doc, err := parseXML(xml)

	if err != nil {
		return err
	}

	if err := raiseIfWSManFault(doc); err != nil {
		return err
	}

	if err := raiseIfWMIError(doc); err != nil {
		return err
	}

	if err := raiseIfSOAPFault(doc); err != nil {
		return err
	}

	return fmt.Errorf("http transport error: status=%d, body=%s", statusCode, xml)
}

func parseXML(xml string) (*xmlquery.Node, error) {

	doc, err := xmlquery.Parse(strings.NewReader(xml))
	if err != nil {
		return nil, fmt.Errorf("unable to parse WinRM response: %v", err)
	}

	return doc, nil
}

func raiseIfWMIError(xml *xmlquery.Node) error {
	node := xmlquery.FindOne(xml, "//*[local-name()='Envelope']/*[local-name()='Body']/*[local-name()='Fault']//*[local-name()='MSFT_WmiError']")
	if node != nil {
		codeNode := xmlquery.FindOne(node, ".//*[local-name()='error_Code']")
		code := ""
		if codeNode != nil {
			code = codeNode.InnerText()
		}
		return fmt.Errorf("WMI error found: code=%s, details=%s", code, node.OutputXML(true))
	}
	return nil
}

func raiseIfWSManFault(xml *xmlquery.Node) error {
	node := xmlquery.FindOne(xml, "//*[local-name()='Envelope']/*[local-name()='Body']/*[local-name()='Fault']//*[local-name()='WSManFault']")
	if node != nil {
		code := node.SelectAttr("Code")
		return fmt.Errorf("WSMan fault found: code=%s, details=%s", code, node.OutputXML(true))
	}
	return nil
}

func raiseIfSOAPFault(xml *xmlquery.Node) error {
	fault := xmlquery.FindOne(xml, "//*[local-name()='Envelope']/*[local-name()='Body']/*[local-name()='Fault']")
	if fault == nil {
		return nil
	}

	code := getNodeText(xmlquery.FindOne(fault, ".//*[local-name()='Code']/*[local-name()='Value']"))
	subcode := getNodeText(xmlquery.FindOne(fault, ".//*[local-name()='Subcode']/*[local-name()='Value']"))
	reason := getNodeText(xmlquery.FindOne(fault, ".//*[local-name()='Reason']/*[local-name()='Text']"))

	if code != "" || subcode != "" || reason != "" {
		return fmt.Errorf("SOAP fault found: code=%s, subcode=%s, reason=%s", code, subcode, reason)
	}

	return nil
}
func getNodeText(n *xmlquery.Node) string {
	if n == nil {
		return ""
	}
	return n.InnerText()
}
