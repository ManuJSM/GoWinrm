package shell

import (
	"GoWinrm/internal/utils"
	"encoding/base64"
	"strconv"
	"strings"

	"github.com/antchfx/xmlquery"
)

type COHandler struct {
	Stdout   strings.Builder
	Stderr   strings.Builder
	ExitCode int
}

func NewCOHandler() *COHandler {
	return &COHandler{}
}

func isCommandDone(respDoc *xmlquery.Node) bool {
	if respDoc == nil {
		return false
	}
	nodes := xmlquery.Find(respDoc, "//*[@State='http://schemas.microsoft.com/wbem/wsman/1/windows/shell/CommandState/Done']")
	return len(nodes) > 0
}

func (coh *COHandler) HandleOutput(xml []byte) (finish bool) {

	outputs := make(map[string][]string)

	doc, _ := utils.ParseXML(string(xml))
	nodes := xmlquery.Find(doc, "/*[local-name()='Envelope']/*[local-name()='Body']/*[local-name()='ReceiveResponse']/*[local-name()='Stream']")
	finish = isCommandDone(doc)

	if finish {
		coh.ExitCode = getExitCode(doc)
	}

	for _, node := range nodes {
		if node.InnerText() == "" {
			continue
		}

		nameAttr := node.SelectAttr("Name")
		if nameAttr == "" {
			continue
		}
		output := utils.GetNodeText(node)
		outputs[nameAttr] = append(outputs[nameAttr], output)

	}

	for _, o := range outputs["stdout"] {
		data, _ := base64.StdEncoding.DecodeString(o)
		coh.Stdout.Write(data)
	}
	for _, o := range outputs["stderr"] {
		data, _ := base64.StdEncoding.DecodeString(o)
		coh.Stderr.Write(data)
	}

	return

}

func getExitCode(doc *xmlquery.Node) int {
	node := xmlquery.FindOne(doc, "//*[local-name()='ExitCode']")
	if node != nil {
		if code, err := strconv.Atoi(node.InnerText()); err == nil {
			return code
		}
	}
	return 69 //Esto no deberia pasar nunca jajajaja
}
