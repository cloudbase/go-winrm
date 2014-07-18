package winrm

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"syscall"

	gc "launchpad.net/gocheck"
)

type TransportSuite struct{}

var _ = gc.Suite(TransportSuite{})

//tests that GetHttpHeader outputs the desired values
func (TransportSuite) TestGetHttpHeader(c *gc.C) {
	req := SoapRequest{}
	header := req.GetHttpHeader()
	want := make(map[string]string)
	want["Content-Type"] = "application/soap+xml;charset=UTF-8"
	want["User-Agent"] = "Go WinRM client"
	c.Assert(header, gc.DeepEquals, want)
}

// tests for completely invalid protocol in HttpCertAuth
func (TransportSuite) TestHttpCertAuthInvalidProtocol(c *gc.C) {
	req := SoapRequest{
		AuthType: "CertAuth",
		Endpoint: "nothttp://whatever.com",
	}

	resp, err := req.HttpCertAuth(nil)
	c.Assert(resp, gc.IsNil)
	c.Assert(err, gc.ErrorMatches, "Invalid protocol. Expected http or https")
}

// tests that https protocol is specifically asked for in HttpCertAuth
func (TransportSuite) TestHttpCertAuthNotHttps(c *gc.C) {
	req := SoapRequest{
		AuthType: "CertAuth",
		Endpoint: "http://something.smth",
	}

	resp, err := req.HttpCertAuth(nil)
	c.Assert(resp, gc.IsNil)
	c.Assert(err, gc.ErrorMatches, "Invalid protocol for this transport type")
}

// test for invalid key-value pair in HttpCertAuth
func (TransportSuite) TestHttpCertAuthKeypairFailure(c *gc.C) {
	// must insert invalid certificate fields into this one:
	cert := CertificateCredentials{}
	req := SoapRequest{
		AuthType: "CertAuth",
		CertAuth: &cert,
		Endpoint: "https://something.good",
	}

	resp, err := req.HttpCertAuth(nil)
	c.Assert(resp, gc.IsNil)
	c.Assert(err, gc.NotNil)
}

// test for completely invalid protocol in HttpBasicAuth
func (TransportSuite) TestHttpBasicAuthInvalidProtocol(c *gc.C) {
	req := SoapRequest{
		AuthType: "BasicAuth",
		Endpoint: "nothttp://whatevs",
	}

	resp, err := req.HttpBasicAuth(nil)
	c.Assert(resp, gc.IsNil)
	c.Assert(err, gc.ErrorMatches, "Invalid protocol. Expected http or https")
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

// tests that alert is raised in case of BasicAuth request with empty user/pass in SendMessage
func (TransportSuite) TestSendMessageEmptyAuthRequest(c *gc.C) {
	req := SoapRequest{AuthType: "BasicAuth"}
	envelope := &Envelope{}

	resp, err := req.SendMessage(envelope)
	c.Assert(resp, gc.IsNil)
	c.Assert(err, gc.ErrorMatches, "AuthType BasicAuth needs Username and Passwd")
}

// tests that valid BasicAuth case is recognized in SendMessage
func (TransportSuite) TestSendMessageBasicAuth(c *gc.C) {
	req := SoapRequest{
		AuthType: "BasicAuth",
		Username: "Leeroy",
		Passwd:   "Jenkins"}
	envelope := &Envelope{}

	resp, err := req.SendMessage(envelope)
	xmld, _ := xml.MarshalIndent(envelope, "  ", "    ")
	expresp, experr := req.HttpBasicAuth(xmld)

	c.Assert(resp, gc.DeepEquals, expresp)
	c.Assert(err, gc.DeepEquals, experr)
}

// test that valid CertAuth case is recognized in SendMessage
func (TransportSuite) TestSendMessageCertAuth(c *gc.C) {
	req := SoapRequest{AuthType: "CertAuth"}
	envelope := &Envelope{}

	resp, err := req.SendMessage(envelope)
	xmld, _ := xml.MarshalIndent(envelope, "  ", "    ")
	expresp, experr := req.HttpCertAuth(xmld)

	c.Assert(resp, gc.DeepEquals, expresp)
	c.Assert(err, gc.DeepEquals, experr)
}

// tests that SoapRequest with bogus AuthType is rejected in SendMessage
func (TransportSuite) TestSendMessageBadAuth(c *gc.C) {
	req := SoapRequest{
		AuthType: "SomeStupidShit",
	}
	envelope := &Envelope{}

	resp, err := req.SendMessage(envelope)
	c.Assert(resp, gc.IsNil)
	c.Assert(err, gc.ErrorMatches, fmt.Sprintf("Invalid transport: %s", req.AuthType))
}

func (TransportSuite) TestHttpBasicRequestOK(c *gc.C) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.Assert(r.Method, gc.Equals, "POST")

		body, err := ioutil.ReadAll(r.Body)
		c.Assert(err, gc.IsNil)
		c.Assert(string(body), gc.Equals, "trololol")

		c.Assert(r.ContentLength, gc.Equals, int64(len(body)))

		c.Assert(r.Header.Get("User-Agent"), gc.Equals, "Go WinRM client")
		c.Assert(r.Header.Get("Content-Type"), gc.Equals, "application/soap+xml;charset=UTF-8")
		c.Assert(r.Header.Get("Authorization"), gc.Equals, "Basic bGVlcm95OmplbmtpbnM=")
	}))
	defer server.Close()

	req := SoapRequest{
		Endpoint:   server.URL,
		AuthType:   "BasicAuth",
		HttpClient: nil,
		Username:   "leeroy",
		Passwd:     "jenkins",
	}
	body := []byte("trololol")

	resp, err := req.HttpBasicAuth(body)
	c.Assert(err, gc.IsNil)
	c.Assert(resp, gc.NotNil)
}

func (TransportSuite) TestHttpBasicServerError(c *gc.C) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "fail", http.StatusInternalServerError)
	}))
	defer server.Close()

	req := SoapRequest{
		Endpoint:   server.URL,
		AuthType:   "BasicAuth",
		HttpClient: nil,
		Username:   "leeroy",
		Passwd:     "jenkins",
	}
	body := []byte("trololol")

	resp, err := req.HttpBasicAuth(body)
	c.Assert(err, gc.ErrorMatches, "Remote host returned error status code: 500")
	c.Assert(resp, gc.IsNil)
}

func (TransportSuite) TestHttpBasicError(c *gc.C) {
	req := SoapRequest{
		Endpoint:   "http://doesnotexist",
		AuthType:   "BasicAuth",
		HttpClient: nil,
		Username:   "leeroy",
		Passwd:     "jenkins",
	}

	body := []byte("trololol")

	resp, err := req.HttpBasicAuth(body)
	c.Assert(err, gc.ErrorMatches, "dial tcp: lookup doesnotexist: no such host")
	c.Assert(resp, gc.IsNil)
}

func (TransportSuite) TestHttpCertRequestOK(c *gc.C) {

	pem, err := ioutil.TempFile("", "pem")
	if err != nil {
		panic(err)
	}
	defer syscall.Unlink(pem.Name())
	ioutil.WriteFile(pem.Name(), []byte(cert_pem), 0644)

	key, err := ioutil.TempFile("", "key")
	if err != nil {
		panic(err)
	}
	defer syscall.Unlink(key.Name())
	ioutil.WriteFile(key.Name(), []byte(cert_key), 0644)

	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.Assert(r.Method, gc.Equals, "POST")

		body, err := ioutil.ReadAll(r.Body)
		c.Assert(err, gc.IsNil)
		c.Assert(string(body), gc.Equals, "trololol")

		c.Assert(r.ContentLength, gc.Equals, int64(len(body)))

		c.Assert(r.Header.Get("User-Agent"), gc.Equals, "Go WinRM client")
		c.Assert(r.Header.Get("Content-Type"), gc.Equals, "application/soap+xml;charset=UTF-8")
		c.Assert(r.Header.Get("Authorization"), gc.Equals, "http://schemas.dmtf.org/wbem/wsman/1/wsman/secprofile/https/mutual")
	}))
	defer server.Close()

	req := SoapRequest{
		Endpoint:   server.URL,
		AuthType:   "CertAuth",
		HttpClient: nil,
		CertAuth: &CertificateCredentials{
			Cert: pem.Name(),
			Key:  key.Name(),
		}}
	body := []byte("trololol")

	resp, err := req.HttpCertAuth(body)
	c.Assert(err, gc.IsNil)
	c.Assert(resp, gc.NotNil)
}

func (TransportSuite) TestHttpCertRequestServerError(c *gc.C) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "fail", http.StatusInternalServerError)
	}))
	defer server.Close()

	pem, err := ioutil.TempFile("", "pem")
	if err != nil {
		panic(err)
	}
	defer syscall.Unlink(pem.Name())
	ioutil.WriteFile(pem.Name(), []byte(cert_pem), 0644)

	key, err := ioutil.TempFile("", "key")
	if err != nil {
		panic(err)
	}
	defer syscall.Unlink(key.Name())
	ioutil.WriteFile(key.Name(), []byte(cert_key), 0644)

	req := SoapRequest{
		Endpoint:   server.URL + "/a",
		AuthType:   "CertAuth",
		HttpClient: nil,
		CertAuth: &CertificateCredentials{
			Cert: pem.Name(),
			Key:  key.Name(),
		}}
	body := []byte("trololol")

	resp, err := req.HttpCertAuth(body)
	c.Assert(err, gc.ErrorMatches, "Remote host returned error status code: 500")
	c.Assert(resp, gc.IsNil)
}

func (TransportSuite) TestHttpCertRequestError(c *gc.C) {
	pem, err := ioutil.TempFile("", "pem")
	if err != nil {
		panic(err)
	}
	defer syscall.Unlink(pem.Name())
	ioutil.WriteFile(pem.Name(), []byte(cert_pem), 0644)

	key, err := ioutil.TempFile("", "key")
	if err != nil {
		panic(err)
	}
	defer syscall.Unlink(key.Name())
	ioutil.WriteFile(key.Name(), []byte(cert_key), 0644)

	req := SoapRequest{
		Endpoint:   "https://doesnotexist",
		AuthType:   "CertAuth",
		HttpClient: nil,
		CertAuth: &CertificateCredentials{
			Cert: pem.Name(),
			Key:  key.Name(),
		}}
	body := []byte("trololol")

	resp, err := req.HttpCertAuth(body)
	c.Assert(err, gc.ErrorMatches, "dial tcp: lookup doesnotexist: no such host")
	c.Assert(resp, gc.IsNil)
}

var cert_pem = `-----BEGIN CERTIFICATE-----
MIIC8DCCAdigAwIBAwICA+gwDQYJKoZIhvcNAQEFBQAwGDEWMBQGA1UEAxQNdWJ1
bnR1QHVidW50dTAeFw0xNDA3MTgxMjEyMzlaFw0yNDA3MTUxMjEyMzlaMBgxFjAU
BgNVBAMUDXVidW50dUB1YnVudHUwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEK
AoIBAQDJ6tBRTRs6VZDYB3tVY0r2CPtzxrSNbKjvsaYdgqYwx3eX7bBhpM9/ijLZ
iycMM/O8SBxzTxZv5nF82RXVXY7dAUq8lPbacwRoo1PAAQX1CZrsmsdyEjXn7M9n
E65FPKTz3yEyEVYK9f43lmVwAntACHenl0bH+0dxgo/eQTGVrVYPJp5aa4OpIUQD
Ph2txDx4oyYgdGP78W4/4tPAwxyljy00F6tdo5OD7Wjw1C01sp50iO6dRb83twAN
KnNxVwQqnWgn3XuzXaNFwOqlQ1yalgUk2Ky/yrMe2FhmNIaGnYv5z0lBS9jBbsmx
l0FZ5btN6QTJAM/IlC+XZaa3ftlrAgMBAAGjRDBCMBYGA1UdJQEB/wQMMAoGCCsG
AQUFBwMCMCgGA1UdEQQhMB+gHQYKKwYBBAGCNxQCA6APDA11YnVudHVAdWJ1bnR1
MA0GCSqGSIb3DQEBBQUAA4IBAQAih7AxVMFHevX2GbcjlU9s5i6+ZKKhh5Jo0iMH
voS8dk0E3V36TABZwqc4jY3BmHvit0esBkpQOP2I4F634ByUEe7462rtUBrgIBHd
WFPEHx/dwq7S+iktOyOvnk2uEyGCH8B4EMeiCpsLzC9g3bjsuySAB5l1HJJROAVH
sjKgBsCgJvdFah0UKv1xXzSZdBjMWw8b4tVdIGJ6N3S4bLfZbwrR8c/Ym9SJdXfM
PLWfcd5kJD8awgKIDUrmCn4LJpCLFSfRrbi1Y+AzsdVvsq3bpzfD5ZMFqxa7m6sr
YKmtPl60PWQcfxP3yWUJNQB9hPTav1u3+NlUxUcP/Vw+6A+U
-----END CERTIFICATE-----`

var cert_key = `-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDJ6tBRTRs6VZDY
B3tVY0r2CPtzxrSNbKjvsaYdgqYwx3eX7bBhpM9/ijLZiycMM/O8SBxzTxZv5nF8
2RXVXY7dAUq8lPbacwRoo1PAAQX1CZrsmsdyEjXn7M9nE65FPKTz3yEyEVYK9f43
lmVwAntACHenl0bH+0dxgo/eQTGVrVYPJp5aa4OpIUQDPh2txDx4oyYgdGP78W4/
4tPAwxyljy00F6tdo5OD7Wjw1C01sp50iO6dRb83twANKnNxVwQqnWgn3XuzXaNF
wOqlQ1yalgUk2Ky/yrMe2FhmNIaGnYv5z0lBS9jBbsmxl0FZ5btN6QTJAM/IlC+X
Zaa3ftlrAgMBAAECggEBAICLRpebmOvoMU/2Y2QW1FARo9Mu+x7VwC7oT7KVzCtd
sRs9rH5dJ+QwHPM1jWRNZqvE1Kfr/4K5mCI9KZMt/pdgDS5FP2oOsw3SfKzNefdn
aAOc/b/3K+48akVa2CUn2HOQ51cyhi5wMKk+y9ElI0W+nj5JJjyGEhOHZQO/SUvZ
blcLLjdyIpj+g5eKkNTDSXo80M3lkedPewFb9ihFH4NF9wa5RP93DbiA85FFdCHc
0k6nKReKJAXyc6VUTJ3O0Bl1m64ByAonuIBMqJ9LYKzk9gCwIAGeTnj153x7wc7o
eebtxfjRB5TMTlh/wn9t855aGMbTr6+URbQh2z6GvxECgYEA6BYR+0K+NgZZsDCp
e87BkMAWnrR/2Ihc89G17G0JVF0rgxpFkcbwO9MkrG658gttG+auTInw9cgWqFj6
0ZXK8PPelIpdoPIIrVJivZNFoN249Z/uD1WmNekw/wj6nCPvo8LOfP4yHQ+BjWwF
3IuFas109B1QfBuF2TEXHFiAGycCgYEA3rjznlPE81iT+8zl2BqZ4Zi6A4iIpBAo
jxc8miL5grltW3Wl5iQk1hRdaFj2IaTouPt6/TBI7Dzw+F3IBckRycD2oZ9npGRW
oCvcE/y5ar4qowb4n/she2K2jQzZ6Yjd+48HsGbXgzvl7T/S+tt0gdtQl80NfZBX
KYUAQ11Zyh0CgYEA5W8MD6yHhbj5aShyJCbdTE/ZDMO7rz//RDoI8tVH59LDdTO/
msFkNIAjPSOpRxLspiyCGsAzKYbIf1yXeCHxIgqz+3xd2wHqeg1795VjvAf1FT0p
hpdRXPJOsZEazsjn2qh2oTJaMEhn9nrXwJNdLZw3Bi0Ep+w9gdz5z9fdrPkCgYBa
4lAPQJGy12dzrdXwzFIU29S0Emfnwuw6D7pcD3+Pl4kHdEehVQhvD1padUriyb9p
lL1ISgbH18phHyu7KKSIlqRNqZWKYKN0stEYmt0ysK0HX5Xe+oRcLBjgD+lwQbiL
qX7yvdSdqbiWip/WW+z7/HmzqCokHd1jhPFpi9NTBQKBgEcVlVHVqSe9WOTYzx8q
O/tYCr2cszggryOeI9CgqX8KzvaGpXWc/w6iWM8JBJ6qZd+hoKKnxuhjPLxyrZPj
9zeMBJcPd67cHzd8ZC55dMVucHQYMr4rHeua8qouat3iPLAfCtr37pY888+jDZvn
eB6Glbz65UWCjd0GL8v8WRfV
-----END PRIVATE KEY-----`
