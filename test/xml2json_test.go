package main

import (
	"bytes"
	"container/list"
	"encoding/xml"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestXMLTOJSON(test *testing.T) {
	strs := "<tracesList>\n" +
		" <tracesElement>\n" +
		"  <mailNo>TT6600189673233</mailNo>\n" +
		"  <logisticProviderID>TTKD</logisticProviderID>\n" +
		"  <traces>\n" +
		"   <trace>\n" +
		"   <city>北京</city>\n" +
		"    <previous>\n" +
		"     <facilityName>previous_facilityName</facilityName>\n" +
		"     <facilityType>previous_facilityType</facilityType>\n" +
		"     <facilityNo>previous_facilityNo</facilityNo>\n" +
		"     <city>previous_city</city>\n" +
		"    </previous>\n" +
		"    <next>\n" +
		"     <facilityName>facilityName_next小明</facilityName>\n" +
		"     <facilityType>facilityType_next</facilityType>\n" +
		"     <facilityNo>facilityNo_next</facilityNo>\n" +
		"     <city>city_next</city>\n" +
		"    </next>\n" +
		"    <contacter>2325</contacter>\n" +
		"    <action>ARRIVAL</action>\n" +
		"    <facilityNo>L999</facilityNo>\n" +
		"    \n" +
		"   </trace>\n" +
		"  </traces>\n" +
		" </tracesElement>\n" +
		"</tracesList>"

	xml_str2 := "<tracesList>\n" +
		" <tracesElement>\n" +
		"  <mailNo>TT6600189673233</mailNo>\n" +
		"  <logisticProviderID>TTKD</logisticProviderID>\n" +
		"  <traces>\n" +
		"   <trace>\n" +
		"    <previous>\n" +
		"     <facilityName>previous_facilityName</facilityName>\n" +
		"     <facilityType>previous_facilityType</facilityType>\n" +
		"     <facilityNo>previous_facilityNo</facilityNo>\n" +
		"     <city>previous_city</city>\n" +
		"    </previous>\n" +
		"    <city>北京</city>\n" +
		"    <next>\n" +
		"     <facilityName>next_facilityName</facilityName>\n" +
		"     <facilityType>next_facilityType</facilityType>\n" +
		"     <facilityNo>next_facilityNo</facilityNo>\n" +
		"     <city>next_city</city>\n" +
		"    </next>\n" +
		"    <contacter>2325</contacter>\n" +
		"    <action>ARRIVAL</action>\n" +
		"    <facilityNo>L999</facilityNo>\n" +
		"   </trace>\n" +
		"  </traces>\n" +
		" </tracesElement>\n" +
		"</tracesList>";

	t := time.Now()

	l:=OPToJson([]byte(strs))
	l2:=OPToJson([]byte(xml_str2))

	fmt.Println(l.Back())
	fmt.Println(l2.Back())
	elapsed := time.Since(t)

	fmt.Println("app elapsed:", elapsed)

}

func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}

var flag2 = true;

var flag1 = true;

func OPToJson(data []byte) (*list.List) {

	cityString := []string{"city", "nextCity", "previousCity"}
	facilityNos := []string{"facilityNo", "nextFacilityNo", "previousFacilityNo"}
	facilityTypes := []string{"facilityType", "nextFacilityType", "previousFacilityType"}
	facilityNames := []string{"facilityName", "nextFacilityName", "previousFacilityName"}
	names := []string{"logisticProviderID", "mailNo", "time", "desc", "city", "facilityType", "facilityNo", "facilityName", "action", "contactPhone", "contacter", "staffCode", "operateType", "extendedInfo"}
	var tagName = ""

	var prevCity = 0
	var nextCity = 0

	var nextFacilityType = 0
	var prevFacilityType = 0

	var nextFacilityNo = 0
	var prevFacilityNo = 0

	var prevFacilityName = 0
	var nextFacilityName = 0
	list := list.New()

	xmlFile := bytes.NewReader(data)

	decoder := xml.NewDecoder(xmlFile)
	var inElement string
	var sb strings.Builder
	var ssb strings.Builder

	for {
		// Read tokens from the XML document in a stream.
		t, _ := decoder.Token()
		if t == nil {
			break
		}
		switch se := t.(type) {

		case xml.StartElement:
			// If we just read a StartElement token
			inElement = se.Name.Local

			if strings.EqualFold(inElement, "tracesElement") {
				sb.WriteString("{")
			} else if strings.EqualFold(inElement, "traces") || strings.EqualFold(inElement, "trace") {
				prevCity = 0
				nextCity = 0

				nextFacilityType = 0
				prevFacilityType = 0

				nextFacilityNo = 0
				prevFacilityNo = 0

				prevFacilityName = 0
				nextFacilityName = 0

			} else if strings.EqualFold(inElement, "next") {
				nextCity = 1
				nextFacilityType = 1
				nextFacilityNo = 1
				nextFacilityName = 1
			} else if strings.EqualFold(inElement, "previous") {
				prevCity = 1
				prevFacilityNo = 1
				prevFacilityType = 1
				prevFacilityName = 1

				// 分段
			} else if prevCity == 1 && strings.EqualFold(inElement, "city") {
				sb.WriteString("\"")
				sb.WriteString(cityString[2])
				sb.WriteString("\":\"")
				prevCity = 0
			} else if nextCity == 1 && strings.EqualFold(inElement, "city") {
				sb.WriteString("\"")
				sb.WriteString(cityString[1])
				sb.WriteString("\":\"")
				nextCity = 0
			} else if strings.EqualFold(inElement, "city") {
				sb.WriteString("\"")
				sb.WriteString(cityString[0])
				sb.WriteString("\":\"")

				// 分段解析
			} else if nextFacilityNo == 1 && strings.EqualFold(inElement, "facilityNo") {
				sb.WriteString("\"")
				sb.WriteString(facilityNos[1])
				sb.WriteString("\":\"")
				nextFacilityNo = 0
			} else if prevFacilityNo == 1 && strings.EqualFold(inElement, "facilityNo") {
				sb.WriteString("\"")
				sb.WriteString(facilityNos[2])
				sb.WriteString("\":\"")
				prevFacilityNo = 0
			} else if strings.EqualFold(inElement, "facilityNo") {
				sb.WriteString("\"")
				sb.WriteString(facilityNos[0])
				sb.WriteString("\":\"")
				// 分段

			} else if nextFacilityType == 1 && strings.EqualFold(inElement, "facilityType") {
				sb.WriteString("\"")
				sb.WriteString(facilityTypes[1])
				sb.WriteString("\":\"")
				nextFacilityType = 0
			} else if prevFacilityType == 1 && strings.EqualFold(inElement, "facilityType") {
				sb.WriteString("\"")
				sb.WriteString(facilityTypes[2])
				sb.WriteString("\":\"")
				prevFacilityType = 0
			} else if strings.EqualFold(inElement, "facilityType") {
				sb.WriteString("\"")
				sb.WriteString(facilityTypes[0])
				sb.WriteString("\":\"")
				// 分段

			} else if nextFacilityName == 1 && strings.EqualFold(inElement, "facilityName") {
				sb.WriteString("\"")
				sb.WriteString(facilityNames[1])
				sb.WriteString("\":\"")
				nextFacilityName = 0
			} else if prevFacilityName == 1 && strings.EqualFold(inElement, "facilityName") {
				sb.WriteString("\"")
				sb.WriteString(facilityNames[2])
				sb.WriteString("\":\"")
				prevFacilityName = 0
			} else if strings.EqualFold(inElement, "facilityName") {
				sb.WriteString("\"")
				sb.WriteString(facilityNames[0])
				sb.WriteString("\":\"")
				// 分段
			} else if !strings.EqualFold(inElement, "tracesList") && !strings.EqualFold(inElement, "") && inElement != "" {
				sb.WriteString("\"")
				sb.WriteString(inElement)
				sb.WriteString("\":\"")
			}

		case xml.CharData:
			sb.WriteString(fmt.Sprintf("%s",string(se.Copy())));
		case xml.EndElement:
			flag2 = false
			tagName = se.Name.Local
			if flag1 && contains(names, tagName) {
				sb.WriteString("\",")
			}
			if !flag1 && contains(names, tagName) {
				sb.WriteString("\"")
				sb.WriteString(strings.Replace(strings.Replace(ssb.String(), "\\", "\\\\", -1), "\"", "", -1))
				sb.WriteString("\",");
			}
			if strings.EqualFold(tagName, "tracesElement") {
				resp := sb.String()
				sb.Reset()
				resp = strings.Replace(resp, " ", "", -1)
				// 去除换行符
				resp = strings.Replace(resp, "\n", "", -1)
				resp = string(resp[:len(resp)-2])
				sb.WriteString(resp)
				sb.WriteString("}")
				list.PushBack(sb.String())
			}
			tagName = ""

		default:
		}

	}
	return list
}