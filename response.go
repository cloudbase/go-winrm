package winrm

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"io"
	"io/ioutil"
)

type ResponseSelector struct {
	Value string `xml:",innerxml"`
	Name  string `xml:"Name,attr"`
}

type ResponseSelectorSet struct {
	Selector *ResponseSelector `xml"w:Selector"`
}

type ReferenceParameters struct {
	ResourceURI string               `xml"w:ResourceURI"`
	SelectorSet *ResponseSelectorSet `xml"w:SelectorSet"`
}

type ResourceCreated struct {
	Address             string               `xml"a:Address"`
	ReferenceParameters *ReferenceParameters `xml"a:ReferenceParameters"`
}

type ResponseHeader struct {
	Action    string `xml"a:Action"`
	MessageID string `xml"a:MessageID"`
	To        string `xml"a:To"`
	RelatesTo string `xml"a:RelatesTo"`
}

type CommandResponse struct {
	CommandId string `xml"rsp:CommandId"`
}

type ResponseShell struct {
	// xmlnsRsp        string  `xml:"xmlns:rsp,attr"`
	ShellId         string `xml"rsp:ShellId`
	ResourceUri     string `xml"rsp:ResourceUri"`
	Owner           string `xml"rsp:Owner"`
	ClientIP        string `xml"rsp:ClientIP"`
	IdleTimeOut     string `xml"rsp:IdleTimeOut"`
	InputStreams    string `xml"rsp:InputStreams"`
	OutputStreams   string `xml"rsp:OutputStreams"`
	ShellRunTime    string `xml"rsp:OutputStreams"`
	ShellInactivity string `xml"rsp:OutputStreams"`
}

type ResponseStream struct {
	Value string `xml:",innerxml"`
	Name  string `xml:"Name,attr"`
	End   string `xml:"End,attr"`
}

type ResponseCommandState struct {
	ExitCode  int    `xml"rsp:ExitCode"`
	CommandId string `xml:"CommandId,attr"`
	State     string `xml:"State,attr"`
}

type ReceiveResponse struct {
	Stream       []ResponseStream      `xml"rsp:Stream"`
	CommandState *ResponseCommandState `xml"rsp:CommandState"`
}

type ResponseBody struct {
	CommandResponse *CommandResponse `xml"rsp:CommandResponse"`
	ResourceCreated *ResourceCreated `xml"x:ResourceCreated"`
	Shell           *ResponseShell   `xml"rsp:Shell"`
	ReceiveResponse *ReceiveResponse `xml"rsp:ReceiveResponse"`
}

type ResponseEnvelope struct {
	XMLName xml.Name `xml"s:Envelope"`
	// xmlnsS      string       `xml"xmlns:s,attr"`
	// xmlnsA      string       `xml"xmlns:a,attr"`
	// xmlnsX      string       `xml"xmlns:x,attr"`
	// xmlnsW      string       `xml"xmlns:w,attr"`
	// xmlnsRsp    string       `xml"xmlns:rsp,attr"`
	// xmlnsP      string       `xml"xmlns:p,attr"`
	// xmlnsLang   string       `xml"xmlns:lang,attr"`
	Header *ResponseHeader `xml"s:Header"`
	Body   *ResponseBody   `xml"s:Body"`
}

func GetObjectFromXML(XMLinput io.Reader) (ResponseEnvelope, error) {
	b, err := ioutil.ReadAll(XMLinput)
	var response ResponseEnvelope
	if err != nil {
		return response, err
	} else {
		xml.Unmarshal(b, &response)
	}
	x := ResponseEnvelope{}
	if response == x {
		return x, errors.New("Invalid server response")
	}
	return response, nil
}

var parseXML = GetObjectFromXML

func ParseCommandOutput(XMLinput io.Reader) (stdout, stderr string, exitcode int, err error) {
	//fmt.Printf("%s\n\n\n", XMLinput)
	object, err := parseXML(XMLinput)
	//fmt.Printf("%s", object.Body.ReceiveResponse.Stream)
	if err != nil {
		return "", "", 0, errors.New("Error parsing XML")
	}
	var stdout_b bytes.Buffer
	var stderr_b bytes.Buffer
	for _, value := range object.Body.ReceiveResponse.Stream {
		if value.End == "" && value.Name == "stdout" {
			tmp, err := base64.StdEncoding.DecodeString(value.Value)
			if err != nil {
				return "", "", 0, errors.New("Error decoding stdout")
			}
			stdout_b.Write(tmp)
		} else if value.End == "" && value.Name == "stderr" {
			tmp, err := base64.StdEncoding.DecodeString(value.Value)
			if err != nil {
				return "", "", 0, errors.New("Error decoding stderr")
			}
			stderr_b.Write(tmp)
		} else {
			break
		}
	}
	exitcode = object.Body.ReceiveResponse.CommandState.ExitCode
	stdout = stdout_b.String()
	stderr = stderr_b.String()
	err = nil
	return
}
