package sbexml

import "encoding/xml"

func newMessageSchema() messageSchema {
	return messageSchema{
		XMLNSSBE:        "http://fixprotocol.io/2016/sbe",
		ID:              1,
		Version:         0,
		SemanticVersion: "0.0.1",
	}
}

type messageSchema struct {
	XMLName         xml.Name      `xml:"sbe:messageSchema"`
	XMLNSSBE        string        `xml:"xmlns:sbe,attr"`
	Package         string        `xml:"package,attr,omitempty"`
	ID              int           `xml:"id,attr"`
	Version         int           `xml:"version,attr"`
	SemanticVersion string        `xml:"semanticVersion,attr"`
	Types           typeSection   `xml:"types"`
	Messages        []messageType `xml:"message"`
}

type typeSection struct {
	Primitives []primitiveType `xml:"type"`
	Composites []compositeType `xml:"composite"`
	Enums      []enumType      `xml:"enum"`
}

type primitiveType struct {
	Name              string `xml:"name,attr"`
	PrimitiveType     string `xml:"primitiveType,attr"`
	CharacterEncoding string `xml:"characterEncoding,attr,omitempty"`
}

type compositeType struct {
	Name   string          `xml:"name,attr"`
	Types  []primitiveType `xml:"type"`
	Fields []fieldTypeXML  `xml:"field"`
}

type enumType struct {
	Name         string       `xml:"name,attr"`
	EncodingType string       `xml:"encodingType,attr"`
	Values       []validValue `xml:"validValue"`
}

type validValue struct {
	Name  string `xml:"name,attr"`
	Value int32  `xml:",chardata"`
}

type messageType struct {
	Name   string         `xml:"name,attr"`
	ID     int            `xml:"id,attr"`
	Fields []fieldTypeXML `xml:"field"`
	Groups []groupType    `xml:"group"`
}

type fieldTypeXML struct {
	Name         string `xml:"name,attr"`
	ID           int    `xml:"id,attr"`
	Type         string `xml:"type,attr"`
	Presence     string `xml:"presence,attr,omitempty"`
	SemanticType string `xml:"semanticType,attr,omitempty"`
	Epoch        string `xml:"epoch,attr,omitempty"`
	TimeUnit     string `xml:"timeUnit,attr,omitempty"`
}

type groupType struct {
	Name          string         `xml:"name,attr"`
	ID            int            `xml:"id,attr"`
	DimensionType string         `xml:"dimensionType,attr"`
	Fields        []fieldTypeXML `xml:"field"`
}
