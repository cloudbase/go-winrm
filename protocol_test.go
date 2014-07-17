package winrm

import (
	"testing"
	gc "launchpad.net/gocheck"
)

func Test_protocol(t *testing.T) { gc.TestingT(t) }

type ProtocolSuite struct {}

var _ = gc.Suite(ProtocolSuite{})


// tests adding ONLY ResourceURI to Envelope in GetSoapHeaders
func (ProtocolSuite) TestGetSoapHeadersResourceURI(c *gc.C) {
	env := Envelope{}
	params := HeaderParams{ResourceURI:"something"}
	exp := &ValueMustUnderstand{params.ResourceURI, "true"}

	env.GetSoapHeaders(params)
	c.Assert(env.Headers.ResourceURI, gc.DeepEquals, exp)
	c.Assert(env.Headers.MessageID, gc.Not(gc.Equals), "")
	c.Assert(env.Headers.SelectorSet, gc.IsNil)
	c.Assert(env.Headers.Action, gc.IsNil)
}

// tests adding ONLY Action to Envelope in GetSoapHeaders
func (ProtocolSuite) TestGetSoapHeadersAction(c *gc.C) {
	env := Envelope{}
	params := HeaderParams{Action:"yes"}
	exp := &ValueMustUnderstand{params.Action, "true"}

	env.GetSoapHeaders(params)
	c.Assert(env.Headers.Action, gc.DeepEquals, exp)
	c.Assert(env.Headers.MessageID, gc.Not(gc.Equals), "")
	c.Assert(env.Headers.ResourceURI, gc.IsNil)
	c.Assert(env.Headers.SelectorSet, gc.IsNil)
}

// tests adding ONLY ShellID to Envelope in GetSoapHeaders
func (ProtocolSuite) TestGetSoapHeadersShellID(c *gc.C) {
	env := Envelope{}
	params := HeaderParams{ShellID:"Not power shell"}
	exp := &Selector{ValueName{params.ShellID, "ShellId"}}

	env.GetSoapHeaders(params)
	c.Assert(env.Headers.SelectorSet, gc.DeepEquals, exp)
	c.Assert(env.Headers.MessageID, gc.Not(gc.Equals), "")
	c.Assert(env.Headers.ResourceURI, gc.IsNil)
	c.Assert(env.Headers.Action, gc.IsNil)
}

// test adding ONLY MessageID to Envelope in GetSoapHeaders
func (ProtocolSuite) TestGetSoapHeadersMessageID(c *gc.C) {
	env := Envelope{}
	params := HeaderParams{MessageID:"I am unique!"}

	env.GetSoapHeaders(params)
	c.Assert(params.MessageID, gc.NotNil)
	c.Assert(env.Headers.MessageID, gc.Equals, params.MessageID)
	c.Assert(env.Headers.ResourceURI, gc.IsNil)
	c.Assert(env.Headers.SelectorSet, gc.IsNil)
	c.Assert(env.Headers.Action, gc.IsNil)
}


// tests if Envelope attributes are succesfully configured by GetShell
func (ProtocolSuite) TestGetShellMakeEnvelope(c *gc.C) {
	req := SoapRequest{}
	env := Envelope{}
	params := ShellParams{EnvVars:&Environment{Variable:[]EnvVariable{EnvVariable{Value:"Yes", Name:"Johann"}}}}
	expparams := ShellParams{IStream:"stdin", OStream:"stdout stderr", EnvVars:&Environment{[]EnvVariable{EnvVariable{Value:"Yes", Name:"Johann"}}}}

	_, _ = env.GetShell(params, req)

	// c.Assert(err, gc.IsNil)
	// c.Assert(shell, gc.NotNil)
	c.Assert(env.EnvelopeAttrs, gc.Equals, Namespaces)
	c.Assert(env.Body.Shell, gc.DeepEquals, &Shell{InputStreams:expparams.IStream, OutputStreams: expparams.OStream, Environment:expparams.EnvVars})
}


// tests if missing ShellID parameter is signaled by SendCommand
func (ProtocolSuite) TestSendCommandNoShellId(c *gc.C) {
	env := Envelope{}
	params := CmdParams{}
	req := SoapRequest{}

	comId, err := env.SendCommand(params, req)
	c.Assert(comId, gc.Equals, "")
	c.Assert(err, gc.ErrorMatches, "Invalid ShellId")
}

// tests if Timeout parameter succesfully added to Envelope by SendCommand
func (ProtocolSuite) TestSendCommandTimeoutSetting(c *gc.C) {
	env := Envelope{}
	params := CmdParams{Timeout:"Some timeout", ShellID:"Something"}
	req := SoapRequest{}

	_, _ = env.SendCommand(params, req)
	c.Assert(env.Headers.OperationTimeout, gc.Equals, params.Timeout)
}

// tests if default Timeout is set by SendCommand when none is provided
func (ProtocolSuite) TestSendCommandDefaultTimeout(c *gc.C) {
	env := Envelope{}
	params := CmdParams{ShellID:"Something"}
	req := SoapRequest{}

	_, _ = env.SendCommand(params, req)
	c.Assert(env.Headers.OperationTimeout, gc.Equals, "PT3600S")
}

// tests if empty Command parameter succesfully signaled by SendCommand
func (ProtocolSuite) TestSendCommandMissingCommand(c *gc.C) {
	env := Envelope{}
	params := CmdParams{ShellID:"Absolutely Something"}
	req := SoapRequest{}

	comId, err := env.SendCommand(params, req)
	c.Assert(comId, gc.Equals, "")
	c.Assert(err, gc.ErrorMatches, "Invalid command")
}

// tests if Args parameter succesfully added to envelope by SendCommand
func (ProtocolSuite) TestSendCommandArgsPassing(c *gc.C) {
	env := Envelope{}
	params := CmdParams{ShellID:"Something", Args:"Pepsi is better that Coke", Cmd:"More Something"}
	req := SoapRequest{}

	_, _ = env.SendCommand(params, req)
	c.Assert(env.Body.CommandLine.Arguments, gc.Equals, params.Args)
}


// tests if Envelope attributes succesfully updated by GetCommandOutput
func (ProtocolSuite) TestGetCommandOutputMakeEnvelope(c *gc.C) {
	env := Envelope{}
	req := SoapRequest{}

	expparams := HeaderParams{ResourceURI:"http://schemas.microsoft.com/wbem/wsman/1/windows/shell/cmd", Action:"http://schemas.microsoft.com/wbem/wsman/1/windows/shell/Recieve", ShellID:""}
	exprec := Receive{DesiredStream:DesiredStreamProps{Value:"stdout stderr", Attr:""}}
	expenv := Envelope{EnvelopeAttrs:Namespaces, Body:&BodyStruct{Receive:&exprec}}
	expenv.GetSoapHeaders(expparams)

	_, _, _, _ = env.GetCommandOutput("", "", req)
	c.Assert(env.Body, gc.DeepEquals, expenv.Body)
	c.Assert(env.EnvelopeAttrs, gc.Equals, expenv.EnvelopeAttrs)
	c.Assert(env.Headers.Locale, gc.DeepEquals, expenv.Headers.Locale)
	c.Assert(env.Headers.ReplyTo, gc.DeepEquals, expenv.Headers.ReplyTo)
	c.Assert(env.Headers.DataLocale, gc.DeepEquals, expenv.Headers.DataLocale)
}


// tests if Envelope attributes succesfully updated by CleanupShell
func (ProtocolSuite) TestCleanupShellMakeEnvelope(c *gc.C) {
	env := Envelope{}
	req := SoapRequest{}

	expparams := HeaderParams{ResourceURI:"http://schemas.microsoft.com/wbem/wsman/1/windows/shell/cmd", Action:"http://schemas.microsoft.com/wbem/wsman/1/windows/shell/Signal", ShellID:""}
	expsig := Signal{Attr:"", Code:"http://schemas.microsoft.com/wbem/wsman/1/windows/shell/signal/terminate"}
	expenv := Envelope{EnvelopeAttrs:Namespaces, Body:&BodyStruct{Signal:&expsig}}
	expenv.GetSoapHeaders(expparams)

	_ = env.CleanupShell("", "", req)
	c.Assert(env.Body, gc.DeepEquals, expenv.Body)
	c.Assert(env.EnvelopeAttrs, gc.Equals, expenv.EnvelopeAttrs)
	c.Assert(env.Headers.Locale, gc.DeepEquals, expenv.Headers.Locale)
	c.Assert(env.Headers.ReplyTo, gc.DeepEquals, expenv.Headers.ReplyTo)
	c.Assert(env.Headers.DataLocale, gc.DeepEquals, expenv.Headers.DataLocale)
}


// tests if Envelope attributes succesfully updated by CloseShell
func (ProtocolSuite) TestCloseShellMakeEnvelope(c *gc.C) {
	env := Envelope{}
	req := SoapRequest{}

	expparams := HeaderParams{ResourceURI: "http://schemas.microsoft.com/wbem/wsman/1/windows/shell/cmd", Action:"http://schemas.xmlsoap.org/ws/2004/09/transfer/Delete", ShellID:""}
	expenv := Envelope{EnvelopeAttrs:Namespaces, Body:&BodyStruct{}}
	expenv.GetSoapHeaders(expparams)

	_ = env.CloseShell("", req)
	c.Assert(env.Body, gc.DeepEquals, expenv.Body)
	c.Assert(env.EnvelopeAttrs, gc.Equals, expenv.EnvelopeAttrs)
	c.Assert(env.Headers.Locale, gc.DeepEquals, expenv.Headers.Locale)
	c.Assert(env.Headers.ReplyTo, gc.DeepEquals, expenv.Headers.ReplyTo)
	c.Assert(env.Headers.DataLocale, gc.DeepEquals, expenv.Headers.DataLocale)
}
