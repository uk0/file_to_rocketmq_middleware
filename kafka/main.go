package main

import (
	"fmt"
	"sync"
)

var waitgroup sync.WaitGroup

func main() {

	_ = InitKafka()
	go JSONBuilder()

	waitgroup.Wait()
}

func JSONBuilder() {
	for {
		waitgroup.Add(1)
		JSONTempLate_OP := `<?xml version="1.0" encoding="GBK"?><tracesList><tracesElement><logisticProviderID>YTO</logisticProviderID><mailNo>700383296498</mailNo><traces><trace><time>2017-09-27 13:57:02</time><desc/><city/><facilityType>BRANCH</facilityType><facilityNo>834006</facilityNo><facilityName/><action>171</action><contacter/><contactPhone/><next><city/><facilityType/><facilityNo></facilityNo><facilityName/></next><previous><city/><facilityType/><facilityNo></facilityNo><facilityName/></previous></trace></traces></tracesElement></tracesList>`
		JSONTempLate_ZT := fmt.Sprint(`<?xml version="1.0" encoding="GBK"?><tracesList><tracesElement><logisticProviderID>YTO</logisticProviderID><mailNo>700383296498</mailNo><traces><trace><time>2017-09-27 13:57:02</time><desc/><city/><facilityType>BRANCH</facilityType><facilityNo>834006</facilityNo><facilityName/><action>171</action><contacter/><contactPhone/><next><city/><facilityType/><facilityNo></facilityNo><facilityName/></next><previous><city/><facilityType/><facilityNo></facilityNo><facilityName/></previous></trace></traces></tracesElement></tracesList>`)

		_ = kafkaSender.addMessage(JSONTempLate_OP, "test_op")
		_ = kafkaSender.addMessage(JSONTempLate_ZT, "test_zt")
		waitgroup.Done()
	}
}
