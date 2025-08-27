package wsmv

import (
	"bytes"
	"fmt"
	"maps"
)

func BuildXML(headers map[string]any, body map[string]any) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")

	// <s:Envelope ...>
	fmt.Fprintf(&buf, "<%s:Envelope", NS_SOAP_ENV)
	for k, v := range Namespaces {
		fmt.Fprintf(&buf, ` %s="%s"`, k, v)
	}
	buf.WriteString(">\n")

	// Header
	buf.WriteString(fmt.Sprintf("  <%s:Header>\n", NS_SOAP_ENV))
	if err := writeMap(&buf, headers, "    "); err != nil {
		return nil, err
	}
	buf.WriteString(fmt.Sprintf("  </%s:Header>\n", NS_SOAP_ENV))

	// Body
	buf.WriteString(fmt.Sprintf("  <%s:Body>\n", NS_SOAP_ENV))
	if err := writeMap(&buf, body, "    "); err != nil {
		return nil, err
	}
	buf.WriteString(fmt.Sprintf("  </%s:Body>\n", NS_SOAP_ENV))

	// </s:Envelope>
	buf.WriteString(fmt.Sprintf("</%s:Envelope>", NS_SOAP_ENV))

	return buf.Bytes(), nil
}

func writeMap(buf *bytes.Buffer, data map[string]any, indent string) error {
	for key, val := range data {
		if key == ":attributes!" {
			continue
		}

		switch v := val.(type) {
		case map[string]any:
			if err := writeComplexElement(buf, key, v, indent); err != nil {
				return err
			}
		case []map[string]any:
			for _, item := range v {
				if err := writeComplexElement(buf, key, item, indent); err != nil {
					return err
				}
			}
		default:
			writeSimpleElement(buf, key, fmt.Sprintf("%v", v), data, indent)
		}
	}
	return nil
}

func writeComplexElement(buf *bytes.Buffer, tag string, content map[string]any, indent string) error {
	attrs := extractAttributes(tag, content, nil)
	buf.WriteString(indent)
	buf.WriteString("<" + tag)
	writeAttrs(buf, attrs)
	buf.WriteString(">")

	if text, ok := content["_"]; ok {
		fmt.Fprintf(buf, "%v", text)
	} else {
		buf.WriteString("\n")
		if err := writeMap(buf, content, indent+"  "); err != nil {
			return err
		}
		buf.WriteString(indent)
	}

	fmt.Fprintf(buf, "</%s>\n", tag)
	return nil
}

func writeSimpleElement(buf *bytes.Buffer, tag, value string, parent map[string]any, indent string) {
	attrs := extractAttributes(tag, nil, parent)
	buf.WriteString(indent)
	buf.WriteString("<" + tag)
	writeAttrs(buf, attrs)
	buf.WriteString(">")
	buf.WriteString(value)
	fmt.Fprintf(buf, "</%s>\n", tag)
}

func writeAttrs(buf *bytes.Buffer, attrs map[string]any) {
	for k, v := range attrs {
		fmt.Fprintf(buf, ` %s="%v"`, k, v)
	}
}

func extractAttributes(tag string, self map[string]any, parent map[string]any) map[string]any {
	out := map[string]any{}

	if parent != nil {
		if parentAttrs, ok := parent[":attributes!"].(map[string]any); ok {
			if attrs, ok := parentAttrs[tag].(map[string]any); ok {
				maps.Copy(out, attrs)
			}
		}
	}

	if self != nil {
		if attrs, ok := self[":attributes!"].(map[string]any); ok {
			maps.Copy(out, attrs)
		}
	}

	return out
}
