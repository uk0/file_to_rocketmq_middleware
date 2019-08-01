package main

import (
	. "github.com/clbanning/mxj"
	"log"
	"strings"
)

const HEADER  = `<?xml version="1.0" encoding="utf-8"?>`

// FromJson() --> ToXml().
func JsonToXml(jsonVal []byte) ([]byte, error) {
	m, err := NewMapJson(jsonVal)
	if err != nil {
		return nil, err
	}
	return m.Xml()
}

func DeleteSymbol(line string) string {
	// 特殊json
	if strings.HasSuffix(line, "}\"") && strings.HasPrefix(line, "\"{") {
		log.Println("is Unusual Json")
		var tempJSONLine = strings.Replace(line, "}\"", "}", -1)
		tempJSONLine = strings.Replace(tempJSONLine, "\"{", "{", -1)
		tempJSONLine = strings.Replace(tempJSONLine, "\"\"", "\"", -1)
		r,_:=JsonToXml([]byte(tempJSONLine))

		return HEADER + string(r)
	}
	// 正常json
	if strings.HasSuffix(line, "}") && strings.HasPrefix(line, "{") {
		log.Println("is normal Json")
		r,_:=JsonToXml([]byte(line))
		return HEADER + string(r)
	}
	// xml
	if strings.HasSuffix(line, ">\"") && strings.HasPrefix(line, "\"<") {
		log.Println("is Xml")
		var tempXmlLine = strings.Replace(line, ">\"", ">", -1)
		tempXmlLine = strings.Replace(tempXmlLine, "\"<", "<", -1)
		tempXmlLine = strings.Replace(tempXmlLine, "\"\"", "\"", -1)
		return tempXmlLine
	}
	// 如果不满足直接发出去
	return ""
}