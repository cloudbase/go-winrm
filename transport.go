package winrm

import (
    "launchpad.net/gwacl/fork/http"
    "encoding/xml"
    "strings"
    "launchpad.net/gwacl/fork/tls"
    // "io/ioutil"
    // "crypto/x509"
    "fmt"
    "bytes"
    "errors"
)

type CertificateCredentials struct {
    Cert    string
    Key     string
    CA      string
}

type SoapRequest struct {
    Endpoint        string
    AuthType        string
    Username        string
    Passwd          string
    HttpInsecure    bool
    CertAuth        *CertificateCredentials
    HttpClient      *http.Client
}

func (conf *SoapRequest) SendMessage(envelope *Envelope) (*http.Response, error){
    output, err := xml.MarshalIndent(envelope, "  ", "    ")
    if err != nil{
        return nil, err
    }

    if conf.AuthType == "BasicAuth"{
        if conf.Username == "" || conf.Passwd == "" {
            // fmt.Errorf("AuthType BasicAuth needs Username and Passwd")
            return nil, errors.New("AuthType BasicAuth needs Username and Passwd")
        }
        return conf.HttpBasicAuth(output)
    }else if conf.AuthType == "CertAuth" {
        return conf.HttpCertAuth(output)
    }
    return nil, errors.New(fmt.Sprintf("Invalid transport: %s", conf.AuthType))
}

func (conf *SoapRequest) GetHttpHeader () (map[string]string) {
    header := make(map[string]string)
    header["Content-Type"] = "application/soap+xml;charset=UTF-8"
    header["User-Agent"] = "Go WinRM client"
    return header
}

func (conf *SoapRequest) HttpCertAuth(data []byte) (*http.Response, error) {
    protocol := strings.Split(conf.Endpoint, ":")
    if protocol[0] != "http" && protocol[0] != "https"{
        return nil, errors.New("Invalid protocol. Expected http or https")
    }
    header := conf.GetHttpHeader()
    header["Authorization"] = "http://schemas.dmtf.org/wbem/wsman/1/wsman/secprofile/https/mutual"

    if conf.HttpClient == nil{
        conf.HttpClient = &http.Client{}
    }

    if protocol[0] != "https" {
        return nil, errors.New("Ivalid protocol for this transport type")
    }

    cert, err := tls.LoadX509KeyPair(conf.CertAuth.Cert, conf.CertAuth.Key)
    if err != nil {
        return nil, err
    }

    tlsConfig := &tls.Config{
        InsecureSkipVerify: true,
        Certificates: []tls.Certificate{
            cert,
        },
    }

    tr := &http.Transport{
        TLSClientConfig: tlsConfig,
    }
    conf.HttpClient.Transport = tr
    body := bytes.NewBuffer(data)
    req, err := http.NewRequest("POST", conf.Endpoint, body)
    req.ContentLength = int64(len(data))
    for k, v := range header {
        req.Header.Add(k, v)
    }
    resp, err := conf.HttpClient.Do(req)
    if err != nil{
        return nil, err
    }
    if resp.StatusCode != 200 {
        return nil, errors.New(fmt.Sprintf("Remote host returned error status code: %d", resp.StatusCode))
    }
    fmt.Printf("%v\n%v\n", resp, err)
    return resp, err
}

func (conf *SoapRequest) HttpBasicAuth (data []byte) (*http.Response, error){
    protocol := strings.Split(conf.Endpoint, ":")
    if protocol[0] != "http" && protocol[0] != "https"{
        return nil, errors.New("Invalid protocol. Expected http or https")
    }

    header := conf.GetHttpHeader()
    
    if conf.HttpClient == nil{
        conf.HttpClient = &http.Client{}
    }
    // Ignore SSL certificate errors
    if protocol[0] == "https" {
        tr := &http.Transport{
            TLSClientConfig: &tls.Config{InsecureSkipVerify: conf.HttpInsecure},
        }
        conf.HttpClient.Transport = tr
    }
    body := bytes.NewBuffer(data)
    req, err := http.NewRequest("POST", conf.Endpoint, body)
    req.ContentLength = int64(len(data))
    req.SetBasicAuth(conf.Username, conf.Passwd)

    for k, v := range header {
        req.Header.Add(k, v)
    }

    resp, err := conf.HttpClient.Do(req)
    if err != nil{
        return nil, err
    }
    if resp.StatusCode != 200 {
        return nil, errors.New(fmt.Sprintf("Remote host returned error status code: %d", resp.StatusCode))
    }
    return resp, err
}
