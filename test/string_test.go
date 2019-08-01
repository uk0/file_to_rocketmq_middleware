package main

import (
	"bufio"
	"fmt"
	. "github.com/clbanning/mxj"
	"io"
	"os"
	"strings"
	"testing"
)

func GetFileName(path string)string{
	arrayStr := strings.Split(path,"/")
	return arrayStr[len(arrayStr)-1]
}

func TestSplit(test2 *testing.T){
		fmt.Println(GetFileName("/Users/zhangjianxin/home/GO_LIB/src/github.com/uk0/readCSV/test.json"))
}

func TestString2(test *testing.T)  {
	fmt.Println("start test")
	fi, err := os.Open("/Users/zhangjianxin/home/GO_LIB/src/github.com/uk0/readCSV/test.json")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer fi.Close()

	br := bufio.NewReader(fi)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		var dataLine  = string([]rune(string(a)));
		fmt.Println(DeleteSymbol2(dataLine))
	}


	//var json = `"{""tracesList"":{""tracesElement"":[{""logisticProviderID"":""PJ"",""mailNo"":""aaa"",""traces"":{""trace"":{""time"":""2019-07-21 08:00:36"",""desc"":""快件在【达州文星站】已操作进站扫描，操作员：【aaa】"",""city"":""达州市"",""facilityType"":""1"",""facilityNo"":""aaaa"",""facilityName"":""aaa"",""action"":""1"",""contacter"":"""",""contactPhone"":"""",""next"":{""city"":""达州市"",""facilityType"":""1"",""facilityNo"":"""",""facilityName"":""""},""previous"":{""city"":"""",""facilityType"":"""",""facilityNo"":"""",""facilityName"":""""}}}}]}}"`
	//var xml = `"<?xml version=""1.0"" encoding=""utf-8""?><tracesList><tracesElement><logisticProviderID>STO</logisticProviderID><mailNo>aaa</mailNo><traces><trace><time>2019-07-24 09:55:59</time><desc>【浙江杭州乔司公司】已进行【aa】扫描,【aa】</desc><city>aaaa</city><facilityType>1</facilityType><facilityNo>aaaa</facilityNo><facilityName>浙江杭州乔司公司</facilityName><action>22</action><contacter></contacter><contactPhone></contactPhone><next><city></city><facilityType></facilityType><facilityNo></facilityNo><facilityName></facilityName></next><previous><city></city><facilityType></facilityType><facilityNo></facilityNo><facilityName></facilityName></previous></trace></traces></tracesElement></tracesList>"`

	//fmt.Println(DeleteSymbol2(json))
	//fmt.Println(DeleteSymbol2(xml))
}

const HEADER  = `<?xml version="1.0" encoding="utf-8"?>`

// FromJson() --> ToXml().
func JsonToXml2(jsonVal []byte) ([]byte, error) {
	m, err := NewMapJson(jsonVal)
	if err != nil {
		return nil, err
	}
	return m.Xml()
}

func DeleteSymbol2(line string) string {
	// json
	if strings.HasSuffix(line, "}\"") && strings.HasPrefix(line, "\"{") {
		fmt.Println("is Json")
		var tempJSONLine = strings.Replace(line, "}\"", "}", -1)
		tempJSONLine = strings.Replace(tempJSONLine, "\"{", "{", -1)
		tempJSONLine = strings.Replace(tempJSONLine, "\"\"", "\"", -1)
		r,_:=JsonToXml2([]byte(tempJSONLine))

		return HEADER + string(r)
	}
	// xml
	if strings.HasSuffix(line, ">\"") && strings.HasPrefix(line, "\"<") {
		fmt.Println("is Xml")
		var tempXmlLine = strings.Replace(line, ">\"", ">", -1)
		tempXmlLine = strings.Replace(tempXmlLine, "\"<", "<", -1)
		tempXmlLine = strings.Replace(tempXmlLine, "\"\"", "\"", -1)
		return tempXmlLine
	}
	return ""
}