package winrm

import (
	"bytes"
	"errors"
	"io"

	gc "launchpad.net/gocheck"
)

type responseSuite struct{}

var _ = gc.Suite(responseSuite{})

func (responseSuite) TestGetFromXML(c *gc.C) {
	xmlin := `<s:Envelope xml:lang="en-US" xmlns:s="http://www.w3.org/2003/05/soap-envelope" xmlns:a="http://schemas.xmlsoap.org/ws/2004/08/addressing" xmlns:x="http://schemas.xmlsoap.org/ws/2004/09/transfer" xmlns:w="http://schemas.dmtf.org/wbem/wsman/1/wsman.xsd" xmlns:rsp="http://schemas.microsoft.com/wbem/wsman/1/windows/shell" xmlns:p="http://schemas.microsoft.com/wbem/wsman/1/wsman.xsd"><s:Header><a:Action>http://schemas.microsoft.com/wbem/wsman/1/windows/shell/CommandResponse</a:Action><a:MessageID>uuid:EC452E31-2872-4921-8C0C-C76398695407</a:MessageID><a:To>http://schemas.xmlsoap.org/ws/2004/08/addressing/role/anonymous</a:To><a:RelatesTo>uuid:7261e275-6d36-a627-8de0-e382e3a3cc5a</a:RelatesTo></s:Header><s:Body><rsp:CommandResponse><rsp:CommandId>6D0A426F-4B4A-44F8-AF20-C35365258FEB</rsp:CommandId></rsp:CommandResponse></s:Body></s:Envelope>`
	res, err := GetObjectFromXML(bytes.NewBufferString(xmlin))
	c.Assert(err, gc.IsNil)
	c.Assert(res.Body.CommandResponse.CommandId, gc.Equals, "6D0A426F-4B4A-44F8-AF20-C35365258FEB")
}

func (responseSuite) TestGetFromXMLShellID(c *gc.C) {
	xmlin := `<s:Envelope xml:lang="en-US" xmlns:s="http://www.w3.org/2003/05/soap-envelope" xmlns:a="http://schemas.xmlsoap.org/ws/2004/08/addressing" xmlns:x="http://schemas.xmlsoap.org/ws/2004/09/transfer" xmlns:w="http://schemas.dmtf.org/wbem/wsman/1/wsman.xsd" xmlns:rsp="http://schemas.microsoft.com/wbem/wsman/1/windows/shell" xmlns:p="http://schemas.microsoft.com/wbem/wsman/1/wsman.xsd"><s:Header><a:Action>http://schemas.xmlsoap.org/ws/2004/09/transfer/CreateResponse</a:Action><a:MessageID>uuid:986227D1-1F7F-410D-8FCE-D45971E61B81</a:MessageID><a:To>http://schemas.xmlsoap.org/ws/2004/08/addressing/role/anonymous</a:To><a:RelatesTo>uuid:c7f03012-ba76-3191-6990-74cea9dd327d</a:RelatesTo></s:Header><s:Body><x:ResourceCreated><a:Address>http://windows-host:5985/wsman</a:Address><a:ReferenceParameters><w:ResourceURI>http://schemas.microsoft.com/wbem/wsman/1/windows/shell/cmd</w:ResourceURI><w:SelectorSet><w:Selector Name="ShellId">9731F5BD-E90B-403B-A8DB-010396CEBB4D</w:Selector></w:SelectorSet></a:ReferenceParameters></x:ResourceCreated><rsp:Shell xmlns:rsp="http://schemas.microsoft.com/wbem/wsman/1/windows/shell"><rsp:ShellId>9731F5BD-E90B-403B-A8DB-010396CEBB4D</rsp:ShellId><rsp:ResourceUri>http://schemas.microsoft.com/wbem/wsman/1/windows/shell/cmd</rsp:ResourceUri><rsp:Owner>ubuntu@ubuntu</rsp:Owner><rsp:ClientIP>192.168.1.1</rsp:ClientIP><rsp:IdleTimeOut>PT7200.000S</rsp:IdleTimeOut><rsp:InputStreams>stdin</rsp:InputStreams><rsp:OutputStreams>stdout stderr</rsp:OutputStreams><rsp:ShellRunTime>P0DT0H0M0S</rsp:ShellRunTime><rsp:ShellInactivity>P0DT0H0M0S</rsp:ShellInactivity></rsp:Shell></s:Body></s:Envelope>`
	res, err := GetObjectFromXML(bytes.NewBufferString(xmlin))
	c.Assert(err, gc.IsNil)
	c.Assert(res.Body.Shell.ShellId, gc.Equals, "9731F5BD-E90B-403B-A8DB-010396CEBB4D")
}

func (responseSuite) TestGerFromXMLError(c *gc.C) {
	xmlin := "Random junk"
	res, err := GetObjectFromXML(bytes.NewBufferString(xmlin))
	var _ = res
	c.Assert(res, gc.Equals, ResponseEnvelope{})
	c.Assert(err, gc.ErrorMatches, "Invalid server response")
}

func MockStdOut(XMLinput io.Reader) (ResponseEnvelope, error) {
	a := make([]ResponseStream, 3)
	a[0] = ResponseStream{Value: "c3VjaCBncmVhdA==", Name: "stdout", End: ""}
	a[1] = ResponseStream{Value: "", Name: "stdout", End: "true"}
	a[2] = ResponseStream{Value: "", Name: "stderr", End: "true"}
	return ResponseEnvelope{Body: &ResponseBody{ReceiveResponse: &ReceiveResponse{Stream: a, CommandState: &ResponseCommandState{ExitCode: 666}}}}, nil
}

func (responseSuite) TestStdOutParseCommandOutput(c *gc.C) {
	parseXML = MockStdOut
	stdout, stderr, exitcode, err := ParseCommandOutput(bytes.NewBufferString("mocked"))
	c.Assert(err, gc.IsNil)
	c.Assert(stdout, gc.Equals, "such great")
	c.Assert(stderr, gc.Equals, "")
	c.Assert(exitcode, gc.Equals, 666)
}

func MockStdErr(XMLinput io.Reader) (ResponseEnvelope, error) {
	a := make([]ResponseStream, 3)
	a[0] = ResponseStream{Value: "J3R5cGVyJyBpcyBub3QgcmVjb2duaXplZCBhcyBhbiBpbnRlcm5hbCBvciBleHRlcm5hbCBjb21tYW5kLA0Kb3BlcmFibGUgcHJvZ3JhbSBvciBiYXRjaCBmaWxlLg0K", Name: "stderr", End: ""}
	a[1] = ResponseStream{Value: "", Name: "stdout", End: "true"}
	a[2] = ResponseStream{Value: "", Name: "stderr", End: "true"}
	return ResponseEnvelope{Body: &ResponseBody{ReceiveResponse: &ReceiveResponse{Stream: a, CommandState: &ResponseCommandState{ExitCode: 666}}}}, nil
}

func (responseSuite) TestStdErrParseCommandOutput(c *gc.C) {
	parseXML = MockStdErr
	stdout, stderr, exitcode, err := ParseCommandOutput(bytes.NewBufferString("mocked"))
	c.Assert(err, gc.IsNil)
	c.Assert(stdout, gc.Equals, "")
	c.Assert(stderr, gc.Equals, "'typer' is not recognized as an internal or external command,\r\noperable program or batch file.\r\n")
	c.Assert(exitcode, gc.Equals, 666)
}

func MockMultipleStdOut(XMLinput io.Reader) (ResponseEnvelope, error) {
	a := make([]ResponseStream, 4)
	a[0] = ResponseStream{Value: "c3VjaCBncmVhdA0K", Name: "stdout", End: ""}
	a[1] = ResponseStream{Value: "bmVlZHMgbW9yZSBsaW5lcw==", Name: "stdout", End: ""}
	a[2] = ResponseStream{Value: "", Name: "stdout", End: "true"}
	a[3] = ResponseStream{Value: "", Name: "stderr", End: "true"}
	return ResponseEnvelope{Body: &ResponseBody{ReceiveResponse: &ReceiveResponse{Stream: a, CommandState: &ResponseCommandState{ExitCode: 666}}}}, nil
}

func (responseSuite) TestMultipleStdOutParseCommandOutput(c *gc.C) {
	parseXML = MockMultipleStdOut
	stdout, stderr, exitcode, err := ParseCommandOutput(bytes.NewBufferString("mocked"))
	c.Assert(err, gc.IsNil)
	c.Assert(stdout, gc.Equals, "such great\r\nneeds more lines")
	c.Assert(stderr, gc.Equals, "")
	c.Assert(exitcode, gc.Equals, 666)
}

func MockStdErrAndOut(XMLinput io.Reader) (ResponseEnvelope, error) {
	a := make([]ResponseStream, 4)
	a[0] = ResponseStream{Value: "c3VjaCBncmVhdA0K", Name: "stdout", End: ""}
	a[1] = ResponseStream{Value: "bmVlZHMgbW9yZSBsaW5lcw==", Name: "stderr", End: ""}
	a[2] = ResponseStream{Value: "", Name: "stdout", End: "true"}
	a[3] = ResponseStream{Value: "", Name: "stderr", End: "true"}
	return ResponseEnvelope{Body: &ResponseBody{ReceiveResponse: &ReceiveResponse{Stream: a, CommandState: &ResponseCommandState{ExitCode: 666}}}}, nil
}

func (responseSuite) TestStdErrAndOutParseCommandOutput(c *gc.C) {
	parseXML = MockStdErrAndOut
	stdout, stderr, exitcode, err := ParseCommandOutput(bytes.NewBufferString("mocked"))
	c.Assert(err, gc.IsNil)
	c.Assert(stdout, gc.Equals, "such great\r\n")
	c.Assert(stderr, gc.Equals, "needs more lines")
	c.Assert(exitcode, gc.Equals, 666)
}

func MockBreak(XMLinput io.Reader) (ResponseEnvelope, error) {
	a := make([]ResponseStream, 4)
	a[0] = ResponseStream{Value: "c3VjaCBncmVhdA0K", Name: "stdout", End: ""}
	a[1] = ResponseStream{Value: "", Name: "stdout", End: "true"}
	a[2] = ResponseStream{Value: "bmVlZHMgbW9yZSBsaW5lcw==", Name: "stdout", End: ""}
	a[3] = ResponseStream{Value: "", Name: "stderr", End: "true"}
	return ResponseEnvelope{Body: &ResponseBody{ReceiveResponse: &ReceiveResponse{Stream: a, CommandState: &ResponseCommandState{ExitCode: 666}}}}, nil
}

func (responseSuite) TestBreakParseCommandOutput(c *gc.C) {
	parseXML = MockBreak
	stdout, stderr, exitcode, err := ParseCommandOutput(bytes.NewBufferString("mocked"))
	c.Assert(err, gc.IsNil)
	c.Assert(stdout, gc.Equals, "such great\r\n")
	c.Assert(stderr, gc.Equals, "")
	c.Assert(exitcode, gc.Equals, 666)
}

func MockInvalidInput(XMLinput io.Reader) (ResponseEnvelope, error) {
	return ResponseEnvelope{}, errors.New("mock")
}

func (responseSuite) TestInvalidInputParseCommandOutput(c *gc.C) {
	parseXML = MockInvalidInput
	stdout, stderr, exitcode, err := ParseCommandOutput(bytes.NewBufferString("mocked"))
	c.Assert(err, gc.ErrorMatches, "Error parsing XML")
	c.Assert(stdout, gc.Equals, "")
	c.Assert(stderr, gc.Equals, "")
	c.Assert(exitcode, gc.Equals, 0)
}

func MockInvalidStdOut(XMLinput io.Reader) (ResponseEnvelope, error) {
	a := make([]ResponseStream, 3)
	a[0] = ResponseStream{Value: "0", Name: "stdout", End: ""}
	a[1] = ResponseStream{Value: "", Name: "stdout", End: "true"}
	a[2] = ResponseStream{Value: "", Name: "stderr", End: "true"}
	return ResponseEnvelope{Body: &ResponseBody{ReceiveResponse: &ReceiveResponse{Stream: a, CommandState: &ResponseCommandState{ExitCode: 666}}}}, nil
}

func (responseSuite) TestInvalidStdOutParseCommandOutput(c *gc.C) {
	parseXML = MockInvalidStdOut
	stdout, stderr, exitcode, err := ParseCommandOutput(bytes.NewBufferString("mocked"))
	c.Assert(err, gc.ErrorMatches, "Error decoding stdout")
	c.Assert(stdout, gc.Equals, "")
	c.Assert(stderr, gc.Equals, "")
	c.Assert(exitcode, gc.Equals, 0)
}

func MockInvalidStdErr(XMLinput io.Reader) (ResponseEnvelope, error) {
	a := make([]ResponseStream, 3)
	a[0] = ResponseStream{Value: "0", Name: "stderr", End: ""}
	a[1] = ResponseStream{Value: "", Name: "stdout", End: "true"}
	a[2] = ResponseStream{Value: "", Name: "stderr", End: "true"}
	return ResponseEnvelope{Body: &ResponseBody{ReceiveResponse: &ReceiveResponse{Stream: a, CommandState: &ResponseCommandState{ExitCode: 666}}}}, nil
}

func (responseSuite) TestInvalidStdErrParseCommandOutput(c *gc.C) {
	parseXML = MockInvalidStdErr
	stdout, stderr, exitcode, err := ParseCommandOutput(bytes.NewBufferString("mocked"))
	c.Assert(err, gc.ErrorMatches, "Error decoding stderr")
	c.Assert(stdout, gc.Equals, "")
	c.Assert(stderr, gc.Equals, "")
	c.Assert(exitcode, gc.Equals, 0)
}
