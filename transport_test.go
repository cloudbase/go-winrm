package winrm

import (
	"fmt"
	"errors"
	"testing"
	"encoding/xml"

	gc "launchpad.net/gocheck"
)

func Test_transport(t *testing.T) { gc.TestingT(t) }

type TransportSuite struct{}

var _ = gc.Suite(TransportSuite{})

func (TransportSuite) TestGetHttpHeader(c *gc.C) {
	req := SoapRequest{}
	header := req.GetHttpHeader()
	want := make(map[string]string)
	want["Content-Type"] = "application/soap+xml;charset=UTF-8"
	want["User-Agent"] = "Go WinRM client"
	c.Assert(header, gc.DeepEquals, want)
}

// test for completely invalid protocol at Endpoint
func (TransportSuite) TestHttpCertAuthInvalidProtocol(c *gc.C) {
	req := SoapRequest{AuthType:"CertAuth", Endpoint:"nothttp://whatever.com"}
 
	resp, err := req.HttpCertAuth(nil)
	c.Assert(resp, gc.IsNil)
	c.Assert(err, gc.DeepEquals, errors.New("Invalid protocol. Expected http or https"))
}

// test for https protocol specifically
func (TransportSuite) TestHttpCertAuthNotHttps(c *gc.C) {
	req := SoapRequest{AuthType:"CertAuth", Endpoint:"http://something.smth"}

	resp, err := req.HttpCertAuth(nil)
	c.Assert(resp, gc.IsNil)
	c.Assert(err, gc.DeepEquals, errors.New("Invalid protocol for this transport type"))
}

// test for invalid key-value pair
func (TransportSuite) TestHttpCertAuthKeypairFailure(c *gc.C) {
	// must insert invalid certificate fields into this one:
	cert := CertificateCredentials{}
	req := SoapRequest{AuthType:"CertAuth", CertAuth:&cert, Endpoint:"https://something.good"}

	resp, err := req.HttpCertAuth(nil)
	c.Assert(resp, gc.IsNil)
	c.Assert(err, gc.NotNil)
}

// test for completely invalid protocol at Endpoint
func (TransportSuite) TestHttpBasicAuthInvalidProtocol(c *gc.C) {
	req := SoapRequest{AuthType:"BasicAuth", Endpoint:"nothttp://whatevs"}

	resp, err := req.HttpBasicAuth(nil)
	c.Assert(resp, gc.IsNil)
	c.Assert(err, gc.DeepEquals, errors.New("Invalid protocol. Expected http or https"))
}

// IRRELEVANT TEST
// compiler won't allow sending of envelope which is not a struct, thus the 
// MarshalIndent never throws the error
// ...
// when Envelope is anything but a struct(will pass even if not of type Envelope)
// func (TransportSuite) TestSendBadEnvelope(c *gc.C) {
// 	req := SoapRequest{}
// 	envelope := make(chan int)

// 	res, err := req.SendMessage(&envelope)
// 	c.Assert(res, gc.IsNil)
// 	c.Assert(err, gc.NotNil)
// }

// when BasicAuth is requested with no user/pass input
func (TransportSuite) TestSendEmptyAuthRequest(c *gc.C) {
	req := SoapRequest{AuthType:"BasicAuth"}
	envelope := &Envelope{}

	resp, err := req.SendMessage(envelope)
	c.Assert(resp, gc.IsNil)
	c.Assert(err, gc.DeepEquals, errors.New("AuthType BasicAuth needs Username and Passwd"))
}

// test that valid BasicAuth case is recognized
func (TransportSuite) TestSendBasicAuth(c *gc.C) {
	req := SoapRequest{AuthType:"BasicAuth", Username:"Leeroy", 
		Passwd:"Jenkins"}
	envelope := &Envelope{}

	resp, err := req.SendMessage(envelope)
	xmld, _ := xml.MarshalIndent(envelope, "  ", "    ")
	expresp, experr := req.HttpBasicAuth(xmld)

	c.Assert(resp, gc.DeepEquals, expresp)
	c.Assert(err, gc.DeepEquals, experr)
}

// test that valid CertAuth case is recognized
func (TransportSuite) TestSendCertAuth(c *gc.C) {
	req := SoapRequest{AuthType:"CertAuth"}
	envelope := &Envelope{}

	resp, err := req.SendMessage(envelope)
	xmld, _ := xml.MarshalIndent(envelope, "  ", "    ")
	expresp, experr := req.HttpCertAuth(xmld)

	c.Assert(resp, gc.DeepEquals, expresp)
	c.Assert(err, gc.DeepEquals, experr)
}

// test SoapRequest with bogus AuthType
func (TransportSuite) TestSendBadAuth(c *gc.C) {
	req := SoapRequest{AuthType:"SomeStupidShit"}
	envelope := &Envelope{}

	resp, err := req.SendMessage(envelope)
	c.Assert(resp, gc.IsNil)
	c.Assert(err, gc.DeepEquals, errors.New(fmt.Sprintf("Invalid transport: %s", req.AuthType)))
}

