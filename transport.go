package winrm

import (
    "net/http"
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
}

func (conf *SoapRequest) SendMessage(data []byte) (*http.Response, error){
    if conf.AuthType == "BasicAuth"{
        if conf.Username == "" || conf.Passwd == "" {
            fmt.Errorf("AuthType BasicAuth needs Username and Passwd")
        }
        return conf.HttpBasicAuth(data)
    }
    return nil, nil
}

func (conf *SoapRequest) GetHttpHeader () (map[string]string) {
    header := make(map[string]string)
    header["Content-Type"] = "application/soap+xml;charset=UTF-8"
    header["User-Agent"] = "Python WinRM client"
    return header
}

func (conf *SoapRequest) HttpBasicAuth (data []byte) (*http.Response, error){
    protocol := strings.Split(conf.Endpoint, ":")
    if protocol[0] != "http" && protocol[0] != "https"{
        return nil, errors.New("Invalid protocol. Expected http or https")
    }

    header := conf.GetHttpHeader()
    
    client := &http.Client{}
    // Ignore SSL certificate errors
    if protocol[0] == "https" {
        tr := &http.Transport{
            TLSClientConfig: &tls.Config{InsecureSkipVerify: conf.HttpInsecure},
        }
        client.Transport = tr
    }
    body := bytes.NewBuffer(data)
    req, err := http.NewRequest("POST", conf.Endpoint, body)
    req.ContentLength = int64(len(data))
    req.SetBasicAuth(conf.Username, conf.Passwd)

    for k, v := range header {
        req.Header.Add(k, v)
    }

    resp, err := client.Do(req)
    return resp, err
}
