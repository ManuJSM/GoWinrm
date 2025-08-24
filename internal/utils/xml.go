package utils

import (
	"fmt"
	"strings"

	"github.com/antchfx/xmlquery"
)

func ParseXML(xml string) (*xmlquery.Node, error) {

	doc, err := xmlquery.Parse(strings.NewReader(xml))
	if err != nil {
		return nil, fmt.Errorf("unable to parse WinRM response: %v", err)
	}

	return doc, nil
}
func GetNodeText(n *xmlquery.Node) string {
	if n == nil {
		return ""
	}
	return n.InnerText()
}

func GetShellId(xml string) (id string) {
	doc, _ := ParseXML(xml)
	id = GetNodeText(xmlquery.FindOne(doc, "//*[@Name='ShellId']"))
	return
}

func GetCommandId(xml string) (id string) {
	doc, _ := ParseXML(xml)
	id = GetNodeText(xmlquery.FindOne(doc, "//*[local-name() = 'CommandId']"))
	return
}
