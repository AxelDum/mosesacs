package cwmp

import (
	"encoding/xml"
	"strings"
	"time"
)

type SoapEnvelope struct {
	XMLName xml.Name
	Header  SoapHeader
	Body    SoapBody
}

type SoapHeader struct{}
type SoapBody struct {
	CWMPMessage CWMPMessage `xml:",any"`
}

type CWMPMessage struct {
	XMLName xml.Name
}

type EventStruct struct {
	EventCode  string
	CommandKey string
}

type ParameterValueStruct struct {
	Name  string
	Value string
}

type ParameterInfoStruct struct {
	Name     string
	Writable string
}

type GetParameterValues_ struct {
	ParameterNames []string `xml:"Body>GetParameterValues>ParameterNames>string"`
}

type GetParameterValuesResponse struct {
	ParameterList []ParameterValueStruct `xml:"Body>GetParameterValuesResponse>ParameterList>ParameterValueStruct"`
}

type GetParameterNamesResponse struct {
	ParameterList []ParameterInfoStruct `xml:"Body>GetParameterNamesResponse>ParameterList>ParameterInfoStruct"`
}

type CWMPInform struct {
	DeviceId      DeviceID               `xml:"Body>Inform>DeviceId"`
	Events        []EventStruct          `xml:"Body>Inform>Event>EventStruct"`
	ParameterList []ParameterValueStruct `xml:"Body>Inform>ParameterList>ParameterValueStruct"`
}

func (s *SoapEnvelope) KindOf() string {
	return s.Body.CWMPMessage.XMLName.Local
}

func (i *CWMPInform) GetEvents() string {
	res := ""
	for idx := range i.Events {
		res += i.Events[idx].EventCode
	}

	return res
}

func (i *CWMPInform) GetConnectionRequest() string {
	for idx := range i.ParameterList {
		// valid condition for both tr98 and tr181
		if strings.HasSuffix(i.ParameterList[idx].Name, "Device.ManagementServer.ConnectionRequestURL") {
			return i.ParameterList[idx].Value
		}
	}

	return ""
}

func (i *CWMPInform) GetSoftwareVersion() string {
	for idx := range i.ParameterList {
		if strings.HasSuffix(i.ParameterList[idx].Name, "Device.DeviceInfo.SoftwareVersion") {
			return i.ParameterList[idx].Value
		}
	}

	return ""
}

func (i *CWMPInform) GetHardwareVersion() string {
	for idx := range i.ParameterList {
		if strings.HasSuffix(i.ParameterList[idx].Name, "Device.DeviceInfo.HardwareVersion") {
			return i.ParameterList[idx].Value
		}
	}

	return ""
}

func (i *CWMPInform) GetDataModelType() string {
	if strings.HasPrefix(i.ParameterList[0].Name, "InternetGatewayDevice") {
		return "TR098"
	} else if strings.HasPrefix(i.ParameterList[0].Name, "Device") {
		return "TR181"
	}

	return ""
}

type DeviceID struct {
	Manufacturer string
	OUI          string
	SerialNumber string
}

func InformResponse() string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soapenc="http://schemas.xmlsoap.org/soap/encoding/" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:cwmp="urn:dslforum-org:cwmp-1-0" xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:schemaLocation="urn:dslforum-org:cwmp-1-0 ..\schemas\wt121.xsd" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
  <soap:Header/>
  <soap:Body soap:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
    <cwmp:InformResponse>
      <MaxEnvelopes>1</MaxEnvelopes>
    </cwmp:InformResponse>
  </soap:Body>
</soap:Envelope>`
}

func GetParameterValues(leaf string) string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soapenc="http://schemas.xmlsoap.org/soap/encoding/" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:cwmp="urn:dslforum-org:cwmp-1-0" xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:schemaLocation="urn:dslforum-org:cwmp-1-0 ..\schemas\wt121.xsd" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
  <soap:Header/>
  <soap:Body soap:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
    <cwmp:GetParameterValues>
      <ParameterNames>
      	<string>` + leaf + `</string>
        <string>Device.Time.</string>
      </ParameterNames>
    </cwmp:GetParameterValues>
  </soap:Body>
</soap:Envelope>`
}

func GetParameterNames(leaf string) string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soapenc="http://schemas.xmlsoap.org/soap/encoding/" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:cwmp="urn:dslforum-org:cwmp-1-0" xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:schemaLocation="urn:dslforum-org:cwmp-1-0 ..\schemas\wt121.xsd" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
  <soap:Header/>
  <soap:Body soap:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
    <cwmp:GetParameterNames>
      <ParameterPath>` + leaf + `</ParameterPath>
      <NextLevel>1</NextLevel>
    </cwmp:GetParameterNames>
  </soap:Body>
</soap:Envelope>`
}

func FactoryReset() string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soapenc="http://schemas.xmlsoap.org/soap/encoding/" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:cwmp="urn:dslforum-org:cwmp-1-0" xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:schemaLocation="urn:dslforum-org:cwmp-1-0 ..\schemas\wt121.xsd" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
  <soap:Header/>
  <soap:Body soap:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
    <cwmp:FactoryReset/>
  </soap:Body>
</soap:Envelope>`
}

// CPE side

func Inform(serial string) string {
	return `<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:soap-enc="http://schemas.xmlsoap.org/soap/encoding/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:cwmp="urn:dslforum-org:cwmp-1-0"><soap:Header><cwmp:ID soap:mustUnderstand="1">5058</cwmp:ID></soap:Header>
	<soap:Body><cwmp:Inform><DeviceId><Manufacturer>ADB Broadband</Manufacturer>
<OUI>0013C8</OUI>
<ProductClass>VV5522</ProductClass>
<SerialNumber>PI234550701S199991-`+ serial +`</SerialNumber>
</DeviceId>
<Event soap-enc:arrayType="cwmp:EventStruct[1]">
<EventStruct><EventCode>6 CONNECTION REQUEST</EventCode>
<CommandKey></CommandKey>
</EventStruct>
</Event>
<MaxEnvelopes>1</MaxEnvelopes>
<CurrentTime>` + time.Now().Format(time.RFC3339) + `</CurrentTime>
<RetryCount>0</RetryCount>
<ParameterList soap-enc:arrayType="cwmp:ParameterValueStruct[8]">
<ParameterValueStruct><Name>InternetGatewayDevice.ManagementServer.ConnectionRequestURL</Name>
<Value xsi:type="xsd:string">http://localhost:7547/ConnectionRequest-`+serial+`</Value>
</ParameterValueStruct>
<ParameterValueStruct><Name>InternetGatewayDevice.ManagementServer.ParameterKey</Name>
<Value xsi:type="xsd:string"></Value>
</ParameterValueStruct>
<ParameterValueStruct><Name>InternetGatewayDevice.DeviceSummary</Name>
<Value xsi:type="xsd:string">InternetGatewayDevice:1.2[](Baseline:1,EthernetLAN:1,WiFiLAN:1,ADSLWAN:1,EthernetWAN:1,QoS:1,QoSDynamicFlow:1,Bridging:1,Time:1,IPPing:1,TraceRoute:1,DeviceAssociation:1,UDPConnReq:1),VoiceService:1.0[1](TAEndpoint:1,SIPEndpoint:1)</Value>
</ParameterValueStruct>
<ParameterValueStruct><Name>InternetGatewayDevice.DeviceInfo.HardwareVersion</Name>
<Value xsi:type="xsd:string">`+serial+`</Value>
</ParameterValueStruct>
<ParameterValueStruct><Name>InternetGatewayDevice.DeviceInfo.ProvisioningCode</Name>
<Value xsi:type="xsd:string">ABCD</Value>
</ParameterValueStruct>
<ParameterValueStruct><Name>InternetGatewayDevice.DeviceInfo.SoftwareVersion</Name>
<Value xsi:type="xsd:string">E_8.0.0.0002</Value>
</ParameterValueStruct>
<ParameterValueStruct><Name>InternetGatewayDevice.DeviceInfo.SpecVersion</Name>
<Value xsi:type="xsd:string">1.0</Value>
</ParameterValueStruct>
<ParameterValueStruct><Name>InternetGatewayDevice.WANDevice.1.WANConnectionDevice.1.WANIPConnection.1.ExternalIPAddress</Name>
<Value xsi:type="xsd:string">12.0.0.10</Value>
</ParameterValueStruct>
</ParameterList>
</cwmp:Inform>
</soap:Body></soap:Envelope>`
}

func BuildGetParameterValuesResponse(serial string) string {
	return `<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:soap-enc="http://schemas.xmlsoap.org/soap/encoding/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:cwmp="urn:dslforum-org:cwmp-1-0">
	<soap:Header><cwmp:ID soap:mustUnderstand="1">3</cwmp:ID></soap:Header>
	<soap:Body><cwmp:GetParameterValuesResponse><ParameterList soap-enc:arrayType="cwmp:ParameterValueStruct[20]">
<ParameterValueStruct><Name>InternetGatewayDevice.DeviceInfo.AdditionalHardwareVersion</Name>
<Value xsi:type="xsd:string"></Value>
</ParameterValueStruct>
<ParameterValueStruct><Name>InternetGatewayDevice.DeviceInfo.AdditionalSoftwareVersion</Name>
<Value xsi:type="xsd:string">E_8.0.0.0002</Value>
</ParameterValueStruct>
<ParameterValueStruct><Name>InternetGatewayDevice.DeviceInfo.Description</Name>
<Value xsi:type="xsd:string"></Value>
</ParameterValueStruct>
<ParameterValueStruct><Name>InternetGatewayDevice.DeviceInfo.HardwareVersion</Name>
<Value xsi:type="xsd:string">VV5522</Value>
</ParameterValueStruct>
<ParameterValueStruct><Name>InternetGatewayDevice.DeviceInfo.Manufacturer</Name>
<Value xsi:type="xsd:string">ADB Broadband</Value>
</ParameterValueStruct>
<ParameterValueStruct><Name>InternetGatewayDevice.DeviceInfo.ManufacturerOUI</Name>
<Value xsi:type="xsd:string">0013C8</Value>
</ParameterValueStruct>
<ParameterValueStruct><Name>InternetGatewayDevice.DeviceInfo.ModelName</Name>
<Value xsi:type="xsd:string">`+serial+`</Value>
</ParameterValueStruct>
<ParameterValueStruct><Name>InternetGatewayDevice.DeviceInfo.ProvisioningCode</Name>
<Value xsi:type="xsd:string">ABCD</Value>
</ParameterValueStruct>
<ParameterValueStruct><Name>InternetGatewayDevice.DeviceInfo.SerialNumber</Name>
<Value xsi:type="xsd:string">PI234550701S199991-VV5522</Value>
</ParameterValueStruct>
<ParameterValueStruct><Name>InternetGatewayDevice.DeviceInfo.SoftwareVersion</Name>
<Value xsi:type="xsd:string">E_8.0.0.0002</Value>
</ParameterValueStruct>
<ParameterValueStruct><Name>InternetGatewayDevice.DeviceInfo.SpecVersion</Name>
<Value xsi:type="xsd:string">1.0</Value>
</ParameterValueStruct>
<ParameterValueStruct><Name>InternetGatewayDevice.DeviceInfo.VendorConfigFileNumberOfEntries</Name>
<Value xsi:type="xsd:unsignedInt">1</Value>
</ParameterValueStruct>
<ParameterValueStruct><Name>InternetGatewayDevice.DeviceInfo.ProductClass</Name>
<Value xsi:type="xsd:string">VV5522</Value>
</ParameterValueStruct>
<ParameterValueStruct><Name>InternetGatewayDevice.DeviceInfo.FirstUseDate</Name>
<Value xsi:type="xsd:dateTime">2013-10-15T15:40:33Z</Value>
</ParameterValueStruct>
<ParameterValueStruct><Name>InternetGatewayDevice.DeviceInfo.VendorConfigFile.1.Name</Name>
<Value xsi:type="xsd:string">multi_user</Value>
</ParameterValueStruct>
<ParameterValueStruct><Name>InternetGatewayDevice.DeviceInfo.VendorConfigFile.1.Version</Name>
<Value xsi:type="xsd:string">E_8.0.0.0002</Value>
</ParameterValueStruct>
<ParameterValueStruct><Name>InternetGatewayDevice.DeviceInfo.VendorConfigFile.1.Date</Name>
<Value xsi:type="xsd:dateTime">Tue Oct 15 15:48:15 UTC 2013</Value>
</ParameterValueStruct>
<ParameterValueStruct><Name>InternetGatewayDevice.DeviceInfo.VendorConfigFile.1.Description</Name>
<Value xsi:type="xsd:string">multi_user</Value>
<Value xsi:type="xsd:string"></Value>
</ParameterValueStruct>
<ParameterValueStruct><Name>InternetGatewayDevice.DeviceInfo.UpTime</Name>
<Value xsi:type="xsd:unsignedInt">5062</Value>
</ParameterValueStruct>
</ParameterList>
</cwmp:GetParameterValuesResponse>
</soap:Body></soap:Envelope>`
}