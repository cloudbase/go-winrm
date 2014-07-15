package winrm

import (
	"testing"

	gc "launchpad.net/gocheck"
)

func Test_transport(t *testing.T) { gc.TestingT(t) }

type transportSuite struct{}

var _ = gc.Suite(transportSuite{})

func (transportSuite) TestUUID(c *gc.C) {
	req := SoapRequest{}
	header := req.GetHttpHeader()
	c.Assert(header["Content-Type"], gc.DeepEquals, "application/soap+xml;charset=UTF-8")
	c.Assert(header["User-Agent"], gc.DeepEquals, "Go WinRM client")
}
