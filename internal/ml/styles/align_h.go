package styles

import (
	"encoding/xml"
)

//HAlignType is a type to encode XSD ST_HorizontalAlignment
type HAlignType byte

var (
	ToHAlignType   map[string]HAlignType
	FromHAlignType map[HAlignType]string
)

func (t HAlignType) String() string {
	return FromHAlignType[t]
}

func (t *HAlignType) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	attr := xml.Attr{Name: name}

	if v, ok := FromHAlignType[*t]; ok {
		attr.Value = v
	} else {
		attr = xml.Attr{}
	}

	return attr, nil
}

func (t *HAlignType) UnmarshalXMLAttr(attr xml.Attr) error {
	if v, ok := ToHAlignType[attr.Value]; ok {
		*t = v
	}

	return nil
}
