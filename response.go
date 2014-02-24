package winrm

import (
	"encoding/xml"
    "io"
    "io/ioutil"
)

type Selector struct{
    Value string  `xml:",innerxml"`
    Name  string  `xml:"Name,attr"`
}

type SelectorSet struct{
    Selector    *Selector      `xml"w:Selector"`
}

type ReferenceParameters struct{
    ResourceURI     string         `xml"w:ResourceURI"`
    SelectorSet     *SelectorSet   `xml"w:SelectorSet"`
}

type ResourceCreated struct{
    Address                 string                 `xml"a:Address"`
    ReferenceParameters     *ReferenceParameters   `xml"a:ReferenceParameters"`
}

type Header struct{
    Action      string      `xml"a:Action"`
    MessageID   string      `xml"a:MessageID"`
    To          string      `xml"a:To"`
    RelatesTo   string      `xml"a:RelatesTo"`
}

type Shell struct{
    // xmlnsRsp           string  `xml:"xmlns:rsp,attr"`
    ShellId            string  `xml"rsp:ShellId"`
    ResourceUri        string  `xml"rsp:ResourceUri"`
    Owner              string  `xml"rsp:Owner"`
    ClientIP           string  `xml"rsp:ClientIP"`
    IdleTimeOut        string  `xml"rsp:IdleTimeOut"`
    InputStreams       string  `xml"rsp:InputStreams"`
    OutputStreams      string  `xml"rsp:OutputStreams"`
    ShellRunTime       string  `xml"rsp:OutputStreams"`
    ShellInactivity    string  `xml"rsp:OutputStreams"`
}

type Body struct{
    ResourceCreated *ResourceCreated   `xml"x:ResourceCreated"`
    Shell           *Shell             `xml"rsp:Shell"`
}

type ResponseEnvelope struct{
    XMLName     xml.Name `xml"s:Envelope"`
    // xmlnsS      string   `xml"xmlns:s,attr"`
    // xmlnsA      string   `xml"xmlns:a,attr"`
    // xmlnsX      string   `xml"xmlns:x,attr"`
    // xmlnsW      string   `xml"xmlns:w,attr"`
    // xmlnsRsp    string   `xml"xmlns:rsp,attr"`
    // xmlnsP      string   `xml"xmlns:p,attr"`
    // xmlnsLang   string   `xml"xmlns:lang,attr"`
    Header      *Header  `xml"s:Header"`
    Body        *Body    `xml"s:Body"`
}


func GetObjectFromXML(XMLinput io.Reader) (ResponseEnvelope) {    
    b, _ := ioutil.ReadAll(XMLinput)
    var response ResponseEnvelope
    xml.Unmarshal(b, &response)
    return response
}