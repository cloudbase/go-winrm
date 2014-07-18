package winrm

import (
	"regexp"
	"testing"

	gc "launchpad.net/gocheck"

	jc "launchpad.net/juju-core/testing/checkers"
)

type utilSuite struct{}

var _ = gc.Suite(utilSuite{})

func IsValidUUID(s string) bool {
	var validUUID = regexp.MustCompile("[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}")

	return validUUID.MatchString(s)
}

func (utilSuite) TestUUID(c *gc.C) {
	uuid, err := Uuid()
	c.Assert(err, gc.IsNil)
	c.Assert(uuid, jc.Satisfies, IsValidUUID)
}
