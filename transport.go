package winrm

import (
    "net/http"
    "encoding/xml"
    "strings"
    "crypto/tls"
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
    }
    return nil, errors.New(fmt.Sprintf("Invalid transport: %s", conf.AuthType))
}

func (conf *SoapRequest) GetHttpHeader () (map[string]string) {
    header := make(map[string]string)
    header["Content-Type"] = "application/soap+xml;charset=UTF-8"
    header["User-Agent"] = "Go WinRM client"
    return header
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
    return resp, err
}
