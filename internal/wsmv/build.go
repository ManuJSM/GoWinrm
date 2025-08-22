package wsmv

import (
	"bytes"
	"encoding/xml"
	"fmt"
)

// BuildXML genera el XML del mensaje WSMV
func BuildXML(headers map[string]any, body map[string]any) (string, error) {
	var buf bytes.Buffer
	enc := xml.NewEncoder(&buf)
	enc.Indent("", "  ")
	buf.WriteString(xml.Header)

	// <s:Envelope ...namespaces>
	startElem := xml.StartElement{
		Name: xml.Name{Local: NS_SOAP_ENV + ":Envelope"},
		Attr: []xml.Attr{},
	}

	// Añadimos atributos xmlns para namespaces
	for prefix, uri := range Namespaces {
		startElem.Attr = append(startElem.Attr, xml.Attr{
			Name:  xml.Name{Local: prefix},
			Value: uri,
		})
	}

	if err := enc.EncodeToken(startElem); err != nil {
		return "", err
	}

	// Crear header
	if len(headers) > 0 {
		if err := createHeader(enc, headers); err != nil {
			return "", err
		}
	}

	// Crear body
	if len(body) > 0 {
		if err := createBody(enc, body); err != nil {
			return "", err
		}
	}

	// Cerrar Envelope
	if err := enc.EncodeToken(xml.EndElement{Name: xml.Name{Local: NS_SOAP_ENV + ":Envelope"}}); err != nil {
		return "", err
	}

	if err := enc.Flush(); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// createHeader crea el elemento Header
func createHeader(enc *xml.Encoder, headers map[string]any) error {
	startElem := xml.StartElement{Name: xml.Name{Local: NS_SOAP_ENV + ":Header"}}
	if err := enc.EncodeToken(startElem); err != nil {
		return err
	}

	if err := encodeContent(enc, headers); err != nil {
		return err
	}

	return enc.EncodeToken(xml.EndElement{Name: xml.Name{Local: NS_SOAP_ENV + ":Header"}})
}

// createBody crea el elemento Body
func createBody(enc *xml.Encoder, body map[string]any) error {
	startElem := xml.StartElement{Name: xml.Name{Local: NS_SOAP_ENV + ":Body"}}
	if err := enc.EncodeToken(startElem); err != nil {
		return err
	}

	if err := encodeContent(enc, body); err != nil {
		return err
	}

	return enc.EncodeToken(xml.EndElement{Name: xml.Name{Local: NS_SOAP_ENV + ":Body"}})
}

// encodeContent codifica el contenido de un mapa
func encodeContent(enc *xml.Encoder, data map[string]any) error {
	for key, value := range data {
		if key == ":attributes!" {
			continue
		}

		if err := encodeElement(enc, key, value, data); err != nil {
			return err
		}
	}
	return nil
}

// encodeElement codifica un elemento individual
func encodeElement(enc *xml.Encoder, tagName string, value any, parentData map[string]any) error {
	switch v := value.(type) {
	case []map[string]any:
		// Slice de elementos complejos
		for _, item := range v {
			if err := encodeComplexElement(enc, tagName, item, parentData); err != nil {
				return err
			}
		}
		return nil

	case map[string]any:
		// Elemento complejo único
		return encodeComplexElement(enc, tagName, v, parentData)

	default:
		// Valor simple
		return encodeSimpleElement(enc, tagName, value, parentData)
	}
}

// encodeComplexElement codifica un elemento complejo
func encodeComplexElement(enc *xml.Encoder, tagName string, item map[string]any, parentData map[string]any) error {
	// Obtener atributos
	attrs := getAttributes(tagName, item, parentData)

	startElem := xml.StartElement{
		Name: xml.Name{Local: tagName},
		Attr: attrs,
	}

	if err := enc.EncodeToken(startElem); err != nil {
		return err
	}

	// Manejar contenido especial "_" que representa el texto del elemento
	if content, ok := item["_"]; ok {
		if err := enc.EncodeToken(xml.CharData(fmt.Sprintf("%v", content))); err != nil {
			return err
		}
	} else {
		// Codificar elementos hijos
		for key, value := range item {
			if key == ":attributes!" || key == "_" {
				continue
			}
			if err := encodeElement(enc, key, value, item); err != nil {
				return err
			}
		}
	}

	return enc.EncodeToken(xml.EndElement{Name: xml.Name{Local: tagName}})
}

// encodeSimpleElement codifica un elemento simple
func encodeSimpleElement(enc *xml.Encoder, tagName string, value any, parentData map[string]any) error {
	// Obtener atributos del padre para este elemento
	var attrs []xml.Attr
	if parentAttrs, ok := parentData[":attributes!"].(map[string]any); ok {
		if elemAttrs, ok := parentAttrs[tagName].(map[string]any); ok {
			for attrName, attrValue := range elemAttrs {
				attrs = append(attrs, createAttr(attrName, attrValue))
			}
		}
	}

	startElem := xml.StartElement{
		Name: xml.Name{Local: tagName},
		Attr: attrs,
	}

	if err := enc.EncodeToken(startElem); err != nil {
		return err
	}

	if err := enc.EncodeToken(xml.CharData(fmt.Sprintf("%v", value))); err != nil {
		return err
	}

	return enc.EncodeToken(xml.EndElement{Name: xml.Name{Local: tagName}})
}

// getAttributes obtiene los atributos para un elemento
func getAttributes(tagName string, item map[string]any, parentData map[string]any) []xml.Attr {
	var attrs []xml.Attr

	// Atributos del padre
	if parentAttrs, ok := parentData[":attributes!"].(map[string]any); ok {
		if elemAttrs, ok := parentAttrs[tagName].(map[string]any); ok {
			for attrName, attrValue := range elemAttrs {
				attrs = append(attrs, createAttr(attrName, attrValue))
			}
		}
	}

	// Atributos del elemento mismo
	if itemAttrs, ok := item[":attributes!"].(map[string]any); ok {
		for attrName, attrValue := range itemAttrs {
			attrs = append(attrs, createAttr(attrName, attrValue))
		}
	}

	return attrs
}

// createAttr crea un atributo XML
func createAttr(attrName string, attrValue any) xml.Attr {
	// Manejar atributos booleanos (mustUnderstand)
	if b, ok := attrValue.(bool); ok {
		if b {
			attrValue = "true"
		} else {
			attrValue = "false"
		}
	}

	// Atributo sin namespace
	return xml.Attr{
		Name:  xml.Name{Local: attrName},
		Value: fmt.Sprintf("%v", attrValue),
	}
}
